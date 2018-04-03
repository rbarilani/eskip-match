package cli

import (
	"flag"
	"testing"

	"github.com/urfave/cli"
)

func TestApp(t *testing.T) {
	app := NewApp()
	if app == nil {
		t.Error("NewApp return nil")
	}

	tests := []struct {
		name string
		args []string
		err  bool
	}{
		{
			name: "match",
			args: []string{"eskip-match", "test", "testdata/routes.eskip", "-p", "/bar"},
		},
		{
			name: "dont-match",
			err:  true,
			args: []string{"eskip-match", "test", "testdata/routes.eskip", "-p", "/foofoo"},
		},
		{
			name: "headers",
			args: []string{
				"eskip-match",
				"test",
				"testdata/routes.eskip",
				"-p", "/bar",
				"-H",
				"foo=bar",
				"-H",
				"bar=foo",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := app.Run(tt.args)
			if tt.err && err == nil {
				t.Error("expecting error but got nil")
			}
			if tt.err == false && err != nil {
				t.Error("expecting match but got error", err)
			}
		})
	}
}

func TestNewTestCommand(t *testing.T) {
	testcommand := newTestCommand(&options{
		ConfigLoader: newConfigLoader(""),
	})
	fn, ok := testcommand.Action.(func(c *cli.Context) error)

	if !ok {
		t.Error("cannot cast test command action to function")
		return
	}

	tests := []struct {
		name     string
		args     []string
		err      bool
		setFlags func(*flag.FlagSet)
	}{
		{
			name:     "must provide a routes file error",
			err:      true,
			setFlags: func(set *flag.FlagSet) {},
		},
		{
			name:     "routes file doesnt exist",
			args:     []string{"blue.eskip"},
			err:      true,
			setFlags: func(set *flag.FlagSet) {},
		},
		{
			name: "success",
			args: []string{"testdata/routes.eskip"},
			setFlags: func(set *flag.FlagSet) {
				set.String("p", "/bar", "")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			set := flag.NewFlagSet("test", 0)
			set.Parse(tt.args)
			tt.setFlags(set)

			ctx := cli.NewContext(nil, set, nil)
			ctx.Command = testcommand

			err := fn(ctx)

			if tt.err == false && err != nil {
				t.Error("not expected error occurred", err)
			}
		})
	}
}
