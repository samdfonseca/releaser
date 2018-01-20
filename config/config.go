package config

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"path"
)

type ReleaserConfig struct {
	WikiUrl              string `json:"wiki_url"`
	PivotalApiToken      string `json:"pivotal_api_token"`
	PivotalProjectIds    []int  `json:"pivotal_project_ids"`
	GithubOrg            string `json:"github_org"`
	RelNotesTemplateFile string `json:"rel_notes_template_file"`
}

func New(r io.Reader) (*ReleaserConfig, error) {
	contents, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	var config ReleaserConfig
	if err = json.Unmarshal(contents, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

func (rc *ReleaserConfig) GetFullRelNotesTemplateFilePath() string {
	if path.IsAbs(rc.RelNotesTemplateFile) {
		return rc.RelNotesTemplateFile
	}
	return path.Join(os.ExpandEnv("${HOME}/.releaser"), rc.RelNotesTemplateFile)
}
