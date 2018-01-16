package ptracker

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/samdfonseca/go-tracker"
)

func GetPrLinkRegexp(githubOrg string) (*regexp.Regexp, error) {
	return regexp.Compile(fmt.Sprintf(`https:\/\/github.com\/%s\/[a-zA-Z-]+\/pull/[0-9]+`, githubOrg))
}

func NewClient(apiToken string) *tracker.Client {
	return tracker.NewClient(apiToken)
}

func GetStoriesWithLabelQuery(label string) tracker.StoriesQuery {
	return tracker.StoriesQuery{
		Label: label,
	}
}

func getStoriesWithLabel(projClient tracker.ProjectClient, label string, limit, offset int) ([]tracker.Story, error) {
	query := GetStoriesWithLabelQuery(label)
	query.Limit = limit
	query.Offset = offset
	stories, pagination, err := projClient.Stories(query)
	if err != nil {
		return nil, err
	}
	if len(stories) != pagination.Total {
		for o := pagination.Offset; len(stories) >= pagination.Total; o = o + pagination.Limit {
			nextStories, err := getStoriesWithLabel(projClient, label, limit, o)
			if err != nil {
				return nil, err
			}
			stories = append(stories, nextStories...)
		}
	}
	return stories, err
}

func GetStoriesWithLabel(projClient tracker.ProjectClient, label string) ([]tracker.Story, error) {
	return getStoriesWithLabel(projClient, label, 25, 0)
}

func ParseStoryUrl(storyUrl string) (map[string]int, error) {
	parsedUrl, err := url.Parse(storyUrl)
	if err != nil {
		return nil, err
	}
	urlPath := parsedUrl.EscapedPath()
	urlPathRegexp := regexp.MustCompile(`^\/n\/projects\/[0-9]+\/stories/[0-9]+$`)
	if !urlPathRegexp.MatchString(urlPath) {
		return nil, fmt.Errorf("Url path %s does not match regex %s", urlPath, urlPathRegexp.String())
	}
	urlPathParts := strings.Split(urlPath, "/")
	projectId, err := strconv.ParseInt(urlPathParts[3], 10, 0)
	if err != nil {
		return nil, err
	}
	storyId, err := strconv.ParseInt(urlPathParts[5], 10, 0)
	if err != nil {
		return nil, err
	}
	return map[string]int{
		"projectId": int(projectId),
		"storyId":   int(storyId),
	}, nil
}

func GetPrLinksFromStory(story tracker.Story, prLinkRegexp *regexp.Regexp) []string {
	return prLinkRegexp.FindAllString(story.Description, -1)
}
