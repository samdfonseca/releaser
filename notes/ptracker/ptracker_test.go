package ptracker

//import (
//	"fmt"
//	"os"
//	"testing"
//
//	"github.com/axialmarket/releaser/config"
//	"github.com/samdfonseca/go-tracker"
//	"github.com/stretchr/testify/assert"
//)
//
//var (
//	TEST_CONFIG = config.ReleaserConfig{
//		PivotalApiToken:   os.Getenv("PIVOTAL_API_TOKEN"),
//		PivotalProjectIds: []int{2077921, 2123086, 1974589, 2099793},
//		GithubOrg:         "axialmarket",
//	}
//)

//func TestGetStoriesByLabel(t *testing.T) {
//	TestActualStoryName := "SmartShare window yes/no buttons not reacting (only the X button reacts)"
//	TestActualStoryLabels := []string{"rc-2018-01-15", "smartshare"}
//	client := tracker.NewClient(TEST_CONFIG.PivotalApiToken)
//	var stories []tracker.Story
//	for _, projId := range TEST_CONFIG.PivotalProjectIds {
//		projClient := client.InProject(projId)
//		projStories, err := GetStoriesWithLabel(projClient, "rc-2018-01-15")
//		assert.Equal(t, err, nil, "err should be nil")
//		stories = append(stories, projStories...)
//	}
//	assert.Equal(t, 1, len(stories), "len(stories) should be 1")
//	assert.Equal(t, TestActualStoryName, stories[0].Name)
//	for _, label := range stories[0].Labels {
//		assert.True(t, label.Name == TestActualStoryLabels[0] || label.Name == TestActualStoryLabels[1])
//	}
//}
//
//func TestGetPrLinkFromStoryDescription(t *testing.T) {
//	ACTUAL_PROJECT_LINK := "https://www.pivotaltracker.com/n/projects/1974589/stories/154134429"
//	ExpectedStoryName := "Status filters should update when statuses changes are successful"
//	ExpectedPrLink := "https://github.com/axialmarket/axial-FE-app/pull/805"
//	urlParts, err := ParseStoryUrl(ACTUAL_PROJECT_LINK)
//	assert.Equal(t, err, nil, "err should be nil")
//	client := tracker.NewClient(TEST_CONFIG.PivotalApiToken)
//	projClient := client.InProject(urlParts["projectId"])
//	stories, _, err := projClient.Stories(tracker.StoriesQuery{Filter: []string{fmt.Sprintf("id:%d", urlParts["storyId"])}})
//	assert.Equal(t, err, nil, "err should be nil")
//	assert.Equal(t, ExpectedStoryName, stories[0].Name)
//	prLinkRegexp, err := GetPrLinkRegexp(TEST_CONFIG.GithubOrg)
//	assert.Equal(t, err, nil, "err should be nil")
//	prLink := GetPrLinksFromStory(stories[0], prLinkRegexp)[0]
//	assert.Equal(t, ExpectedPrLink, prLink)
//}
