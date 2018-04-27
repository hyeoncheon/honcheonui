package models

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/hyeoncheon/honcheonui/utils"
)

// Resource is struct for storing atomic monitoring target resources.
// For support general resources from different kind of services,
// it just contains common attributes and provider specific attributes
// are stored separately on Attribute model.
type Resource struct {
	ID                 uuid.UUID  `json:"id" db:"id"`
	CreatedAt          time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at" db:"updated_at"`
	Provider           string     `json:"provider" db:"provider"`
	Type               string     `json:"type" db:"type"`
	OriginalID         string     `json:"original_id" db:"original_id"`
	UUID               uuid.UUID  `json:"uuid" db:"uuid"`
	Name               string     `json:"name" db:"name"`
	Notes              string     `json:"notes" db:"notes"`
	GroupID            string     `json:"group_id" db:"group_id"`
	ResourceCreatedAt  time.Time  `json:"resource_created_at" db:"resource_created_at"`
	ResourceModifiedAt time.Time  `json:"resource_modified_at" db:"resource_modified_at"`
	IPAddress          string     `json:"ip_address" db:"ip_address"`
	Location           string     `json:"location" db:"location"`
	IsConn             bool       `json:"is_conn" db:"is_conn"`
	IsOn               bool       `json:"is_on" db:"is_on"`
	Tags               Tags       `many_to_many:"resources_tags"`
	Attributes         Attributes `has_many:"attributes"`
	Providers          Providers  `many_to_many:"providers_resources"`
	Incidents          Incidents  `many_to_many:"incidents_resources"`
}

// ResourcesTags is struct for mapping resources to tags.
// This model is supporting many_to_many association.
type ResourcesTags struct {
	ID         uuid.UUID `db:"id"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
	ResourceID uuid.UUID `db:"resource_id"`
	TagID      uuid.UUID `db:"tag_id"`
}

// ResourcesUsers is struct for mapping resources to users.
// User part of this relationship model is not present currently.
type ResourcesUsers struct {
	ID         uuid.UUID `db:"id"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
	ResourceID uuid.UUID `db:"resource_id"`
	UserID     string    `db:"user_id"`
	Provider   Provider  `belongs_to:"providers" fk_id:"user_id"`
}

// String represents its name and provider.
func (r Resource) String() string {
	return r.Name + "@" + r.Provider
}

// JSON returns json marshalled string for this model.
func (r Resource) JSON() string {
	jr, _ := json.Marshal(r)
	return string(jr)
}

// Resources is an array of resources
type Resources []Resource

// RTMaps is an array of resources-tags map
type RTMaps []ResourcesTags

// RUMaps is an array of resources-users map
type RUMaps []ResourcesUsers

//*** relational methods

// Services returns indirectly associated services for the resource.
// resources have association with services via tags and matching rule of
// services.
func (r *Resource) Services() *Services {
	svcs := &Services{}
	services := &Services{}
	if len(r.Tags) < 1 {
		return services
	}

	var IDs []interface{}
	for _, t := range r.Tags {
		IDs = append(IDs, t.ID)
	}

	query := DB.Q().
		Join("services_tags", "services_tags.service_id = services.id").
		Where("services_tags.tag_id in (?)", IDs...).
		GroupBy("services.id")
	if err := query.All(svcs); err != nil {
		logger.Errorf("could not get resources. error: %v", err)
	}

	for _, svc := range *svcs {
		if svc.HasResource(r) {
			*services = append(*services, svc)
		}
	}
	DB.Load(services, "Member")

	return services
}

// AddAttribute creates and saves attribute of the resource.
func (r *Resource) AddAttribute(name, value string) error {
	attr := &Attribute{
		ResourceID: r.ID,
		Name:       name,
		Value:      value,
	}
	verrs, err := DB.ValidateAndCreate(attr)
	if verrs.HasAny() {
		return errors.New("validation error")
	}
	return err
}

// LinkTags makes a link map of resource and user.
// Since this models are mirrored from its origin, this function should
// handle duplications and updates by checking existing one.
func (r *Resource) LinkTags(ts []string) error {
	requestedTags, err := utils.ToInterface(utils.Cleaner(ts))
	if err != nil {
		return errors.New("could not convert argument to interface slice")
	}
	logger.Debugf("link tags requested: %v", requestedTags)

	hasError := false
	// get existing tag map and remove them from given list.
	existingMaps := &RTMaps{}
	if err := DB.Where("resource_id = ?", r.ID).All(existingMaps); err != nil {
		logger.Errorf("database selection failed! error: %v", err)
	}
	for _, m := range *existingMaps {
		existingTag := &Tag{}
		err := DB.Find(existingTag, m.TagID)
		if err != nil { // in case of broken link map, but why? safety?
			err := DB.Destroy(&m)
			if err != nil {
				logger.Errorf("found broken map but could not delete: id:%v", m.ID)
			} else {
				logger.Warnf("found broken map and deleted: id:%v", m.ID)
			}
		}

		if utils.Has(requestedTags, existingTag.Name) {
			requestedTags = utils.Remove(requestedTags, existingTag.Name)
		} else {
			logger.Debugf("removing %v from map. no longer exists", existingTag.Name)
			if err := DB.Destroy(&m); err != nil {
				logger.Errorf("could not remove  %v from the map", m)
				hasError = true
			}
		}
	}

	logger.Debugf("adding new tags...: %v", requestedTags)
	for _, t := range requestedTags {
		name := strings.TrimSpace(t.(string)) //! check me

		// search existing tag entry or create new one.
		tag := &Tag{}
		if err := DB.Where("name = ?", name).First(tag); err != nil { // if none
			logger.Debugf("create new tag %v...", name)
			tag.Name = name
			verrs, err := DB.ValidateAndCreate(tag)
			if verrs.HasAny() {
				logger.Errorf("could not save resource-tag map %v: %v", tag, verrs)
			}
			if err != nil {
				logger.Errorf("could not save resource-tag map %v: %v", tag, err)
			}
		}

		logger.Debugf("creating map for %v on %v", name, r)
		rtmap := &ResourcesTags{
			ResourceID: r.ID,
			TagID:      tag.ID,
		}
		if err := DB.Save(rtmap); err != nil {
			logger.Errorf("could not save resource-tag map %v: %v", rtmap, err)
			hasError = true
		}
	}
	if hasError {
		return errors.New("linking done with error(s)")
	}
	return nil
}

// LinkUsers makes a link map of resource and user.
// the logic is same as `LinkTags()`.
func (r *Resource) LinkUsers(us []string) error {
	requestedUsers, err := utils.ToInterface(us)
	if err != nil {
		return errors.New("could not convert argument to interface slice")
	}
	logger.Debugf("link users requested: %v", requestedUsers)

	hasError := false
	// get existing user map and remove them from given list
	maps := &RUMaps{}
	if err := DB.Where("resource_id = ?", r.ID).All(maps); err != nil {
		logger.Errorf("database selection failed! error: %v", err)
	}
	for _, m := range *maps {
		if utils.Has(requestedUsers, m.UserID) {
			requestedUsers = utils.Remove(requestedUsers, m.UserID)
		} else {
			logger.Debugf("removing %v from map. no longer exists", m.UserID)
			if err := DB.Destroy(&m); err != nil {
				logger.Errorf("could not remove user %v from the map", m)
				hasError = true
			}
		}
	}

	logger.Debugf("adding new users...: %v", requestedUsers)
	for _, u := range requestedUsers {
		logger.Debugf("creating map for %v on %v", u, r)
		rumap := &ResourcesUsers{
			ResourceID: r.ID,
			UserID:     u.(string), //! check me
		}
		if err := DB.Save(rumap); err != nil {
			logger.Errorf("could not save resource-user map %v: %v", rumap, err)
			hasError = true
		}
	}
	if hasError {
		return errors.New("linking done with error(s)")
	}
	return nil
}

//*** common database functions and methods

// Save stores the resource
//! be smart!
func (r *Resource) Save() error {
	if err := DB.Create(r); err != nil {
		return DB.Update(r)
	}
	return nil
}

//*** validators

// Validate gets run every time you call a "pop.Validate*" method.
func (r *Resource) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: r.Provider, Name: "Provider"},
		&validators.StringIsPresent{Field: r.Type, Name: "Type"},
		&validators.StringIsPresent{Field: r.Name, Name: "Name"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
func (r *Resource) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
func (r *Resource) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
