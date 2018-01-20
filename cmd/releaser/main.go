package main

import (
	"os"

	"gopkg.in/urfave/cli.v2"
	"github.com/axialmarket/releaser/notes/format"
	"github.com/axialmarket/releaser/notes/ptracker"
	"github.com/axialmarket/releaser/notes/wiki"
)

func main() {
	app := &cli.App{
		Name:  "releaser",
		Usage: "Axial release automation tasks",
		Commands: []*cli.Command{
			format.Command(),
			ptracker.Command(),
			wiki.Command(),
		},
	}

	app.Run(os.Args)
}
