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

// Resource is
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
}

// ResourcesTags is structure for mapping resources to tags
type ResourcesTags struct {
	ID         uuid.UUID `db:"id"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
	ResourceID uuid.UUID `db:"resource_id"`
	TagID      uuid.UUID `db:"tag_id"`
}

// ResourcesUsers is structure for mapping resources to users
type ResourcesUsers struct {
	ID         uuid.UUID `db:"id"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
	ResourceID uuid.UUID `db:"resource_id"`
	UserID     string    `db:"user_id"`
	Provider   Provider  `belongs_to:"providers" fk_id:"user_id"`
}

// String represents its name and provider
func (r Resource) String() string {
	return r.Name + "@" + r.Provider
}

// JSON returns json marshalled object
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

// AddAttribute creates and saves attribute of the resource
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

// LinkTags makes a link map of resource and user
func (r *Resource) LinkTags(ts []string) error {
	tags, err := utils.ToInterface(utils.Cleaner(ts))
	if err != nil {
		return errors.New("could not convert argument to interface slice")
	}
	logger.Debugf("link tags requested: %v", tags)
	hasError := false

	// get existing tag map and remove them from given list
	maps := &RTMaps{}
	if err := DB.Where("resource_id = ?", r.ID).All(maps); err != nil {
		logger.Errorf("database selection failed! error: %v", err)
	}
	for _, m := range *maps {
		tag := &Tag{}
		err := DB.Find(tag, m.TagID)
		if err != nil { // in case of broken link map
			err := DB.Destroy(&m)
			if err != nil {
				logger.Errorf("found broken map but could not delete: id:%v", m.ID)
			} else {
				logger.Warnf("found broken map and deleted: id:%v", m.ID)
			}
		}

		if utils.Has(tags, tag.Name) {
			tags = utils.Remove(tags, tag.Name)
		} else {
			logger.Debugf("removing %v from map. no longer exists", tag.Name)
			if err := DB.Destroy(&m); err != nil {
				logger.Errorf("could not remove  %v from the map", m)
				hasError = true
			}
		}
	}

	logger.Debugf("adding new tags...: %v", tags)
	for _, u := range tags {
		name := strings.TrimSpace(u.(string)) //! check me

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

// LinkUsers makes a link map of resource and user
func (r *Resource) LinkUsers(us []string) error {
	users, err := utils.ToInterface(us)
	if err != nil {
		return errors.New("could not convert argument to interface slice")
	}
	logger.Debugf("link users requested: %v", users)
	hasError := false

	// get existing user map and remove them from given list
	maps := &RUMaps{}
	if err := DB.Where("resource_id = ?", r.ID).All(maps); err != nil {
		logger.Errorf("database selection failed! error: %v", err)
	}
	for _, m := range *maps {
		if utils.Has(users, m.UserID) {
			users = utils.Remove(users, m.UserID)
		} else {
			logger.Debugf("removing %v from map. no longer exists", m.UserID)
			if err := DB.Destroy(&m); err != nil {
				logger.Errorf("could not remove user %v from the map", m)
				hasError = true
			}
		}
	}

	logger.Debugf("adding new users...: %v", users)
	for _, u := range users {
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
