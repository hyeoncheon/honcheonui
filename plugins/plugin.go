package plugins

import (
	"errors"
	"io/ioutil"
	"os"
	"plugin"
	"strings"
	"time"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/uuid"
)

// HoncheonuiResource structure
type HoncheonuiResource struct {
	Provider           string
	Type               string
	OriginalID         string
	UUID               uuid.UUID
	Name               string
	Notes              string
	GroupID            string
	ResourceCreatedAt  time.Time
	ResourceModifiedAt time.Time
	IPAddress          string
	Location           string
	IsConn             bool
	IsOn               bool
	Attributes         map[string]string
	IntegerAttributes  map[string]int
	Tags               []string
	UserIDs            []string
	Raw                interface{}
}

//*** provider plugin

// TODO: plugin caching
// TODO: reload if updated
// TODO: plugin types

// Provider interface
type Provider interface {
	Init() error
	CheckAccount(user, pass string) (int, int, error)
	GetResources(user, pass string) ([]interface{}, error)
	GetStatuses(user, pass string) ([]interface{}, error)
}

// GetPlugin returns provider via plugin
func GetPlugin(name, class string) (Provider, error) {
	// TODO: am I need to seperate provider plugin and messaging plugins?
	name = class + "-" + name + ".so"
	plug, err := plugin.Open(os.Getenv("HCU_HOME") + "/plugins/" + name)
	if err != nil {
		return nil, err
	}

	symbol, err := plug.Lookup("Provider")
	if err != nil {
		return nil, err
	}

	var provider Provider
	provider, ok := symbol.(Provider)
	if !ok {
		return nil, errors.New("invalid plugin")
	}
	// TODO: version check and logging
	return provider, nil
}

func checkPlugin(name, class string) error {
	if _, err := GetPlugin(name, class); err != nil {
		return err
	}
	return nil
}

// GetPluginList returns an array of plugin names
func GetPluginList(c buffalo.Context, class string) []string {
	var plugins []string

	pluginHome := os.Getenv("HCU_HOME") + "/plugins"
	c.Logger().WithField("category", "plugin").Debugf("searching on %v", pluginHome)
	files, err := ioutil.ReadDir(pluginHome)
	if err != nil {
		// error
		return plugins
	}
	for _, f := range files {
		name := f.Name()
		if strings.HasPrefix(name, class+"-") && strings.HasSuffix(name, ".so") {
			name = strings.TrimSuffix(strings.TrimPrefix(name, class+"-"), ".so")
			if err := checkPlugin(name, class); err == nil {
				plugins = append(plugins, name)
			}
		}
	}
	c.Logger().WithField("category", "plugin").Debugf("plugins found: %v", plugins)
	return plugins
}
