package ptracker

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"gopkg.in/urfave/cli.v2"

	releaserConfig "github.com/axialmarket/releaser/config"
	"github.com/axialmarket/releaser/notes"
	"github.com/axialmarket/releaser/flags"
)

func getReleaseStories(c *cli.Context) error {
	configPath := c.String("config")
	f, err := os.Open(configPath)
	if err != nil {
		log.Fatalf("Unable to find config file: %s", err)
	}
	config, err := releaserConfig.New(f)
	if err != nil {
		return err
	}
	label := c.String("label")
	if label == "" {
		label = fmt.Sprintf("rc-%s", time.Now().Format("2006-01-02"))
	}
	relDate := c.String("relDate")
	if relDate == "" {
		relDate = time.Now().Format("2006-01-02")
	}
	relNotesVars := notes.RelNotesVars{
		ReleaseDate: relDate,
		Projects:       []notes.RelNotesProject{},
	}
	ptClient := NewClient(config.PivotalApiToken)
	prLinkRegexp, err := GetPrLinkRegexp(config.GithubOrg)
	if err != nil {
		return err
	}
	for _, projId := range config.PivotalProjectIds {
		projClient := ptClient.InProject(projId)
		proj, err := projClient.Project()
		if err != nil {
			return err
		}
		relNotesTeam := notes.RelNotesProject{
			ProjectName:  proj.Name,
			ProjectStories: []notes.RelNotesStory{},
		}
		projStories, err := GetStoriesWithLabel(projClient, label)
		if err != nil {
			return err
		}
		for _, story := range projStories {
			relNotesItem := notes.RelNotesStory{
				StoryLink:   story.URL,
				StoryName:   story.Name,
				StoryRepo:   "no PR",
				StoryPrLink: "no PR",
			}
			prUrls := GetPrLinksFromStory(story, prLinkRegexp)
			if len(prUrls) > 0 {
				parsedPrUrl, err := url.Parse(prUrls[0])
				if err != nil {
					return err
				}
				prUrlPath := parsedPrUrl.EscapedPath()
				repo := strings.Split(prUrlPath, "/")[2]
				relNotesItem.StoryRepo = repo
				relNotesItem.StoryPrLink = prUrls[0]
			}
			relNotesTeam.ProjectStories = append(relNotesTeam.ProjectStories, relNotesItem)
		}
		if len(relNotesTeam.ProjectStories) > 0 {
			relNotesVars.Projects = append(relNotesVars.Projects, relNotesTeam)
		}
	}
	notesJson, err := json.Marshal(relNotesVars)
	if err != nil {
		return err
	}
	if _, err := os.Stdout.Write(notesJson); err != nil {
		return err
	}
	return nil
}

func Command() *cli.Command {
	return &cli.Command{
		Name:    "relnotes",
		Aliases: []string{"rn"},
		Usage:   "generate the wiki page release notes",
		Action:  getReleaseStories,
		Flags: []cli.Flag{
			flags.Config,
			&cli.StringFlag{
				Name:  "label",
				Usage: "pivotal label for stories in release",
			},
			&cli.StringFlag{
				Name:  "relDate",
				Usage: "date of release in yyyy-mm-dd format",
			},
		},
	}
}
