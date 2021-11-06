package models

import (
	"encoding/json"
	"log"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/logger"
	"github.com/gobuffalo/pop/v5"
)

// DB is a connection to your database to be used throughout your application.
var DB *pop.Connection

// default logger and security logger
var mlogger = logger.NewLogger("Debug").WithField("category", "models")
var slogger = mlogger.WithField("category", "security")

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
	mlogger = l.WithField("category", "models")
	slogger = mlogger.WithField("category", "security")
}

func inspect(desc string, data interface{}) {
	mlogger.Debugf("inspect: %v: %v", desc, JSON(data))
}

// JSON returns json formatted object
func JSON(d interface{}) string {
	ba, _ := json.Marshal(d)
	return string(ba)
}
