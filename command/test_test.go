package command

import (
	"flag"
	"github.com/urfave/cli"
	"testing"
)

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
