package wiki

import (
	"github.com/sadbox/mediawiki"
)

type WikiClient struct {
	WikiUrl string
	Client  *mediawiki.MWApi
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
