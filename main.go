// Package eskip-match provides a cli tool and utilities to
// helps you test Skipper (https://github.com/zalando/skipper)
// `.eskip` files routing matching logic.
package main

import (
	"log"
	"os"

	"github.com/rbarilani/eskip-match/cli"
)

var logFatal = log.Fatal

func main() {
	app := cli.NewApp()
	app.Version = "0.3.3"
	err := app.Run(os.Args)
	if err != nil {
		logFatal(err)
	}
}
