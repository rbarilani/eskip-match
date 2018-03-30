package command

import (
	"fmt"
	"github.com/rbarilani/eskip-match/config"
	"github.com/rbarilani/eskip-match/matcher"
	"github.com/urfave/cli"
	"log"
	"strings"
)

// NewTest creates `test` cli command
func newTest(o *Options) cli.Command {
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
			cli.BoolFlag{
				Name:  "verbose, v",
				Usage: "Print verbose output",
			},
		},
		Action: func(c *cli.Context) error {
			conf := config.Load(o.ConfigFile)
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
			reqAttrs := &matcher.RequestAttributes{
				Method: strings.ToUpper(c.String("m")),
				Path:   c.String("p"),
			}
			res, err := m.Test(reqAttrs)
			log.Printf("request: %s %s", res.Attributes().Method, res.Attributes().Path)
			if err != nil {
				return err
			}
			route := res.Route()
			if route == nil {
				return fmt.Errorf("no match")
			}
			log.Println("matching route id:", route.Id)
			log.Printf("matching route:\n```\n%s```", res.PrettyPrintRoute())

			return nil
		},
	}
}
