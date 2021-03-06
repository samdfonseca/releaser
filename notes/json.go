package notes

type RelNotesStory struct {
	StoryLink   string `json:"url"`
	StoryName   string `json:"name"`
	StoryRepo   string `json:"repo"`
	StoryPrLink string `json:"pr_url"`
}

type RelNotesProject struct {
	ProjectName    string `json:"name"`
	ProjectStories []RelNotesStory `json:"stories"`
}

type RelNotesVars struct {
	ReleaseDate string            `json:"release_date"`
	Projects    []RelNotesProject `json:"projects"`
}
