package actions

import (
	"net/http"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/pkg/errors"

	"github.com/hyeoncheon/honcheonui/models"
)

// ProvidersResource is the resource for the Provider model
type ProvidersResource struct {
	buffalo.Resource
}

// List gets all Providers.
//! MANAGER ONLY, currently not protected
func (v ProvidersResource) List(c buffalo.Context) error {
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	providers := &models.Providers{}
	q := tx.PaginateFromParams(c.Params())
	if err := q.All(providers); err != nil {
		return errors.WithStack(err)
	}

	c.Set("pagination", q.Paginator)
	return c.Render(200, r.Auto(c, providers))
}

// Create adds a Provider to the DB.
//! currently not implemented. just for testing now.
func (v ProvidersResource) Create(c buffalo.Context) error {
	provider := &models.Provider{}
	if err := c.Bind(provider); err != nil {
		return errors.WithStack(err)
	}

	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	provider.GroupID = "1234"
	provider.UserID = "123456"
	provider.MemberID = c.Session().Get("member_id").(uuid.UUID)
	verrs, err := tx.ValidateAndCreate(provider)
	if err != nil {
		return errors.WithStack(err)
	}
	if verrs.HasAny() {
		c.Set("errors", verrs)
		return c.Render(http.StatusUnprocessableEntity, r.String("value error: %v", verrs))
	}

	return c.Render(http.StatusCreated, r.String("provider created"))
}

// Destroy deletes a Provider from the DB.
func (v ProvidersResource) Destroy(c buffalo.Context) error {
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	provider := &models.Provider{}
	if err := tx.Find(provider, c.Param("provider_id")); err != nil {
		return c.Error(404, err)
	}

	if err := tx.Destroy(provider); err != nil {
		return errors.WithStack(err)
	}

	c.Flash().Add("success", "Provider was destroyed successfully")
	return c.Redirect(http.StatusSeeOther, "/settings")
}
