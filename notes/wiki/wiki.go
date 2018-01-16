package wiki

import (
	"io"
	"text/template"

	"github.com/sadbox/mediawiki"
)

var (
	REL_NOTES_PAGE_TEMPLATE = `Release Date: {{ .ReleaseDate }}
{{range .Teams}}{{"{{RelNotesTeam|"}}{{ .TeamName }}|{{len .TeamItems }}{{"}}"}}
{{range .TeamItems}}{{"{{RelNotesTicket|"}}{{ .StoryLink }}|{{ .StoryName }}|{{ .StoryRepo }}|{{ .StoryPrLink }}{{"}}"}}
{{end}}
{{- end}}`
)

type WikiClient struct {
	WikiUrl string
	Client  *mediawiki.MWApi
}

type RelNotesItem struct {
	StoryLink   string
	StoryName   string
	StoryRepo   string
	StoryPrLink string
}

type RelNotesTeam struct {
	TeamName  string
	TeamItems []RelNotesItem
}

type RelNotesVars struct {
	ReleaseDate string
	Teams       []RelNotesTeam
}

func NewWikiClient(wikiUrl string) (*WikiClient, error) {
	mwApi, err := mediawiki.New(wikiUrl, "releaser-bot")
	if err != nil {
		return nil, err
	}
	return &WikiClient{
		WikiUrl: wikiUrl,
		Client:  mwApi,
	}, nil
}

func (wc *WikiClient) UpdatePageText(pageName, pageText string) error {
	editConfig := mediawiki.Values{
		"title": pageName,
		"text":  pageText,
	}
	return wc.Client.Edit(editConfig)
}

func (wc *WikiClient) ReadPage(pageName string) (string, error) {
	revision, err := wc.Client.Read(pageName)
	if err != nil {
		return "", err
	}
	return revision.Body, nil
}

func GenerateRelNotesText(relNotesVars RelNotesVars, w io.Writer) error {
	t := template.Must(template.New("relnotes").Parse(REL_NOTES_PAGE_TEMPLATE))
	if err := t.Execute(w, relNotesVars); err != nil {
		return err
	}
	return nil
}
