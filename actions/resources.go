package actions

import (
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v5"
	"github.com/pkg/errors"

	"github.com/hyeoncheon/honcheonui/models"
)

// ResourcesResource is the resource for the Resource model
type ResourcesResource struct {
	buffalo.Resource
}

// List gets all Resources.
func (v ResourcesResource) List(c buffalo.Context) error {
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	resources := &models.Resources{}
	q := tx.PaginateFromParams(c.Params())
	if err := q.All(resources); err != nil {
		return errors.WithStack(err)
	}

	c.Set("pagination", q.Paginator)
	return c.Render(200, r.Auto(c, resources))
}

// Show gets the data for one Resource.
func (v ResourcesResource) Show(c buffalo.Context) error {
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	resource := &models.Resource{}
	if err := tx.Eager().Find(resource, c.Param("resource_id")); err != nil {
		return c.Error(404, err)
	}

	tx.Load(&resource.Providers, "Member")
	c.Set("services", resource.Services())
	return c.Render(200, r.Auto(c, resource))
}

// Update changes a Resource in the DB.
func (v ResourcesResource) Update(c buffalo.Context) error {
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	resource := &models.Resource{}
	if err := tx.Find(resource, c.Param("resource_id")); err != nil {
		return c.Error(404, err)
	}

	// TODO: sync via plugin

	verrs, err := tx.ValidateAndUpdate(resource)
	if err != nil {
		return errors.WithStack(err)
	}

	if verrs.HasAny() {
		c.Set("errors", verrs)
		return c.Render(422, r.Auto(c, resource))
	}

	c.Flash().Add("success", "Resource was updated successfully")
	return c.Render(200, r.Auto(c, resource))
}

// Destroy deletes a Resource from the DB.
func (v ResourcesResource) Destroy(c buffalo.Context) error {
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	resource := &models.Resource{}
	if err := tx.Find(resource, c.Param("resource_id")); err != nil {
		return c.Error(404, err)
	}

	if err := tx.Destroy(resource); err != nil {
		return errors.WithStack(err)
	}

	c.Flash().Add("success", "Resource was destroyed successfully")
	return c.Render(200, r.Auto(c, resource))
}
