package wiki

import (
	"fmt"
	"io"
	"log"
	"os"

	"gopkg.in/urfave/cli.v2"

	releaserConfig "github.com/axialmarket/releaser/config"

	"github.com/axialmarket/releaser/flags"
	"github.com/axialmarket/releaser/logging"
)

var (
	logger = logging.New("wiki")
)

func updateWikiPage(c *cli.Context) error {
	configPath := c.String("config")
	logger.Debugf("Using config file: %s\n", configPath)
	f, err := os.Open(configPath)
	if err != nil {
		log.Fatalf("Unable to find config file: %s", err)
	}
	config, err := releaserConfig.New(f)
	if err != nil {
		return err
	}
	page := c.String("page")
	if page == "" {
		return fmt.Errorf("missing required flag -page")
	}
	logger.Debugf("Creating mediawiki client: %s\n", config.WikiUrl)
	client, err := NewWikiClient(config.WikiUrl)
	if err != nil {
		return err
	}
	logger.Debugf("Updating page: %s\n", page)
	var stdin io.Reader
	stdin = os.Stdin
	if c.Bool("tee") {
		stdin = io.TeeReader(os.Stdin, os.Stdout)
	}
	err = client.UpdatePageTextReader(page, stdin)
	if err != nil {
		return err
	}
	return nil
}

func Command() *cli.Command {
	return &cli.Command{
		Name:    "wikipage",
		Aliases: []string{"wiki", "wp"},
		Usage:   "updates the contents of a wiki page",
		Action:  updateWikiPage,
		Flags: []cli.Flag{
			flags.Config,
			&cli.StringFlag{
				Name:  "page",
				Aliases: []string{"p"},
				Usage: "name of the wiki page to update",
			},
			&cli.StringFlag{
				Name: "template",
				Aliases: []string{"t"},
				Usage: "name of the template file used to generate output",
			},
			&cli.BoolFlag{ Name: "tee",
				Value: false,
				Usage: "print out the generated wiki page text",
			},
		},
	}
}
