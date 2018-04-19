package plugins

import (
	"errors"
	"io/ioutil"
	"os"
	"plugin"
	"strings"

	"github.com/gobuffalo/buffalo"
	spec "github.com/hyeoncheon/honcheonui-spec"
)

// TODO: plugin caching
// TODO: reload if updated
// TODO: plugin types

// GetPlugin returns provider via plugin
func GetPlugin(name, class string) (spec.Provider, error) {
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

	var provider spec.Provider
	provider, ok := symbol.(spec.Provider)
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
