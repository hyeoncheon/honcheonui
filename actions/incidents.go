package actions

import (
	"net/http"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v5"
	"github.com/pkg/errors"

	"github.com/hyeoncheon/honcheonui/models"
)

// IncidentsResource is the resource for the Incident model
type IncidentsResource struct {
	buffalo.Resource
}

// Show gets the data for one Incident.
func (v IncidentsResource) Show(c buffalo.Context) error {
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	incident := &models.Incident{}
	if err := tx.Eager().Find(incident, c.Param("incident_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	return c.Render(http.StatusOK, r.Auto(c, incident))
}
