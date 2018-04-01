package command

import (
	"flag"
	"strings"
	"testing"

	"github.com/rbarilani/eskip-match/config"
	"github.com/urfave/cli"
)

func TestApp(t *testing.T) {
	app := NewApp()
	if app == nil {
		t.Error("NewApp return nil")
	}

	scenarios := []struct {
		title    string
		args     []string
		expError bool
	}{
		{
			title: "match",
			args:  []string{"eskip-match", "test", "command_test.eskip", "-p", "/bar"},
		},
		{
			title:    "dont-match",
			expError: true,
			args:     []string{"eskip-match", "test", "command_test.eskip", "-p", "/foofoo"},
		},
		{
			title: "headers",
			args: []string{
				"eskip-match",
				"test",
				"command_test.eskip",
				"-p", "/bar",
				"-H",
				"foo=bar",
				"-H",
				"bar=foo",
			},
		},
	}

	for _, s := range scenarios {
		t.Run(strings.Join(s.args[:], ", "), func(t *testing.T) {
			err := app.Run(s.args)
			if s.expError && err == nil {
				t.Error("expecting error but got nil")
			}
			if s.expError == false && err != nil {
				t.Error("expecting match but got error", err)
			}
		})
	}
}

func TestNewTestCommand(t *testing.T) {
	testcommand := newTestCommand(&options{
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
