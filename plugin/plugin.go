package plugin

import (
	"gopkg.in/yaml.v2"

	"io/ioutil"
	"log"
)

// Plugins represents the configuration information.
type Plugins struct {
	Bin BinaryPlugin   `yaml:"binary"`
	Doc DocumantPlugin `yaml:"document"`
}

// BinaryPlugin represents the Email configuration details
type BinaryPlugin struct {
	Name struct {
		Enabled string `yaml:"enabled"`
		Image   string `yaml:"image"`
	}
}

// DocumantPlugin represents the Database configuration details
type DocumantPlugin struct {
	Name struct {
		Enabled string `yaml:"enabled"`
		Image   string `yaml:"image"`
	}
}

// Plugin represents the Malice regiestered Plugins
var Plugin Plugins

func init() {
	// Get the config file
	plugins, err := ioutil.ReadFile("./plugins.yaml")
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	yaml.Unmarshal(plugins, &Plugin)
}
