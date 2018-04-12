package models

import (
	"log"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/pop"
)

// DB is a connection to your database to be used throughout your application.
var DB *pop.Connection

// default logger and security logger
var logger = buffalo.NewLogger("Debug").WithField("category", "models")
var slogger = logger.WithField("category", "security")

func init() {
	var err error
	env := envy.Get("GO_ENV", "development")
	DB, err = pop.Connect(env)
	if err != nil {
		log.Fatal(err)
	}
	pop.Debug = env == "development"
}

// SetLogger sets logger and slogger with external logger
func SetLogger(l buffalo.Logger) {
	logger = l.WithField("category", "models")
	slogger = logger.WithField("category", "security")
}
