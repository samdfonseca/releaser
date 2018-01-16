package wiki

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	wikiUrl       = "https://wiki.axialmarket.com/api.php"
	testPageTitle = "Samtest"
)

func TestCreateWikiClient(t *testing.T) {
	client, err := NewWikiClient(wikiUrl)
	assert.Equal(t, err, nil, "Error should be nil")
	assert.NotEqual(t, client, nil, "Client should not be nil")
}

func TestEditWikiPage(t *testing.T) {
	client, err := NewWikiClient(wikiUrl)
	assert.Equal(t, err, nil, "Error should be nil")
	newPageText := fmt.Sprintf("sam test %d", time.Now().Unix())
	err = client.UpdatePageText(testPageTitle, newPageText)
	assert.Equal(t, err, nil, "Error should be nil")
	pageText, err := client.ReadPage(testPageTitle)
	assert.Equal(t, err, nil, "Error should be nil")
	assert.Equal(t, newPageText, pageText, "Page should be updated with new text")
}

func TestGenerateRelNotesText(t *testing.T) {
	EXPECTED_REL_NOTES_PAGE := `{{RelNotesHeader}}
{{RelNotesTeam|OUTREACH|1}}
{{RelNotesTicket|https://www.pivotaltracker.com/story/show/154134429|Status filters should update when statuses changes are successful|axial-fe-app|https://github.com/axialmarket/axial-FE-app/pull/805}}
`
	i := RelNotesVars{
		Teams: []RelNotesTeam{
			RelNotesTeam{
				TeamName: "OUTREACH",
				TeamItems: []RelNotesItem{
					RelNotesItem{
						StoryLink:   "https://www.pivotaltracker.com/story/show/154134429",
						StoryName:   "Status filters should update when statuses changes are successful",
						StoryRepo:   "axial-fe-app",
						StoryPrLink: "https://github.com/axialmarket/axial-FE-app/pull/805",
					},
				},
			},
		},
	}
	var buf bytes.Buffer
	err := GenerateRelNotesText(i, &buf)
	assert.Equal(t, err, nil, "Error should be nil")
	s, err := ioutil.ReadAll(&buf)
	assert.Equal(t, err, nil, "Error should be nil")
	assert.Equal(t, EXPECTED_REL_NOTES_PAGE, string(s))
}

func TestGenerateRelNotesMultipleStories(t *testing.T) {
	EXPECTED_REL_NOTES_PAGE := `{{RelNotesHeader}}
{{RelNotesTeam|OUTREACH|2}}
{{RelNotesTicket|https://www.pivotaltracker.com/story/show/154134429|Status filters should update when statuses changes are successful|axial-fe-app|https://github.com/axialmarket/axial-FE-app/pull/805}}
{{RelNotesTicket|https://www.pivotaltracker.com/story/show/154134398|I want to see sender's name in Active Deals list|axial-fe-app|https://github.com/axialmarket/axial-FE-app/pull/813}}
{{RelNotesTeam|TAXONIMY|1}}
{{RelNotesTicket|https://www.pivotaltracker.com/story/show/154255967|Implement discard changes option to exit modal on sellside form|axial-fe-app|https://github.com/axialmarket/axial-FE-app/pull/814}}
`
	i := RelNotesVars{
		Teams: []RelNotesTeam{
			RelNotesTeam{
				TeamName: "OUTREACH",
				TeamItems: []RelNotesItem{
					RelNotesItem{
						StoryLink:   "https://www.pivotaltracker.com/story/show/154134429",
						StoryName:   "Status filters should update when statuses changes are successful",
						StoryRepo:   "axial-fe-app",
						StoryPrLink: "https://github.com/axialmarket/axial-FE-app/pull/805",
					},
					RelNotesItem{
						StoryLink:   "https://www.pivotaltracker.com/story/show/154134398",
						StoryName:   "I want to see sender's name in Active Deals list",
						StoryRepo:   "axial-fe-app",
						StoryPrLink: "https://github.com/axialmarket/axial-FE-app/pull/813",
					},
				},
			},
			RelNotesTeam{
				TeamName: "TAXONIMY",
				TeamItems: []RelNotesItem{
					RelNotesItem{
						StoryLink:   "https://www.pivotaltracker.com/story/show/154255967",
						StoryName:   "Implement discard changes option to exit modal on sellside form",
						StoryRepo:   "axial-fe-app",
						StoryPrLink: "https://github.com/axialmarket/axial-FE-app/pull/814",
					},
				},
			},
		},
	}
	var buf bytes.Buffer
	err := GenerateRelNotesText(i, &buf)
	assert.Equal(t, err, nil, "Error should be nil")
	s, err := ioutil.ReadAll(&buf)
	assert.Equal(t, err, nil, "Error should be nil")
	assert.Equal(t, EXPECTED_REL_NOTES_PAGE, string(s))
}
