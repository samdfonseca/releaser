package format

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"testing"

	"github.com/axialmarket/releaser/notes"
	"github.com/stretchr/testify/assert"
	"github.com/davecgh/go-spew/spew"
)

var (
	testStory = notes.RelNotesStory{
		StoryLink: "https://www.pivotaltracker.com/story/show/123456789",
		StoryName: "Test Story",
		StoryPrLinks: []string{
			"https://github.com/testorg/testrepo/pull/123",
			"https://github.com/testorg/testrepo/pull/456",
		},
	}
	testStoryBytes = []byte(fmt.Sprintf(`{
		"url": "%s",
        "name": "%s",
		"pr_urls": ["%s", "%s"]
	}`, testStory.StoryLink, testStory.StoryName, testStory.StoryPrLinks[0], testStory.StoryPrLinks[1]))


	testRnVars = notes.RelNotesVars{
		ReleaseDate: "2018-01-20",
		Projects: []notes.RelNotesProject{
			{
				ProjectName: "Axial Outreach Team",
				ProjectStories: []notes.RelNotesStory{
					{
						StoryLink:   "https://www.pivotaltracker.com/story/show/123456789",
						StoryName:   "Test Story",
						StoryPrLinks: []string{
							"https://github.com/testorg/testrepo/pull/123",
							"https://github.com/testorg/testrepo/pull/456",
						},
					},
					{
						StoryLink:   "https://www.pivotaltracker.com/story/show/234567890",
						StoryName:   "Test Story 2",
						StoryPrLinks: []string{
							"https://github.com/testorg/testrepo2/pull/234",
						},
					},
				},
			},
		},
	}

	testRnVarsBytes = []byte(fmt.Sprintf(`{
		"release_date": "%s",
		"projects": [
		{
			"name": "%s",
			"stories": [
			{
				"url": "%s",
				"name": "%s",
				"pr_urls": [
					"%s",
					"%s"
				]
			},
			{
				"url": "%s",
				"name": "%s",
          		"pr_urls": ["%s"]
			}
			]
		}
		]
	}`,
	testRnVars.ReleaseDate,
	testRnVars.Projects[0].ProjectName,
	testRnVars.Projects[0].ProjectStories[0].StoryLink,
	testRnVars.Projects[0].ProjectStories[0].StoryName,
	testRnVars.Projects[0].ProjectStories[0].StoryPrLinks[0],
	testRnVars.Projects[0].ProjectStories[0].StoryPrLinks[1],
	testRnVars.Projects[0].ProjectStories[1].StoryLink,
	testRnVars.Projects[0].ProjectStories[1].StoryName,
	testRnVars.Projects[0].ProjectStories[1].StoryPrLinks[0],
	))

	testTemplate = []byte(`Release Date: {{ .ReleaseDate }}
{{range .Projects}}{{"{{RelNotesTeam|"}}{{ .ProjectName }}|{{len .ProjectStories }}{{"}}"}}
{{range .ProjectStories}}{{"{{RelNotesTicket|"}}{{ .StoryLink }}|{{ .StoryName }}|{{ .StoryRepo }}|{{ .StoryPrLink }}{{"}}"}}
{{end}}
{{- end}}`)
)

func TestMarshalJsonDataToStory(t *testing.T) {
	var actualStory notes.RelNotesStory
	expectedStory := testStory
	if err := json.Unmarshal(testStoryBytes, &actualStory); err != nil {
		t.Error(err)
	}
	assert.EqualValues(t, expectedStory, actualStory)
}

func TestMarshalJsonDataToRelNotesVars(t *testing.T) {
	var rnVars notes.RelNotesVars
	if err := json.Unmarshal(testRnVarsBytes, &rnVars); err != nil {
		t.Error(err)
	}
	for i := range testRnVars.Projects {
		expectedProject := testRnVars.Projects[i]
		actualProject := rnVars.Projects[i]
		assert.EqualValues(t, expectedProject.ProjectName, actualProject.ProjectName)
		for j := range expectedProject.ProjectStories {
			expectedStory := expectedProject.ProjectStories[j]
			actualStory := expectedProject.ProjectStories[j]
			assert.EqualValues(t, expectedStory.StoryLink, actualStory.StoryLink)
			assert.EqualValues(t, expectedStory.StoryName, actualStory.StoryName)
			assert.EqualValues(t, expectedStory.StoryPrLinks, actualStory.StoryPrLinks)
		}
	}
	assert.EqualValues(t, testRnVars, rnVars, spew.Sdump(testRnVars, rnVars))
}

func TestReadRelNotesVars(t *testing.T) {
	var rnVars notes.RelNotesVars
	r := bytes.NewReader(testRnVarsBytes)
	if err := ReadRelNotesVars(r, &rnVars); err != nil {
		t.Error(err)
	}
	assert.EqualValues(t, testRnVars, rnVars)
}

func TestCompileRelNotesTemplate(t *testing.T) {
	tmpl := bytes.NewReader(DefaultRelNotesTemplate)
	rnVars := notes.RelNotesVars{
		ReleaseDate: "2018-01-20",
		Projects: []notes.RelNotesProject{
			{
				ProjectName: "Axial Outreach Team",
				ProjectStories: []notes.RelNotesStory{
					{
						StoryLink:   "https://www.pivotaltracker.com/story/show/123456789",
						StoryName:   "Test Story",
						StoryPrLinks: []string{
							"https://github.com/testorg/testrepo/pull/123",
							"https://github.com/testorg/testrepo/pull/234",
						},
					},
				},
			},
		},
	}
	r, w := io.Pipe()
	defer func() {
		r.Close()
	}()
	go func() {
		defer w.Close()

		if err := CompileRelNotesTemplate(rnVars, tmpl, w); err != nil {
			t.Error(err)
		}
	}()
	actual, err := ioutil.ReadAll(r)
	if err != nil {
		t.Error(err)
	}
	expected := []byte(`Release Date: 2018-01-20
{{RelNotesTeam|Axial Outreach Team|1}}
{{RelNotesTicket|https://www.pivotaltracker.com/story/show/123456789|Test Story}}
* https://github.com/testorg/testrepo/pull/123
* https://github.com/testorg/testrepo/pull/234
`)
	assert.Equal(t, expected, actual)
}
