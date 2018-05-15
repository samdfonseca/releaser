package main

import (
	"log"
	"os"

	"gopkg.in/urfave/cli.v2"
	"github.com/axialmarket/releaser/notes/format"
	"github.com/axialmarket/releaser/notes/ptracker"
	"github.com/axialmarket/releaser/notes/wiki"
)

func writeDefaultFile(name, location string) error {
	if _, err := os.Stat(os.ExpandEnv(location)); err != nil && os.IsNotExist(err) {
		cf, err := os.Create(os.ExpandEnv(location))
		if err != nil {
			return err
		}
		defer cf.Close()
		_, err = cf.Write(files[name])
		if err != nil {
			return err
		}
	}
	return nil
}

func init() {
	if _, err := os.Stat(os.ExpandEnv("${HOME}/.releaser")); err != nil && os.IsNotExist(err) {
		if err = os.Mkdir(os.ExpandEnv("${HOME}/.releaser"), 0644); err != nil {
			log.Fatal(err)
		}
	}
	if err := writeDefaultFile("config.default.json", "${HOME}/.releaser/config.json"); err != nil {
		log.Fatal(err)
	}
	if err := writeDefaultFile("wiki.default.tmpl", "${HOME}/.releaser/wiki.tmpl"); err != nil {
		log.Fatal(err)
	}
}

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
