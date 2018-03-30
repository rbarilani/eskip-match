package config

import (
	"github.com/jinzhu/configor"
	"os"
)

// DefaultFile default configuration file
const DefaultFile = ".eskip-match.yml"

// Config root element for configuration
type Config struct {
	CustomFilters []string
}

// Loader ...
type Loader interface {
	Load(file string) Config
}

type configLoader struct {
	defaultFile string
}

// NewLoader ...
func NewLoader(defaultFile string) Loader {
	return &configLoader{
		defaultFile,
	}
}

// Load configuration from a file
func (c *configLoader) Load(file string) Config {
	var config Config
	configorConfig := &configor.Config{ENVPrefix: "EM"}

	if file == "" {
		// search for default configuration file
		file = c.defaultFile
		if _, err := os.Stat(file); os.IsNotExist(err) {
			file = ""
		}
	}

	if file == "" {
		configor.New(configorConfig).Load(&config)
	} else {
		configor.New(configorConfig).Load(&config, file)
	}
	return config
}
