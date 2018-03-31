package command

import (
	"fmt"
	"github.com/rbarilani/eskip-match/config"
	"github.com/rbarilani/eskip-match/matcher"
	"github.com/urfave/cli"
	"log"
	"strings"
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
	options := &Options{
		ConfigFile:   configFile,
		ConfigLoader: loader,
	}

	app := cli.NewApp()
	app.Version = "0.3.0"
	app.Name = "eskip-match"
	app.Usage = "A command line tool that helps you test .eskip files routing matching logic"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "config, c",
			Usage:       "Load configuration from `FILE`",
			Destination: &configFile,
		},
	}

	app.Commands = []cli.Command{
		newTestCommand(options),
	}
	return app
}

func newTestCommand(o *Options) cli.Command {
	return cli.Command{
		Name:      "test",
		Aliases:   []string{"t"},
		ArgsUsage: "<ROUTES_FILE>",
		Usage:     "Given a routes file and request attributes, checks a route matches",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "method, m",
				Usage: "Specify request `METHOD`",
			},
			cli.StringFlag{
				Name:  "path, p",
				Usage: "Specify request `PATH`",
			},
			cli.StringSliceFlag{
				Name:  "header, H",
				Usage: "Specify request `HEADER` as key=value pair",
			},
			cli.BoolFlag{
				Name:  "verbose, v",
				Usage: "Print verbose output",
			},
		},
		Action: func(c *cli.Context) error {
			conf := o.ConfigLoader.Load(o.ConfigFile)
			routesFile := c.Args().First()
			if routesFile == "" {
				return fmt.Errorf("A routes file must be provided")
			}

			m, err := matcher.New(&matcher.Options{
				RoutesFile:    routesFile,
				CustomFilters: matcher.MockFilters(conf.CustomFilters),
				Verbose:       c.Bool("v"),
			})
			if err != nil {
				return err
			}

			res := m.Test(&matcher.RequestAttributes{
				Method:  strings.ToUpper(c.String("m")),
				Path:    c.String("p"),
				Headers: headers(c.StringSlice("H")),
			})

			out := res.PrettyPrintLines()
			route := res.Route()
			for _, line := range out {
				log.Println(line)
			}
			if route == nil {
				return fmt.Errorf("no match")
			}
			return nil
		},
	}
}

// headers parses a list of strings representing http headers
// written in a key=value format (eg. ["Content=text", "Accept=json"])
// and transform them to a map (eg. {"Content":"text", "Accept": "json"}).
// It discards items with an invalid format (eg. ["Content="])
func headers(pairs []string) map[string]string {
	m := make(map[string]string)
	for _, pair := range pairs {
		// Use SplitN to handle "=" symbol in the value
		// bar=bar=foo -> "bar": "bar=foo"
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) >= 2 {
			m[parts[0]] = parts[1]
		}
	}
	return m
}
