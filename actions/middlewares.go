package actions

import (
	"net/http"

	"github.com/gobuffalo/buffalo"
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
		memberID := c.Session().Get("member_id")
		if memberID != nil {
			c.Set("member_id", memberID)
		}
		return next(c)
	}
}
