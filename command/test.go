package command

import (
	"fmt"
	"github.com/rbarilani/eskip-match/config"
	"github.com/rbarilani/eskip-match/matcher"
	"github.com/urfave/cli"
)

// NewTest creates `test` cli command
func NewTest(o *Options) cli.Command {
	return cli.Command{
		Name:    "test",
		Aliases: []string{"t"},
		Usage:   "given a routes file and request attributes, checks if a route match",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "method, m",
				Usage: "Specify an http method",
			},
			cli.StringFlag{
				Name:  "path, p",
				Usage: "Specify an http path",
			},
			cli.BoolFlag{
				Name:  "verbose, v",
				Usage: "Print verbose output",
			},
		},
		Action: func(c *cli.Context) error {
			conf := config.Load(o.ConfigFile)
			routesFile := c.Args().First()
			m, err := matcher.New(&matcher.Options{
				RoutesFile:    routesFile,
				CustomFilters: matcher.MockFilters(conf.CustomFilters),
				Verbose:       c.Bool("v"),
			})
			if err != nil {
				return err
			}
			reqAttrs := &matcher.RequestAttributes{
				Method: c.String("m"),
				Path:   c.String("p"),
			}
			route := m.Test(reqAttrs)

			if route == nil {
				return fmt.Errorf("no match")
			}

			fmt.Print(matcher.PrettyPrintRoute(route))

			return nil
		},
	}
}
