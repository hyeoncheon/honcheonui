package actions

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"

	"github.com/hyeoncheon/honcheonui/models"
	"github.com/hyeoncheon/honcheonui/plugins"
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

	plugin, err := getPlugin(provider.Provider)
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
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	provider := &models.Provider{}
	if err := tx.Find(provider, c.Param("provider_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	plugin, err := getPlugin(provider.Provider)
	if err != nil {
		return errors.WithStack(err)
	}
	resources, err := plugin.GetResources(provider.User, provider.Pass)
	if err != nil {
		return c.Render(http.StatusUnprocessableEntity, r.String("plugin error: %v", err))
	}

	for _, r := range resources {
		res := &models.Resource{}
		if jr, err := json.Marshal(r); err == nil {
			c.Logger().Debugf("------ json: %v", string(jr))
			re := &plugins.HoncheonuiResource{}
			if err := json.Unmarshal(jr, re); err != nil {
				c.Logger().Errorf("error: %v", err)
				c.Flash().Add("danger", t(c, "could.not.interprete.message"))
				return c.Redirect(http.StatusSeeOther, "/settings")
			}
			copier.Copy(res, re)
			// universally unique identifier, uuid is not perfectly uniq but almost.
			// but we can assume it is uniq anyway.
			// buffalo/pop uses uuid version 4 based on random number generator and
			// softlayer seems to use real random string as uuid. :-(
			if res.UUID != uuid.Nil {
				res.ID = res.UUID
			}
			if err := res.Save(); err != nil {
				c.Logger().Errorf("saving error: %v", err)
				c.Logger().Errorf("---- resource: %v", res.JSON())
				c.Flash().Add("warning", t(c, "could.not.save.resource")+res.String())
			}
			for k, v := range re.Attributes {
				res.AddAttribute(k, v)
			}
			for _, tag := range re.Tags {
				if err := res.LinkTag(tag); err != nil {
					c.Logger().Errorf("could not link to tag %v: %v", tag, err)
					c.Flash().Add("warning", t(c, "could.not.link.to.tag")+" "+tag)
				}
			}
			c.Logger().Debugf("------ re.UserIDs: %v", re.UserIDs)
			c.Logger().Debugf("------ re.IntegerAttributes: %v", re.IntegerAttributes)
		}
	}

	c.Flash().Add("success", t(c, "resources.was.synced.successfully"))
	return c.Redirect(http.StatusSeeOther, "/settings")
}
