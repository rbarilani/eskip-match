package main

import (
	"log"
	"os"

	"github.com/rbarilani/eskip-match/cli"
)

func main() {
	app := cli.NewApp()
	app.Version = "0.3.2"
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
