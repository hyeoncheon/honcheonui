package actions

import (
	"net/http"

	"github.com/gobuffalo/buffalo"
)

// HomeHandler is a default handler to serve up
// a home page.
func HomeHandler(c buffalo.Context) error {
	return c.Render(200, r.HTML("index.html"))
}

// LoginHandler renders login page
func LoginHandler(c buffalo.Context) error {
	return c.Render(http.StatusOK, r.HTML("login.html"))
}

// LogoutHandler clears all session information and redirect to root.
func LogoutHandler(c buffalo.Context) error {
	sess := c.Session()
	sess.Clear()
	c.Flash().Add("success", t(c, "you.have.been.successfully.logged.out"))
	return c.Redirect(http.StatusTemporaryRedirect, "/")
}
