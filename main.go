package main

import (
	"log"
	"os"

	"github.com/rbarilani/eskip-match/command"
)

func main() {
	app := command.NewApp()
	app.Version = "0.3.1"
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
