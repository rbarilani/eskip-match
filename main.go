package main

import (
	"github.com/rbarilani/eskip-match/command"
	"log"
	"os"
)

func main() {
	app := command.NewApp()
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
