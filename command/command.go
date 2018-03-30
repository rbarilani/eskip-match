package command

import (
	"github.com/rbarilani/eskip-match/config"
	"github.com/urfave/cli"
)

// Options configurable options
type Options struct {
	// holds global config file cli flag value
	ConfigFile string

	ConfigLoader config.Loader
}

// NewApp creates the cli application
func NewApp() *cli.App {
	var configFile string
	loader := config.NewLoader(config.DefaultFile)

	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "config, c",
			Usage:       "Load configuration from `FILE`",
			Destination: &configFile,
		},
	}
	app.Name = "eskip-match"
	app.Usage = "A command line tool that helps you test .eskip files routing matching logic"

	app.Commands = getCommands(&Options{
		ConfigFile:   configFile,
		ConfigLoader: loader,
	})
	return app
}

// GetCommands returns cli commands list
func getCommands(o *Options) []cli.Command {
	return []cli.Command{
		newTest(o),
	}
}
