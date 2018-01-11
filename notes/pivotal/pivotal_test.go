package pivotal

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/salsita/go-pivotaltracker.v1/v5/pivotal"
)

type MockProjectLister struct {
	OnListProjects func() ([]*pivotal.Project, error)
}

type MockProjectGetter struct {
	OnGetProject func(int) (*pivotal.Project, error)
}

type MockProjectListerGetter struct {
	*MockProjectLister
	*MockProjectGetter
}

type MockStoryGetter struct {
	OnGetStory func(int, int) (*pivotal.Story, error)
}

var (
	mockProjectIds = []int{100, 101, 102, 103, 104}
	mockStoryIds   = []int{200, 201, 202, 203, 204}
)

func NewMockProjectLister() *MockProjectLister {
	return &MockProjectLister{
		OnListProjects: func() ([]*pivotal.Project, error) {
			var projects []*pivotal.Project
			mockStoryId := 200
			for _, projectId := range mockProjectIds {
				var storyIds []int
				for i := 0; i < 5; i++ {
					storyIds = append(storyIds, mockStoryId)
					mockStoryId = mockStoryId + 1
				}
				projects = append(projects, &pivotal.Project{Id: projectId, StoryIds: storyIds})
			}
			return projects, nil
		},
	}
}

func (m MockProjectLister) ListProjects() ([]*pivotal.Project, error) {
	return m.OnListProjects()
}

func NewMockProjectGetter() *MockProjectGetter {
	return &MockProjectGetter{
		OnGetProject: func(projectId int) (*pivotal.Project, error) {
			isExistingProjectId := false
			for _, mockProjectId := range mockProjectIds {
				if mockProjectId == projectId {
					isExistingProjectId = true
				}
			}
			if !isExistingProjectId {
				return nil, fmt.Errorf("Project does not exist")
			}
			var storyIds []int
			for i := 200; i < 205; i++ {
				storyIds = append(storyIds, i)
			}
			return &pivotal.Project{Id: projectId, StoryIds: storyIds}, nil
		},
	}
}

func (m MockProjectGetter) GetProject(projectId int) (*pivotal.Project, error) {
	return m.OnGetProject(projectId)
}

func NewMockProjectListerGetter() *MockProjectListerGetter {
	return &MockProjectListerGetter{NewMockProjectLister(), NewMockProjectGetter()}
}

func NewMockStoryGetter() *MockStoryGetter {
	return &MockStoryGetter{
		OnGetStory: func(projectId, storyId int) (*pivotal.Story, error) {
			isExistingProjectId := false
			for _, mockProjectId := range mockProjectIds {
				if mockProjectId == projectId {
					isExistingProjectId = true
				}
			}
			if !isExistingProjectId {
				return nil, fmt.Errorf("Project does not exist")
			}
			return &pivotal.Story{Id: storyId, ProjectId: projectId, Name: "Mock Story", Description: "sam test"}, nil
		},
	}
}

func (m MockStoryGetter) GetStory(projectId, storyId int) (*pivotal.Story, error) {
	return m.OnGetStory(projectId, storyId)
}

func TestGetStoryProjectId(t *testing.T) {
	pl := NewMockProjectLister()
	storyId := 200
	for _, mockProjectId := range mockProjectIds {
		for i := 0; i < 5; i++ {
			projectId, err := GetStoryProjectId(*pl, storyId)
			assert.Equal(t, err, nil, "err should be nil")
			assert.Equal(t, projectId, mockProjectId, "projectId should equal mockProjectId")
			storyId = storyId + 1
		}
	}
}

func TestGetActualStory(t *testing.T) {
	pivotalApiToken := os.Getenv("PIVOTAL_API_TOKEN")
	assert.NotEqual(t, pivotalApiToken, "", "pivotalApiToken should not be empty")
	client := NewClient(pivotalApiToken)
	story, _, err := client.client.Stories.Get(1974589, 154134429)
	assert.Equal(t, nil, err, "err should be nil")
	assert.NotEqual(t, nil, story, "story should not be nil")
	project, _, err := client.client.Projects.Get(1974589)
	assert.Equal(t, nil, err, "err should be nil")
	assert.Equal(t, 1974589, project.Id, "projectId should equal 1974589")
	assert.NotEqual(t, 0, len(project.StoryIds), "project.StoryIds should have > 0 items")
}
