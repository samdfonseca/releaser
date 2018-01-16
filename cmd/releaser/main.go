package main

import (
	"os"

	"gopkg.in/urfave/cli.v2"
)

func main() {
	app := &cli.App{
		Name:  "releaser",
		Usage: "Axial release automation tasks",
		Commands: []*cli.Command{
			{
				Name:    "relnotes",
				Aliases: []string{"rn"},
				Usage:   "generate the wiki page release notes",
				Action:  getReleaseStories,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "label",
						Usage: "pivotal label for stories in release",
					},
					&cli.StringFlag{
						Name:  "relDate",
						Usage: "date of release in yyyy-mm-dd format",
					},
					&cli.StringFlag{
						Name:  "config",
						Usage: "path to the config.json file",
						Value: "./notes_config.json",
					},
				},
			},
			{
				Name:    "wikipage",
				Aliases: []string{"wiki", "wp"},
				Usage:   "updates the contents of a wiki page",
				Action:  updateWikiPage,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "page",
						Usage: "name of the wiki page to update",
					},
					&cli.StringFlag{
						Name:  "config",
						Usage: "path to the config.json file",
						Value: "./notes_config.json",
					},
				},
			},
		},
	}

	app.Run(os.Args)
}
