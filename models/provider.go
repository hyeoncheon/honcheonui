package models

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
)

// Provider structure
type Provider struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	MemberID  uuid.UUID `json:"member_id" db:"member_id"`
	Provider  string    `json:"provider" db:"provider"`
	User      string    `json:"user" db:"user"`
	Pass      string    `json:"pass" db:"pass"`
	GroupID   string    `json:"group_id" db:"group_id"`
	UserID    string    `json:"user_id" db:"user_id"`
}

// String currently returns provider name and username
func (p Provider) String() string {
	return p.Provider + "/" + p.User
}

// Member returns owner of the provider entity
func (p Provider) Member() *Member {
	member := &Member{}
	if err := DB.Find(member, p.MemberID); err != nil {
		return nil
	}
	return member
}

// Providers is array of providers
type Providers []Provider

// Validate gets run every time you call a "pop.Validate*" method.
func (p *Provider) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: p.Provider, Name: "Provider"},
		&validators.StringIsPresent{Field: p.User, Name: "User"},
		&validators.StringIsPresent{Field: p.Pass, Name: "Pass"},
		&validators.StringIsPresent{Field: p.GroupID, Name: "GroupID"},
		&validators.StringIsPresent{Field: p.UserID, Name: "UserID"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
func (p *Provider) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
func (p *Provider) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
