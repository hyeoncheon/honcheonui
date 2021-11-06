package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
)

// Attribute is structure for attribute of user resource
type Attribute struct {
	ID         uuid.UUID `json:"id" db:"id"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
	Name       string    `json:"name" db:"name"`
	Value      string    `json:"value" db:"value"`
	ResourceID uuid.UUID `json:"resource_id" db:"resource_id"`
	Resource   Resource  `belongs_to:"resource"`
}

// String returns name:value formmatted string
func (a Attribute) String() string {
	return a.Name + ":" + a.Value
}

// JSON returns json formatted string
func (a Attribute) JSON() string {
	ja, _ := json.Marshal(a)
	return string(ja)
}

// Attributes is an array of attributes
type Attributes []Attribute

//*** validators

// Validate gets run every time you call a "pop.Validate*" method.
func (a *Attribute) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.UUIDIsPresent{Field: a.ResourceID, Name: "ResourceID"},
		&validators.StringIsPresent{Field: a.Name, Name: "Name"},
		&validators.StringIsPresent{Field: a.Value, Name: "Value"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
func (a *Attribute) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
func (a *Attribute) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
