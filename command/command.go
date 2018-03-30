package command

import (
	"github.com/urfave/cli"
)

// Options configurable options
type Options struct {
	// holds global config file cli flag value
	ConfigFile string
}

// GetCommands returns cli commands list
func GetCommands(o *Options) []cli.Command {
	return []cli.Command{
		NewTest(o),
	}
}
