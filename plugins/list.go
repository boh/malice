package plugins

import (
	"fmt"
	"strings"

	"github.com/crackcomm/go-clitable"
)

// ListEnabledPlugins lists enabled plugins
func ListEnabledPlugins(detail bool) {
	// TODO: Create a template for this kind of output : http://stackoverflow.com/questions/10747054/special-case-treatment-for-the-last-element-of-a-range-in-google-gos-text-templ
	enabled := getEnabled(Plugs.Plugins)
	if detail {
		ToMarkDownTable(enabled)
	} else {
		for _, plugin := range enabled {
			fmt.Println(plugin.Name)
		}
	}
}

// ListAllPlugins lists all plugins
func ListAllPlugins(detail bool) {
	plugins := Plugs.Plugins
	if detail {
		ToMarkDownTable(plugins)
	} else {
		for _, plugin := range plugins {
			fmt.Println(plugin.Name)
		}
	}
}

// ToMarkDownTable prints plugins out as Markdown table
func ToMarkDownTable(plugins []Plugin) {
	table := clitable.New([]string{"Name", "Description", "Enabled", "Image", "Category", "Mime"})
	for _, plugin := range plugins {
		table.AddRow(map[string]interface{}{
			"Name":        plugin.Name,
			"Description": plugin.Description,
			"Enabled":     plugin.Enabled,
			"Image":       plugin.Image,
			"Category":    plugin.Category,
			"Mime":        plugin.Mime,
		})
	}
	table.Markdown = true
	table.Print()
}

// GetPluginByName will return plugin for the given name
func GetPluginByName(name string) Plugin {
	for _, plugin := range Plugs.Plugins {
		if strings.EqualFold(plugin.Name, name) {
			return plugin
		}
	}
	return Plugin{}
}

// GetIntelPlugins will return all Intel plugins
func GetIntelPlugins(enabled bool) []Plugin {
	if enabled {
		// fmt.Printf("%#v\n", filterPluginsByIntel(Plugs.Plugins))
		return getIntel(getEnabled(getInstalled()))
	}
	return getIntel(getInstalled())
}

// GetPluginsForMime will return all plugins that can consume the mime type file
func GetPluginsForMime(mime string, enabled bool) []Plugin {
	if enabled {
		return getMime(mime, getEnabled(getInstalled()))
	}
	return getMime(mime, getInstalled())
}

func getIntel(plugins []Plugin) []Plugin {
	intel := []Plugin{}
	if plugins == nil {
		plugins = Plugs.Plugins
	}
	for _, plugin := range plugins {
		if strings.Contains(plugin.Category, "intel") {
			intel = append(intel, plugin)
		}
	}
	return intel
}

// getInstalled returns a map[string]plugin of installed plugins
func getInstalled() []Plugin {
	installed := []Plugin{}
	for _, plugin := range Plugs.Plugins {
		if plugin.Installed {
			installed = append(installed, plugin)
		}
	}
	return installed
}

// filterPluginsByEnabled returns a map[string]plugin of plugins
// that work on the given mime type
func getMime(mime string, plugins []Plugin) []Plugin {
	mimeMatch := []Plugin{}
	if plugins == nil {
		plugins = Plugs.Plugins
	}
	for _, plugin := range plugins {
		if strings.Contains(plugin.Mime, mime) || strings.Contains(plugin.Mime, "*") {
			mimeMatch = append(mimeMatch, plugin)
		}
	}
	return mimeMatch
}

// getEnabled returns a map[string]plugin of enabled plugins
func getEnabled(plugins []Plugin) []Plugin {
	enabled := []Plugin{}
	if plugins == nil {
		return Plugs.Plugins
	}
	for _, plugin := range Plugs.Plugins {
		if plugin.Enabled {
			enabled = append(enabled, plugin)
		}
	}
	return enabled
}
