package config

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

type NotesConfig struct {
	WikiUrl           string `json:"wiki_url"`
	PivotalApiToken   string `json:"pivotal_api_token"`
	PivotalProjectIds []int  `json:"pivotal_project_ids"`
	GithubOrg         string `json:"github_org"`
}

func NewNotesConfig(r io.Reader) (*NotesConfig, error) {
	contents, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	var config NotesConfig
	if err = json.Unmarshal(contents, &config); err != nil {
		return nil, err
	}
	return &config, nil
}
