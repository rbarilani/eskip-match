package command

import (
	"flag"
	"github.com/urfave/cli"
	"strings"
	"testing"
)

func TestApp(t *testing.T) {
	app := NewApp()
	if app == nil {
		t.Error("NewApp return nil")
	}

	scenarios := []struct {
		Args     []string
		ExpError bool
	}{
		{
			Args: []string{"eskip-match", "test", "command_test.eskip", "-p", "/bar"},
		},
		{
			ExpError: true,
			Args:     []string{"eskip-match", "test", "command_test.eskip", "-p", "/foofoo"},
		},
	}

	for _, s := range scenarios {
		t.Run(strings.Join(s.Args[:], " "), func(t *testing.T) {
			err := app.Run(s.Args)
			if s.ExpError && err == nil {
				t.Error("expecting error but got nil")
			}
			if s.ExpError == false && err != nil {
				t.Error("expecting match but got error", err)
			}
		})
	}
}

func TestNewTest(t *testing.T) {
	testcommand := newTest(&Options{})
	fn, ok := testcommand.Action.(func(c *cli.Context) error)

	if !ok {
		t.Error("cannot cast test command action to function")
		return
	}

	set := flag.NewFlagSet("test", 0)
	set.Parse([]string{"command_test.eskip"})
	set.Bool("v", true, "")
	set.String("p", "/bar", "")
	ctx := cli.NewContext(nil, set, nil)
	ctx.Command = testcommand

	err := fn(ctx)

	if err != nil {
		t.Error(err)
	}
}
