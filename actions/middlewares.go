package actions

import (
	"net/http"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/uuid"
)

// AuthorizeHandler protect all application pages from unauthorized accesses.
func AuthorizeHandler(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		memberID := c.Session().Get("member_id")
		if memberID == nil {
			c.Logger().Warn("unauthorized access to ", c.Request().RequestURI)
			c.Flash().Add("danger", t(c, "login.required"))
			return c.Redirect(http.StatusTemporaryRedirect, "/login")
		}
		return next(c)
	}
}

func contextHandler(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		c.Set("TIME_FORMAT", "2006-01-02T15:04:05Z07:00")
		memberID := c.Session().Get("member_id")
		if memberID != nil {
			if id, ok := memberID.(uuid.UUID); ok {
				c.Set("member_id", id.String())
			}
			c.Set("member_mail", c.Session().Get("member_mail"))
			c.Set("member_name", c.Session().Get("member_name"))
			c.Set("member_icon", c.Session().Get("member_icon"))
			c.Set("member_roles", c.Session().Get("member_roles"))
		}
		return next(c)
	}
}
