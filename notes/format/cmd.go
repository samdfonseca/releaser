package format

import (
	"io"
	"log"
	"os"

	"gopkg.in/urfave/cli.v2"

	releaserConfig "github.com/axialmarket/releaser/config"

	"github.com/axialmarket/releaser/flags"
	"github.com/axialmarket/releaser/logging"
	"github.com/axialmarket/releaser/notes"
)

var (
	logger = logging.New("fmt")
)

func formatRelNotesData(c *cli.Context) error {
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
	var relNotesVars notes.RelNotesVars
	logger.Debugln("Reading release notes data from stdin")
	if err := ReadRelNotesVars(os.Stdin, &relNotesVars); err != nil {
		return err
	}
	r, w := io.Pipe()
	defer func() {
		r.Close()
		w.Close()
	}()
	tmplFilePath := c.String("template")
	if tmplFilePath == "" {
		tmplFilePath = config.GetFullRelNotesTemplateFilePath()
	}
	logger.Debugf("Using template file: %s\n", tmplFilePath)
	tmplFile, err := os.Open(tmplFilePath)
	defer tmplFile.Close()
	if err != nil {
		return err
	}
	if err := CompileRelNotesTemplate(relNotesVars, tmplFile, os.Stdout); err != nil {
		return err
	}
	return nil
}

func Command() *cli.Command {
	return &cli.Command{
		Name:    "format",
		Aliases: []string{"fmt"},
		Usage:   "applies the release notes json data to the provided template",
		Action:  formatRelNotesData,
		Flags: []cli.Flag{
			flags.Config,
			&cli.StringFlag{
				Name:    "template",
				Aliases: []string{"t"},
				Usage:   "name of the template file used to generate output",
			},
		},
	}
}
