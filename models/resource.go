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
}

// ResourcesTags is structure for mapping resources to tags
type ResourcesTags struct {
	ID         uuid.UUID `db:"id"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
	ResourceID uuid.UUID `db:"resource_id"`
	TagID      uuid.UUID `db:"tag_id"`
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

// Resources is array of resources
type Resources []Resource

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

// LinkTag links the resource to tag with given name.
// if tag with the name does not exist, it create before linking.
func (r *Resource) LinkTag(name string) error {
	tag := &Tag{}
	name = strings.TrimSpace(name)
	err := DB.Where("name = ?", name).First(tag)
	if err != nil {
		// if cannot found, it returns error.
		tag.Name = name
		verrs, err := DB.ValidateAndCreate(tag)
		if verrs.HasAny() {
			// TODO logging
		}
		if err != nil {
			// TODO logging
		}
	}
	resourcesTags := &ResourcesTags{
		ResourceID: r.ID,
		TagID:      tag.ID,
	}
	if err := DB.Save(resourcesTags); err != nil {
		// TODO logging
		return err
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
