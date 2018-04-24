package models

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
)

// Incident is a struct for most atomic incident and event records.
type Incident struct {
	ID         uuid.UUID `json:"id" db:"id"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
	Provider   string    `json:"provider" db:"provider"`
	Type       string    `json:"type" db:"type"`
	OriginalID string    `json:"original_id" db:"original_id"`
	GroupID    string    `json:"group_id" db:"group_id"`
	UserID     string    `json:"user_id" db:"user_id"`
	Title      string    `json:"title" db:"title"`
	Content    string    `json:"content" db:"content"`
	Category   string    `json:"category" db:"category"`
	IssuedBy   string    `json:"issued_by" db:"issued_by"`
	IsOpen     bool      `json:"is_open" db:"is_open"`
	IssuedAt   time.Time `json:"issued_at" db:"issued_at"`
	ModifiedAt time.Time `json:"modified_at" db:"modified_at"`
	Resources  Resources `many_to_many:"incidents_resources"`
}

// IncidentsResources is structure for mapping incidents to resources
type IncidentsResources struct {
	ID         uuid.UUID `db:"id"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
	IncidentID uuid.UUID `db:"incident_id"`
	ResourceID uuid.UUID `db:"resource_id"`
}

// IncidentsUsers is structure for mapping incidents to users
type IncidentsUsers struct {
	ID         uuid.UUID `db:"id"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
	IncidentID uuid.UUID `db:"incident_id"`
	UserID     string    `db:"user_id"`
}

// String returns json marshalled string of incident
func (i Incident) String() string {
	ji, _ := json.Marshal(i)
	return string(ji)
}

// Incidents is an array of incidents
type Incidents []Incident

//*** common database handling

// Save just save given incident record
func (i *Incident) Save() error {
	return TrySave(i)
}

// LinkResourcesByOrigIDs makes a link map for incident to resources.
// If resource is not exist on database, just ignore it.
func (i *Incident) LinkResourcesByOrigIDs(IDs ...string) error {
	success := 0
	for _, id := range IDs {
		resource := &Resource{}
		if err := DB.Where("original_id = ?", id).First(resource); err != nil {
			if strings.Contains(err.Error(), "no rows") {
				logger.Warnf("no resource with original_id %v", id)
			} else {
				logger.Errorf("database error: %v", err)
			}
			continue
		}
		ir := &IncidentsResources{
			ResourceID: resource.ID,
			IncidentID: i.ID,
		}
		if err := TrySave(ir); err == nil {
			success++
		}
	}
	if success < len(IDs) {
		logger.Warnf("only %v entries saved out of %v requests", success, len(IDs))
	}
	return nil
}

// LinkUsers makes a link map for incident to users.
// It just make a link without checking user model since there is no such model currently.
func (i *Incident) LinkUsers(IDs ...string) error {
	for _, id := range IDs {
		if err := TrySave(&IncidentsUsers{
			IncidentID: i.ID,
			UserID:     id,
		}); err != nil {
			logger.Errorf("could not link with %v", id)
		}
	}
	return nil
}

//*** validators

// Validate gets run every time you call a "pop.Validate*" method.
func (i *Incident) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: i.Provider, Name: "Provider"},
		&validators.StringIsPresent{Field: i.Type, Name: "Type"},
		&validators.StringIsPresent{Field: i.OriginalID, Name: "OriginalID"},
		&validators.StringIsPresent{Field: i.GroupID, Name: "GroupID"},
		&validators.StringIsPresent{Field: i.UserID, Name: "UserID"},
		&validators.StringIsPresent{Field: i.Title, Name: "Title"},
		&validators.StringIsPresent{Field: i.Content, Name: "Content"},
		&validators.StringIsPresent{Field: i.Category, Name: "Category"},
		&validators.StringIsPresent{Field: i.IssuedBy, Name: "IssuedBy"},
		&validators.TimeIsPresent{Field: i.IssuedAt, Name: "IssuedAt"},
		&validators.TimeIsPresent{Field: i.ModifiedAt, Name: "ModifiedAt"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
func (i *Incident) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
func (i *Incident) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
