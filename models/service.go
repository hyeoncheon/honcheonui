package models

import (
	"errors"
	"time"

	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
)

// Service is a structure for service.
// Service is user's perspective and it has many resources indirectly via tags
// and matching rule.
type Service struct {
	ID          uuid.UUID `json:"id" db:"id"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	MemberID    uuid.UUID `json:"member_id" db:"member_id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	MatchAll    bool      `json:"match_all" db:"match_all"`
	Member      Member    `belongs_to:"members"`
	Resources   Resources `many_to_many:"services_resources"`
	Tags        Tags      `many_to_many:"services_tags"`
}

// ServicesTags is a link map of tags for services.
type ServicesTags struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	ServiceID uuid.UUID `json:"service_id" db:"service_id"`
	TagID     uuid.UUID `json:"tag_id" db:"tag_id"`
}

// String returns name of the service.
func (s Service) String() string {
	return s.Name
}

// Services is an array of services.
type Services []Service

//*** special functions

// HasResource returns true if resource is associated with the service.
// This relationship is indirect.
func (s *Service) HasResource(r *Resource) bool {
	count, err := DB.Where("service_id = ?", s.ID).Count(&ServicesTags{})
	if err != nil {
		mlogger.Errorf("database error: %v", err)
		return false
	}

	query := DB.Q().
		Join("resources_tags", "resources_tags.resource_id = resources.id").
		Join("tags", "tags.id = resources_tags.tag_id").
		Join("services_tags", "services_tags.tag_id = tags.id").
		Where("services_tags.service_id = ?", s.ID).
		Where("resources.id = ?", r.ID).
		GroupBy("resources.id")
	if s.MatchAll {
		query = query.Having("count(resources.name) = ?", count)
	}
	resources := &Resources{}
	if err := query.All(resources); err != nil {
		mlogger.Errorf("could not get resources. error: %v", err)
	}
	mlogger.Debugf("%v has %v tags and %v is matched", s, count, *resources)

	return len(*resources) == 1
}

//*** relational operations and queries

// LinkTags make link maps for the service.
func (s *Service) LinkTags(tagIDs []string) error {
	hasError := false
	for _, e := range tagIDs {
		id, err := uuid.FromString(e)
		if err != nil {
			mlogger.Errorf("invalid parameter! non UUID tag ID: %v", e)
			return errors.New("invalid request")
		}
		//? should I add tag entry checking?
		if err := TrySave(&ServicesTags{ServiceID: s.ID, TagID: id}); err != nil {
			hasError = true
		}
	}
	if hasError {
		return errors.New("at least one tag is not linked")
	}
	return nil
}

// TaggedResources gets and returns all accessible tags.
// Service has associated resources but the relationship is indirect via tags.
func (s *Service) TaggedResources() *Resources {
	resources := &Resources{}
	if len(s.Tags) < 1 {
		return resources
	}

	var IDs []interface{}
	for _, t := range s.Tags {
		IDs = append(IDs, t.ID)
	}

	query := DB.Q().
		Join("resources_tags", "resources_tags.resource_id = resources.id").
		Where("resources_tags.tag_id in (?)", IDs...).
		GroupBy("resources.id")
	// if matching rule is MatchAll, resources which has all tags will be
	// selected. Otherwise resources with any tags are selected.
	if s.MatchAll {
		//// `IN` bug fixed https://github.com/gobuffalo/pop/issues/65
		////query = query.Having(fmt.Sprintf("count(name) = %v", len(IDs)))
		query = query.Having("count(name) = ?", len(IDs))
	}
	if err := query.All(resources); err != nil {
		mlogger.Errorf("could not get resources. error: %v", err)
	}

	return resources
}

// Incidents returns incidents associated with resources of the service.
func (s *Service) Incidents() *Incidents {
	incidents := &Incidents{}
	if len(s.Tags) < 1 {
		return incidents
	}

	var IDs []interface{}
	for _, t := range s.Tags {
		IDs = append(IDs, t.ID)
	}

	// TODO: select via events?
	query := DB.Q().
		Join("incidents_resources", "incidents_resources.incident_id = incidents.id").
		Join("resources", "resources.id = incidents_resources.resource_id").
		Join("resources_tags", "resources_tags.resource_id = resources.id").
		Where("resources_tags.tag_id in (?)", IDs...)
	if err := query.All(incidents); err != nil {
		mlogger.Errorf("could not get incidents: %v", err)
	}
	return incidents
}

//*** validators

// Validate gets run every time you call a "pop.Validate*" method.
func (s *Service) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: s.Name, Name: "Name"},
		&validators.StringIsPresent{Field: s.Description, Name: "Description"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
func (s *Service) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
func (s *Service) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
