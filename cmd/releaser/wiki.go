package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/urfave/cli.v2"

	"github.com/samdfonseca/releaser/notes/wiki"
)

var (
	WIKI_URL = "https://wiki.axialmarket.com/api.php"
)

func updateWikiPage(c *cli.Context) error {
	page := c.String("page")
	if page == "" {
		return fmt.Errorf("missing required flag -page")
	}
	client, err := wiki.NewWikiClient(WIKI_URL)
	if err != nil {
		return err
	}
	text, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return err
	}
	err = client.UpdatePageText(page, string(text))
	if err != nil {
		return err
	}
	return nil
}
