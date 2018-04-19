package models

import (
	"errors"
	"strings"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
)

// Service is a structure for service
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

// ServicesTags is a link map of tags for services
type ServicesTags struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	ServiceID uuid.UUID `json:"service_id" db:"service_id"`
	TagID     uuid.UUID `json:"tag_id" db:"tag_id"`
}

// String returns name of the service
func (s Service) String() string {
	return s.Name
}

// Services is an array of services
type Services []Service

//*** relational operations and queries

// LinkTags make link maps for the service
func (s *Service) LinkTags(tagIDs []string) error {
	for _, e := range tagIDs {
		id, err := uuid.FromString(e)
		if err != nil {
			logger.Errorf("invalid parameter! non UUID tag ID: %v", e)
			return errors.New("invalid request")
		}
		if err := DB.Save(&ServicesTags{ServiceID: s.ID, TagID: id}); err != nil {
			if strings.Contains(err.Error(), "Duplicate") { // mysql case
				logger.Warnf("duplicated tag map %v to %v: %v", s.ID, id, err)
			} else {
				logger.Errorf("could not make tag map %v to %v: %v", s.ID, id, err)
			}
		}
	}
	return nil
}

// TaggedResources gets and returns all accessible tags
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
	if s.MatchAll {
		//// `IN` bug fixed https://github.com/gobuffalo/pop/issues/65
		////query = query.Having(fmt.Sprintf("count(name) = %v", len(IDs)))
		query = query.Having("count(name) = ?", len(IDs))
	}
	if err := query.All(resources); err != nil {
		logger.Errorf("could not get resources. error: %v", err)
	}

	return resources
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
