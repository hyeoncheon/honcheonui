package actions

import (
	"net/http"
	"strconv"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v5"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/hyeoncheon/honcheonui/models"
	"github.com/hyeoncheon/honcheonui/plugins"
	"github.com/hyeoncheon/honcheonui/workers"
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
	tx.Load(providers, "Member")

	c.Set("pagination", q.Paginator)
	return c.Render(200, r.Auto(c, providers))
}

// Create adds a Provider to the DB.
func (v ProvidersResource) Create(c buffalo.Context) error {
	provider := &models.Provider{}
	if err := c.Bind(provider); err != nil {
		return errors.WithStack(err)
	}

	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	plugin, err := plugins.GetPlugin(provider.Provider, "provider")
	if err != nil {
		return errors.WithStack(err)
	}
	uid, aid, err := plugin.CheckAccount(provider.User, provider.Pass)
	if err != nil {
		return c.Render(http.StatusUnprocessableEntity, r.String("plugin error: %v", err))
	}

	c.Logger().Debugf("check account: %v %v %v", uid, aid, err)
	provider.GroupID = strconv.Itoa(aid)
	provider.UserID = strconv.Itoa(uid)
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

// Sync gets resources from provider via plugin API
func (v ProvidersResource) Sync(c buffalo.Context) error {
	args := map[string]interface{}{
		"provider_id": c.Param("provider_id"),
	}
	if err := workers.Run(workers.WorkerResourceSync, args); err != nil {
		c.Flash().Add("danger", t(c, "Could.not.start.background.sync"))
	} else {
		c.Flash().Add("success", t(c, "Resources.will.be.synced.in.background"))
	}

	return c.Redirect(http.StatusSeeOther, "/settings")
}
