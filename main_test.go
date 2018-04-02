package main

import (
	"os"
	"testing"

	"github.com/urfave/cli"
)

func Test(t *testing.T) {
	osExiterOriginal := cli.OsExiter
	defer func() { cli.OsExiter = osExiterOriginal }()

	logFatalOriginal := logFatal
	defer func() { logFatal = logFatalOriginal }()

	tests := []struct {
		title  string
		experr bool
		args   []string
	}{
		{
			title:  "success",
			experr: false,
			args:   []string{"-h"},
		},
		{
			title:  "error",
			experr: true,
			args:   []string{"not-existent"},
		},
	}

	// mock cli.OsExiter
	exitCode := 0
	cli.OsExiter = func(code int) {
		exitCode = code
	}

	// mock logFatal
	logFatal = func(...interface{}) {}

	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			args := make([]string, 0, len(tt.args)+1)
			args = append(args, "eskip-match")
			args = append(args, tt.args...)

			os.Args = args
			main()

			if exitCode > 0 && !tt.experr {
				t.Errorf("expected process exiting with success but an error occurred, exit code is = %d", exitCode)
			}

			if exitCode == 0 && tt.experr {
				t.Errorf("expected process exiting with error but got success, exit code is = %d", exitCode)
			}
		})
	}

}
