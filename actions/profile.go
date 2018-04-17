package actions

import (
	"os"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
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
	err := tx.Eager().Find(currentMember, c.Session().Get("member_id"))
	if err != nil {
		return errors.WithStack(err)
	}
	tx.Load(&currentMember.Providers, "Member")

	supportedProviders := make(map[string]string)
	for _, p := range getPluginList(c) {
		supportedProviders[p] = p
	}
	c.Set("providers", currentMember.Providers)
	c.Set("provider", &models.Provider{}) // for modal form
	c.Set("uart_url", os.Getenv("UART_URL"))
	c.Set("supported_providers", supportedProviders)
	return c.Render(200, r.HTML("profile/settings.html"))
}

func effectiveMember(c buffalo.Context) *models.Member {
	dummy := &models.Member{}
	if id, ok := c.Value("member_id").(uuid.UUID); ok {
		dummy.ID = id
	}
	return dummy
}
