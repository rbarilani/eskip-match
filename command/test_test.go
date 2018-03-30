package command

import (
	"flag"
	"github.com/rbarilani/eskip-match/config"
	"github.com/urfave/cli"
	"testing"
)

func TestNewTest(t *testing.T) {
	testcommand := newTest(&Options{
		ConfigLoader: config.NewLoader(""),
	})
	fn, ok := testcommand.Action.(func(c *cli.Context) error)

	if !ok {
		t.Error("cannot cast test command action to function")
		return
	}

	scenarios := []struct {
		title    string
		args     []string
		expError bool
		setFlags func(*flag.FlagSet)
	}{
		{
			title:    "must provide a routes file error",
			expError: true,
			setFlags: func(set *flag.FlagSet) {},
		},
		{
			title:    "routes file doesnt exist",
			args:     []string{"blue.eskip"},
			expError: true,
			setFlags: func(set *flag.FlagSet) {},
		},
		{
			title: "success",
			args:  []string{"command_test.eskip"},
			setFlags: func(set *flag.FlagSet) {
				set.String("p", "/bar", "")
			},
		},
	}

	for _, s := range scenarios {
		t.Run(s.title, func(t *testing.T) {
			set := flag.NewFlagSet("test", 0)
			set.Parse(s.args)
			s.setFlags(set)

			ctx := cli.NewContext(nil, set, nil)
			ctx.Command = testcommand

			err := fn(ctx)

			if s.expError == false && err != nil {
				t.Error("not expected error occurred", err)
			}
		})
	}
}
