package main

import (
	"log"
	"os"

	"github.com/mcclurejt/mrkt-backend/tablegen"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:     "tg",
		HelpName: "tablegen generates a table of stock information based on the symbols in 'symbols.config'",
		Action:   tablegen.Action,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
