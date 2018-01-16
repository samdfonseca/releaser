package main

import (
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"gopkg.in/urfave/cli.v2"

	"github.com/samdfonseca/releaser/notes/ptracker"
	"github.com/samdfonseca/releaser/notes/wiki"
)

var (
	PIVOTAL_API_TOKEN     = os.Getenv("PIVOTAL_API_TOKEN")
	PIVOTAL_PROJECT_IDS   = []int{2077921, 2123086, 1974589, 2099793}
	PIVOTAL_PROJECT_NAMES = map[int]string{
		2077921: "TAXONIMY",
		2123086: "FLOATING",
		1974589: "OUTREACH",
		2099793: "BUY-SIDE (GHB)",
	}
	WIKI_URL = "https://wiki.axialmarket.com/api.php"
)

func getReleaseStories(c *cli.Context) error {
	label := c.String("label")
	if label == "" {
		label = fmt.Sprintf("rc-%s", time.Now().Format("2006-01-02"))
	}
	var relNotesVars wiki.RelNotesVars
	ptClient := ptracker.NewClient(PIVOTAL_API_TOKEN)
	for _, projId := range PIVOTAL_PROJECT_IDS {
		relNotesTeam := wiki.RelNotesTeam{
			TeamName:  PIVOTAL_PROJECT_NAMES[projId],
			TeamItems: []wiki.RelNotesItem{},
		}
		projClient := ptClient.InProject(projId)
		projStories, err := ptracker.GetStoriesWithLabel(projClient, label)
		if err != nil {
			return err
		}
		for _, story := range projStories {
			relNotesItem := wiki.RelNotesItem{
				StoryLink:   story.URL,
				StoryName:   story.Name,
				StoryRepo:   "no PR",
				StoryPrLink: "no PR",
			}
			prUrls := ptracker.GetPrLinksFromStory(story)
			if len(prUrls) > 0 {
				parsedPrUrl, err := url.Parse(prUrls[0])
				if err != nil {
					return err
				}
				prUrlPath := parsedPrUrl.EscapedPath()
				repo := strings.Split(prUrlPath, "/")[3]
				relNotesItem.StoryRepo = repo
				relNotesItem.StoryPrLink = prUrls[0]
			}
			relNotesTeam.TeamItems = append(relNotesTeam.TeamItems, relNotesItem)
		}
		if len(relNotesTeam.TeamItems) > 0 {
			relNotesVars.Teams = append(relNotesVars.Teams, relNotesTeam)
		}
	}
	if err := wiki.GenerateRelNotesText(relNotesVars, os.Stdout); err != nil {
		return err
	}
	return nil
}

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
				},
			},
		},
	}

	app.Run(os.Args)
}
