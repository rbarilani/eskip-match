package config

import (
	"github.com/jinzhu/configor"
	"os"
)

// Config root element for configuration
type Config struct {
	CustomFilters []string
}

// Load configuration from a file
func Load(file string) Config {
	var config Config
	configorConfig := &configor.Config{ENVPrefix: "EM"}

	if file == "" {
		// search for default configuration file
		file = "eskip-match.yml"
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

// Default configuration
func Default() Config {
	return Load("")
}
