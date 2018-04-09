package actions

import (
	"os"

	"github.com/gobuffalo/buffalo"
)

// ProfileShow renders current member's profile page.
func ProfileShow(c buffalo.Context) error {
	return c.Render(200, r.HTML("profile/show.html"))
}

// ProfileSettings renders current member's settings page.
func ProfileSettings(c buffalo.Context) error {
	c.Set("uart_url", os.Getenv("UART_URL"))
	return c.Render(200, r.HTML("profile/settings.html"))
}
