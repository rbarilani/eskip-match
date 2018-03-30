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
