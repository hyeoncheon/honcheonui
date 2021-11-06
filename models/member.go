package models

import (
	"strconv"
	"time"

	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
)

// Member is a model for storing information of service user
type Member struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	Email     string    `json:"email" db:"email"`
	Name      string    `json:"name" db:"-"`
	Avatar    string    `json:"avatar" db:"-"`
	Roles     []string  `json:"role" db:"-"`
	Providers Providers `has_many:"providers" order_by:"provider"`
}

// String returns email address for the member
func (m Member) String() string {
	return m.Email
}

// Members is array of members
type Members []Member

// String returns number of members as string
func (m Members) String() string {
	return strconv.Itoa(len(m))
}

//*** relational operations and queries

// GroupTags gets and returns all accessible tags
func (m *Member) GroupTags() *Tags {
	tags := &Tags{}

	// get all tags for same group_id not just mine.
	// NOTE: I am not sure which is performing better even though plan for
	// join is shorter. Anyway, buffalo/pop's query builder support join :-)
	query := DB.Q().
		Join("resources_tags", "resources_tags.tag_id = tags.id").
		Join("resources", "resources.id = resources_tags.resource_id").
		Join("providers", "providers.group_id = resources.group_id").
		Where("providers.member_id = ?", m.ID).
		GroupBy("tags.id").Order("tags.name")
	/*
		query := DB.RawQuery(`SELECT tags.id, tags.name FROM tags WHERE id IN (
				SELECT tag_id FROM resources_tags WHERE resource_id IN (
					SELECT id FROM resources WHERE group_id IN (
						SELECT group_id FROM providers WHERE member_id = ?
					)
				)
			) ORDER BY tags.name`, m.ID)
	*/
	err := query.All(tags)
	if err != nil {
		mlogger.Errorf("get all tags error: %v", err)
	}
	return tags
}

//*** validators

// Validate gets run every time you call a "pop.Validate*" method.
func (m *Member) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: m.Email, Name: "Email"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
func (m *Member) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
func (m *Member) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
