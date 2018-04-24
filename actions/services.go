package actions

import (
	"net/http"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/pkg/errors"

	"github.com/hyeoncheon/honcheonui/models"
)

// ServicesResource is the resource for the Service model
type ServicesResource struct {
	buffalo.Resource
}

// List gets all Services.
func (v ServicesResource) List(c buffalo.Context) error {
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	services := &models.Services{}
	q := tx.Eager("Member").PaginateFromParams(c.Params())
	if err := q.All(services); err != nil {
		return errors.WithStack(err)
	}

	c.Set("pagination", q.Paginator)
	return c.Render(200, r.Auto(c, services))
}

// Show gets the data for one Service.
func (v ServicesResource) Show(c buffalo.Context) error {
	tx, service, err := setService(c)
	if err != nil {
		return err
	}
	tx.Load(service, "Member", "Tags")

	c.Set("incidents", service.Incidents())
	c.Set("resources", service.TaggedResources())
	c.Set("tags", effectiveMember(c).GroupTags())
	return c.Render(200, r.Auto(c, service))
}

// New renders the form for creating a new Service.
func (v ServicesResource) New(c buffalo.Context) error {
	return c.Render(200, r.Auto(c, &models.Service{}))
}

// Create adds a Service to the DB.
func (v ServicesResource) Create(c buffalo.Context) error {
	service := &models.Service{}
	if err := c.Bind(service); err != nil {
		return errors.WithStack(err)
	}
	service.MemberID = effectiveMember(c).ID

	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	verrs, err := tx.ValidateAndCreate(service)
	if err != nil {
		return errors.WithStack(err)
	}

	if verrs.HasAny() {
		c.Set("errors", verrs)
		return c.Render(422, r.Auto(c, service))
	}

	c.Flash().Add("success", "Service was created successfully")
	return c.Render(201, r.Auto(c, service))
}

// Edit renders a edit form for a Service.
func (v ServicesResource) Edit(c buffalo.Context) error {
	_, service, err := setService(c)
	if err != nil {
		return err
	}

	return c.Render(200, r.Auto(c, service))
}

// Update changes a Service in the DB.
func (v ServicesResource) Update(c buffalo.Context) error {
	tx, service, err := setService(c)
	if err != nil {
		return err
	}

	if err := c.Bind(service); err != nil {
		return errors.WithStack(err)
	}

	verrs, err := tx.ValidateAndUpdate(service)
	if err != nil {
		return errors.WithStack(err)
	}

	if verrs.HasAny() {
		c.Set("errors", verrs)
		return c.Render(422, r.Auto(c, service))
	}

	c.Flash().Add("success", "Service was updated successfully")
	return c.Render(200, r.Auto(c, service))
}

// Destroy deletes a Service from the DB.
func (v ServicesResource) Destroy(c buffalo.Context) error {
	tx, service, err := setService(c)
	if err != nil {
		return err
	}

	if err := tx.Destroy(service); err != nil {
		return errors.WithStack(err)
	}

	c.Flash().Add("success", "Service was destroyed successfully")
	return c.Render(200, r.Auto(c, service))
}

// AddTags make a tag maps for service
func (v ServicesResource) AddTags(c buffalo.Context) error {
	var tagIDs []string
	if err := c.Request().ParseForm(); err == nil {
		tagIDs = c.Request().Form["tag_id"]
	}
	c.Logger().Infof("AddTags: %v tags are requested", len(tagIDs))
	c.Logger().Debugf("----- %v", tagIDs)

	_, service, err := setService(c)
	if err != nil {
		return err
	}

	if err := service.LinkTags(tagIDs); err != nil {
		return c.Render(http.StatusUnprocessableEntity, r.String("tag linking error: %v", err))
	}
	return c.Render(http.StatusCreated, r.String("tags are saved"))
}

func setService(c buffalo.Context) (*pop.Connection, *models.Service, error) {
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return nil, nil, errors.WithStack(errors.New("no transaction found"))
	}

	service := &models.Service{}
	if err := tx.Find(service, c.Param("service_id")); err != nil {
		return nil, nil, c.Error(404, err)
	}
	return tx, service, nil
}
