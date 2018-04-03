package actions

import (
	"fmt"
	"os"

	"github.com/gobuffalo/buffalo"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/cloudfoundry"
)

func init() {
	gothic.Store = App().SessionStore

	uartProvider := cloudfoundry.New(
		os.Getenv("UART_URL"),
		os.Getenv("UART_KEY"),
		os.Getenv("UART_SECRET"),
		fmt.Sprintf("%s%s", os.Getenv("HCU_URL"), "/auth/uart/callback"),
		"profile")
	uartProvider.SetName("uart")

	goth.UseProviders(
		uartProvider,
	)
}

// AuthCallback is universal callback handler for goth authorization
func AuthCallback(c buffalo.Context) error {
	user, err := gothic.CompleteUserAuth(c.Response(), c.Request())
	if err != nil {
		return c.Error(401, err)
	}
	// Do something with the user, maybe register them/sign them in
	return c.Render(200, r.JSON(user))
}
