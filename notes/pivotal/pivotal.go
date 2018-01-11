package pivotal

import (
	"fmt"

	"gopkg.in/salsita/go-pivotaltracker.v1/v5/pivotal"
)

type ProjectLister interface {
	ListProjects() ([]*pivotal.Project, error)
}

type ProjectGetter interface {
	GetProject(int) (*pivotal.Project, error)
}

type ProjectListerGetter interface {
	ProjectLister
	ProjectGetter
}

type StoryGetter interface {
	GetStory(int, int) (*pivotal.Story, error)
}

type PivotalClient struct {
	client *pivotal.Client
}

func NewClient(apiToken string) *PivotalClient {
	client := pivotal.NewClient(apiToken)
	return &PivotalClient{
		client: client,
	}
}

func (pc *PivotalClient) ListProjects() ([]*pivotal.Project, error) {
	projects, _, err := pc.client.Projects.List()
	return projects, err
}

func (pc *PivotalClient) GetProject(projectId int) (*pivotal.Project, error) {
	project, _, err := pc.client.Projects.Get(projectId)
	return project, err
}

func (pc *PivotalClient) GetStory(projectId, storyId int) (*pivotal.Story, error) {
	story, _, err := pc.client.Stories.Get(projectId, storyId)
	return story, err
}

func GetStoryProjectId(pl ProjectLister, storyId int) (int, error) {
	projects, err := pl.ListProjects()
	if err != nil {
		return 0, err
	}
	for _, project := range projects {
		for _, projectStoryId := range project.StoryIds {
			if storyId == projectStoryId {
				return project.Id, nil
			}
		}
	}
	return 0, fmt.Errorf("Unable to find story")
}

func GetStoryProject(plg ProjectListerGetter, storyId int) (*pivotal.Project, error) {
	projectId, err := GetStoryProjectId(plg, storyId)
	if err != nil {
		return nil, err
	}
	project, err := plg.GetProject(projectId)
	return project, err
}

func GetStory(pl ProjectLister, sg StoryGetter, storyId int) (*pivotal.Story, error) {
	projectId, err := GetStoryProjectId(pl, storyId)
	if err != nil {
		return nil, err
	}
	story, err := sg.GetStory(projectId, storyId)
	return story, err
}
