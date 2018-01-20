package format

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"testing"

	"github.com/axialmarket/releaser/notes"
	"github.com/stretchr/testify/assert"
)

var (
	testJsonData  = []byte(`{"release_date":"2018-01-20","projects":[{"name":"Axial Outreach Team","stories":[{"url":"https://www.pivotaltracker.com/story/show/123456789","name":"Test Story","repo":"testrepo","pr_url":"https://github.com/testorg/testrepo/pull/123"},{"url":"https://www.pivotaltracker.com/story/show/234567890","name":"Test Story 2","repo":"testrepo2","pr_url":"https://github.com/testorg/testrepo2/pull/234"}]}]}`)
	testRnVars = notes.RelNotesVars{
		ReleaseDate: "2018-01-20",
		Projects: []notes.RelNotesProject{
			{
				ProjectName: "Axial Outreach Team",
				ProjectStories: []notes.RelNotesStory{
					{
						StoryLink: "https://www.pivotaltracker.com/story/show/123456789",
						StoryName: "Test Story",
						StoryPrLink: "https://github.com/testorg/testrepo/pull/123",
						StoryRepo: "testrepo",
					},
					{
						StoryLink: "https://www.pivotaltracker.com/story/show/234567890",
						StoryName: "Test Story 2",
						StoryPrLink: "https://github.com/testorg/testrepo2/pull/234",
						StoryRepo: "testrepo2",
					},
				},
			},
		},
	}
	testTemplate = []byte(`Release Date: {{ .ReleaseDate }}
{{range .Projects}}{{"{{RelNotesTeam|"}}{{ .ProjectName }}|{{len .ProjectStories }}{{"}}"}}
{{range .ProjectStories}}{{"{{RelNotesTicket|"}}{{ .StoryLink }}|{{ .StoryName }}|{{ .StoryRepo }}|{{ .StoryPrLink }}{{"}}"}}
{{end}}
{{- end}}`)
)

func TestMarshalJsonDataToRelNotesVars(t *testing.T) {
	var rnVars notes.RelNotesVars
	if err := json.Unmarshal(testJsonData, &rnVars); err != nil {
		t.Error(err)
	}
	assert.EqualValues(t, testRnVars, rnVars)
}

func TestReadRelNotesVars(t *testing.T) {
	var rnVars notes.RelNotesVars
	r := bytes.NewReader(testJsonData)
	if err := ReadRelNotesVars(r, &rnVars); err != nil {
		t.Error(err)
	}
	assert.EqualValues(t, testRnVars, rnVars)
}

func TestCompileRelNotesTemplate(t *testing.T) {
	tmpl := bytes.NewReader(testTemplate)
	rnVars := notes.RelNotesVars{
		ReleaseDate: "2018-01-20",
		Projects: []notes.RelNotesProject{
			{
				ProjectName: "Axial Outreach Team",
				ProjectStories: []notes.RelNotesStory{
					{
						StoryLink: "https://www.pivotaltracker.com/story/show/123456789",
						StoryName: "Test Story",
						StoryPrLink: "https://github.com/testorg/testrepo/pull/123",
						StoryRepo: "testrepo",
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
{{RelNotesTicket|https://www.pivotaltracker.com/story/show/123456789|Test Story|testrepo|https://github.com/testorg/testrepo/pull/123}}
`)
	assert.Equal(t, expected, actual)
}
