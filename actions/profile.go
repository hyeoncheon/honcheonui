package actions

import (
	"os"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/pkg/errors"

	"github.com/hyeoncheon/honcheonui/models"
)

// ProfileShow renders current member's profile page.
func ProfileShow(c buffalo.Context) error {
	return c.Render(200, r.HTML("profile/show.html"))
}

// ProfileSettings renders current member's settings page.
func ProfileSettings(c buffalo.Context) error {
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	currentMember := &models.Member{}
	err := tx.Find(currentMember, c.Session().Get("member_id"))
	if err != nil {
		return errors.WithStack(err)
	}

	providers := &models.Providers{}
	if err := tx.BelongsTo(currentMember).All(providers); err != nil {
		return errors.WithStack(err)
	}

	supportedProviders := map[string]string{
		"SoftLayer": "softlayer",
	}
	c.Set("providers", providers)
	c.Set("provider", &models.Provider{})
	c.Set("uart_url", os.Getenv("UART_URL"))
	c.Set("supported_providers", supportedProviders)
	return c.Render(200, r.HTML("profile/settings.html"))
}
