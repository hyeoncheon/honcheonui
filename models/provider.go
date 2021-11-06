package models

import (
	"errors"
	"time"

	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
	"github.com/hyeoncheon/honcheonui/utils"
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
	Member    Member    `belongs_to:"member"`
	Resources Resources `many_to_many:"providers_resources"`
}

// ProvidersResources is a map between provider and resource
type ProvidersResources struct {
	ID         uuid.UUID `json:"id" db:"id"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
	ProviderID uuid.UUID `json:"provider_id" db:"provider_id"`
	ResourceID uuid.UUID `json:"resource_id" db:"resource_id"`
}

// String currently returns provider name and username
func (p Provider) String() string {
	return p.Provider + "/" + p.User
}

// Owner returns owner of the provider entity
// DEPRECATED: now buffalo support assotiation more easily with Eager() and Load()
func (p Provider) Owner() *Member {
	mlogger.Warn("using obsoleted method")
	member := &Member{}
	if err := DB.Find(member, p.MemberID); err != nil {
		return nil
	}
	return member
}

// Providers is an array of providers
type Providers []Provider

// ProvidersResourcesMaps is an array of providers resources map
type ProvidersResourcesMaps []ProvidersResources

//*** relationship

// LinkResources makes a link map of provider and resources
func (p *Provider) LinkResources(oids []uuid.UUID) error {
	ids, err := utils.ToInterface(oids)
	if err != nil {
		return errors.New("could not convert argument to interface slice")
	}
	mlogger.Debugf("link resources requested: %v", ids)
	hasError := false

	// get existing mappings and remove them from given list
	maps := &ProvidersResourcesMaps{}
	if err := DB.Where("provider_id = ?", p.ID).All(maps); err != nil {
		mlogger.Errorf("database selection failed! error: %v", err)
	}
	for _, m := range *maps {
		mlogger.Debugf("check %v is on the list...", m.ResourceID)
		if utils.Has(ids, m.ResourceID) {
			ids = utils.Remove(ids, m.ResourceID)
			mlogger.Debugf("removed existing resource from list: %v", ids)
		} else {
			mlogger.Debugf("removing %v from map. no longer exists", m.ResourceID)
			if err := DB.Destroy(&m); err != nil {
				mlogger.Errorf("could not remove resource %v from the map", m)
				hasError = true
			}
		}
	}

	mlogger.Debugf("adding new resources: %v", ids)
	for _, id := range ids {
		mlogger.Debugf("creating map for %v on %v", id, p)
		prmap := &ProvidersResources{
			ProviderID: p.ID,
			ResourceID: id.(uuid.UUID), //! check me
		}
		if err := DB.Save(prmap); err != nil {
			mlogger.Errorf("could not save provider-resource map %v: %v", prmap, err)
			hasError = true
		}
	}
	if hasError {
		return errors.New("linking done with error(s)")
	}
	return nil
}

//*** validators

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
