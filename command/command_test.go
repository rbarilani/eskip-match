package command

import (
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
