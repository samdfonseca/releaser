package ptracker

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xoebus/go-tracker"
)

var (
	PIVOTAL_API_TOKEN  = os.Getenv("PIVOTAL_API_TOKEN")
	ACTUAL_PROJECT_IDS = []int{2077921, 2123086, 1974589, 2099793}
)

func TestGetStoriesByLabel(t *testing.T) {
	// TEST_ACTUAL_PROJECT_ID := 2077921
	// TEST_ACTUAL_STORY_ID := 154266846
	TestActualStoryName := "SmartShare window yes/no buttons not reacting (only the X button reacts)"
	TestActualStoryLabels := []string{"rc-2018-01-15", "smartshare"}
	client := tracker.NewClient(PIVOTAL_API_TOKEN)
	var stories []tracker.Story
	for _, projId := range ACTUAL_PROJECT_IDS {
		projClient := client.InProject(projId)
		projStories, err := GetStoriesWithLabel(projClient, "rc-2018-01-15")
		assert.Equal(t, err, nil, "err should be nil")
		stories = append(stories, projStories...)
	}
	assert.Equal(t, 1, len(stories), "len(stories) should be 1")
	assert.Equal(t, TestActualStoryName, stories[0].Name)
	for _, label := range stories[0].Labels {
		assert.True(t, label.Name == TestActualStoryLabels[0] || label.Name == TestActualStoryLabels[1])
	}
}

func TestGetPrLinkFromStoryDescription(t *testing.T) {
	ACTUAL_PROJECT_LINK := "https://www.pivotaltracker.com/n/projects/1974589/stories/154134429"
	ExpectedStoryName := "Status filters should update when statuses changes are successful"
	ExpectedPrLink := "https://github.com/axialmarket/axial-FE-app/pull/805"
	urlParts, err := ParseStoryUrl(ACTUAL_PROJECT_LINK)
	assert.Equal(t, err, nil, "err should be nil")
	client := tracker.NewClient(PIVOTAL_API_TOKEN)
	projClient := client.InProject(urlParts["projectId"])
	stories, _, err := projClient.Stories(tracker.StoriesQuery{Filter: []string{fmt.Sprintf("id:%d", urlParts["storyId"])}})
	assert.Equal(t, err, nil, "err should be nil")
	assert.Equal(t, ExpectedStoryName, stories[0].Name)
	prLink := GetPrLinksFromStory(stories[0])[0]
	assert.Equal(t, ExpectedPrLink, prLink)
}
