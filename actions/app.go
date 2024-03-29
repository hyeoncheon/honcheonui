package actions

import (
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo-pop/v2/pop/popmw"
	"github.com/gobuffalo/envy"
	csrf "github.com/gobuffalo/mw-csrf"
	i18n "github.com/gobuffalo/mw-i18n"
	paramlogger "github.com/gobuffalo/mw-paramlogger"
	"github.com/gobuffalo/packr/v2"
	"github.com/markbates/goth/gothic"

	"github.com/hyeoncheon/honcheonui/models"
	"github.com/hyeoncheon/honcheonui/workers"
)

// global variables
var (
	ENV = envy.Get("GO_ENV", "development")
	app *buffalo.App
	T   *i18n.Translator
)

// App is where all routes and middleware for buffalo
// should be defined. This is the nerve center of your
// application.
func App() *buffalo.App {
	if app == nil {
		app = buffalo.New(buffalo.Options{
			Env:         ENV,
			SessionName: "_honcheonui_session",
			// TODO: add secure session store. should it be redis?
		})

		if err := workers.InitWorkers(app); err != nil {
			app.Logger.Errorf("error while initializing workers: %v", err)
		}

		if ENV == "development" {
			app.Use(paramlogger.ParameterLogger)
		}

		// https://www.owasp.org/index.php/Cross-Site_Request_Forgery_(CSRF)
		app.Use(csrf.New)

		// Wraps each request in a transaction.
		//  c.Value("tx").(*pop.PopTransaction)
		app.Use(popmw.Transaction(models.DB))
		models.SetLogger(app.Logger)

		// Setup and use translations:
		var err error
		if T, err = i18n.New(packr.NewBox("../locales"), "en-US"); err != nil {
			app.Stop(err)
		}
		app.Use(T.Middleware())

		app.GET("/", HomeHandler)
		app.GET("/login", LoginHandler)
		app.GET("/logout", LogoutHandler)

		// authorization with uart
		auth := app.Group("/auth")
		auth.GET("/{provider}", buffalo.WrapHandlerFunc(gothic.BeginAuthHandler))
		auth.GET("/{provider}/callback", AuthCallback)

		// protect resources and set context for the session
		app.Use(AuthorizeHandler)
		app.Middleware.Skip(AuthorizeHandler, LoginHandler)
		app.Use(contextHandler)

		app.Resource("/members", MembersResource{})
		app.GET("/profile", ProfileShow)
		app.GET("/settings", ProfileSettings)
		app.GET("/providers", ProvidersResource{}.List)
		app.POST("/providers", ProvidersResource{}.Create)
		app.DELETE("/providers/{provider_id}", ProvidersResource{}.Destroy)
		app.GET("/providers/{provider_id}/sync", ProvidersResource{}.Sync)
		app.GET("/resources", ResourcesResource{}.List)
		app.GET("/resources/{resource_id}", ResourcesResource{}.Show)
		app.GET("/resources/{resource_id}/sync", ResourcesResource{}.Update)
		app.DELETE("/resources/{resource_id}", ResourcesResource{}.Destroy)
		app.Resource("/services", ServicesResource{})
		app.POST("/services/{service_id}/add_tags", ServicesResource{}.AddTags)
		app.GET("/incidents/{incident_id}", IncidentsResource{}.Show)

		admin := app.Group("/admin")
		admin.GET("/", AdminHandler)
		admin.GET("/sync/notification", AdminSyncNotification)

		app.ServeFiles("/", assetsBox) // serve files from the public directory
	}

	return app
}
