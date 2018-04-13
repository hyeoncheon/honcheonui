package models

import (
	"strconv"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
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
