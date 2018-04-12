package actions

import (
	"io/ioutil"
	"os"
	"plugin"

	"github.com/gobuffalo/buffalo"
	"github.com/pkg/errors"
)

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

func getPluginList(c buffalo.Context) []string {
	var plugins []string

	pluginHome := os.Getenv("HCU_HOME") + "/plugins"
	c.Logger().WithField("category", "plugin").
		Debugf("searching on %v", pluginHome)
	files, err := ioutil.ReadDir(pluginHome)
	if err != nil {
		// error
		return plugins
	}
	for _, f := range files {
		name := f.Name()
		if ".so" == name[len(name)-3:len(name)] {
			name = name[0 : len(name)-3]
			if err := checkPlugin(name); err == nil {
				plugins = append(plugins, name)
			}
		}
	}
	c.Logger().WithField("category", "plugin").
		Debugf("plugins found: %v", plugins)
	return plugins
}

func checkPlugin(name string) error {
	if _, err := getPlugin(name); err != nil {
		return err
	}
	return nil
}

func getPlugin(name string) (Provider, error) {
	name = name + ".so"
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
