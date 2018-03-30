package command

import (
	"fmt"
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
			reqAttrs := &matcher.RequestAttributes{
				Method:  strings.ToUpper(c.String("m")),
				Path:    c.String("p"),
				Headers: headers(c.StringSlice("H")),
			}
			res, err := m.Test(reqAttrs)
			attrs := res.Attributes()
			log.Printf("request: %s %s", attrs.Method, attrs.Path)
			if len(attrs.Headers) > 0 {
				pairs := make([]string, 0, len(attrs.Headers))
				for key, value := range attrs.Headers {
					pairs = append(pairs, key+"="+value)
				}
				log.Printf("request headers: %s", strings.Join(pairs, ", "))
			}

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

func headers(pairs []string) map[string]string {
	m := make(map[string]string)
	for _, pair := range pairs {
		parts := strings.Split(pair, "=")
		if len(parts) >= 2 {
			m[parts[0]] = parts[1]
		}
	}
	return m
}
