package cli

import (
	"os"

	"github.com/jinzhu/configor"
)

// DefaultFile default configuration file
const configDefaultFile = ".eskip-match.yml"

// Config root element for configuration
type config struct {
	CustomFilters []string
}

// Loader ...
type configLoader interface {
	Load(file string) config
}

type configorLoader struct {
	defaultFile string
}

// NewLoader ...
func newConfigLoader(defaultFile string) configLoader {
	return &configorLoader{
		defaultFile,
	}
}

// Load configuration from a file
func (c *configorLoader) Load(file string) config {
	var conf config
	configorConf := &configor.Config{ENVPrefix: "EM"}

	if file == "" {
		// search for default configuration file
		file = c.defaultFile
		if _, err := os.Stat(file); os.IsNotExist(err) {
			file = ""
		}
	}

	if file == "" {
		configor.New(configorConf).Load(&conf)
	} else {
		configor.New(configorConf).Load(&conf, file)
	}
	return conf
}
