// Copyright 2016 Christopher Brown. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package tracker

import (
	"fmt"
	"time"
)

type Me Person

type Person struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	Initials string `json:"initials"`
	ID       int    `json:"id"`
	Email    string `json:"email"`
}

type Day string

const (
	DayMonday    Day = "Monday"
	DayTuesday   Day = "Tuesday"
	DayWednesday Day = "Wednesday"
	DayThursday  Day = "Thursday"
	DayFriday    Day = "Friday"
	DaySaturday  Day = "Saturday"
	DaySunday    Day = "Sunday"
)

type Date time.Time

func (date *Date) UnmarshalJSON(content []byte) error {
	s := string(content)

	parsingError := func() error {
		return fmt.Errorf(
			"pivotal.Date.UnmarshalJSON: invalid date string: %s", content)
	}

	// Check whether the leading and trailing " is there.
	if len(s) < 2 || s[0] != '"' || s[len(s)-1] != '"' {
		return parsingError()
	}

	// Strip the leading and trailing "
	s = s[:len(s)-1][1:]

	// Parse the rest.
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return parsingError()
	}

	*date = Date(t)
	return nil
}

func (date Date) MarshalJson() ([]byte, error) {
	return []byte((time.Time)(date).Format("2006-01-02")), nil
}

type TimeZone struct {
	OlsonName string `json:"olson_name,omitempty"`
	Offset    string `json:"offset,omitempty"`
}

type AccountingType string

const (
	AccountingTypeUnbillable AccountingType = "unbillable"
	AccountingTypeBillable   AccountingType = "billable"
	AccountingTypeOverhead   AccountingType = "overhead"
)

type Project struct {
	Id                           int            `json:"id"`
	Name                         string         `json:"name"`
	Version                      int            `json:"version"`
	IterationLength              int            `json:"iteration_length"`
	WeekStartDay                 Day            `json:"week_start_day"`
	PointScale                   string         `json:"point_scale"`
	PointScaleIsCustom           bool           `json:"point_scale_is_custom"`
	BugsAndChoresAreEstimatable  bool           `json:"bugs_and_chores_are_estimatable"`
	AutomaticPlanning            bool           `json:"automatic_planning"`
	EnableTasks                  bool           `json:"enable_tasks"`
	StartDate                    *Date          `json:"start_date"`
	TimeZone                     *TimeZone      `json:"time_zone"`
	VelocityAveragedOver         int            `json:"velocity_averaged_over"`
	ShownIterationsStartTime     *time.Time     `json:"shown_iterations_start_time"`
	StartTime                    *time.Time     `json:"start_time"`
	NumberOfDoneIterationsToShow int            `json:"number_of_done_iterations_to_show"`
	HasGoogleDomain              bool           `json:"has_google_domain"`
	Description                  string         `json:"description"`
	ProfileContent               string         `json:"profile_content"`
	EnableIncomingEmails         bool           `json:"enable_incoming_emails"`
	InitialVelocity              int            `json:"initial_velocity"`
	ProjectType                  string         `json:"project_type"`
	Public                       bool           `json:"public"`
	AtomEnabled                  bool           `json:"atom_enabled"`
	CurrentIterationNumber       int            `json:"current_iteration_number"`
	CurrentVelocity              int            `json:"current_velocity"`
	CurrentVolatility            float64        `json:"current_volatility"`
	AccountId                    int            `json:"account_id"`
	AccountingType               AccountingType `json:"accounting_type"`
	Featured                     bool           `json:"featured"`
	StoryIds                     []int          `json:"story_ids"`
	EpicIds                      []int          `json:"epic_ids"`
	MembershipIds                []int          `json:"membership_ids"`
	LabelIds                     []int          `json:"label_ids"`
	IntegrationIds               []int          `json:"integration_ids"`
	IterationOverrideNumbers     []int          `json:"iteration_override_numbers"`
	CreatedAt                    *time.Time     `json:"created_at"`
	UpdatedAt                    *time.Time     `json:"updated_at"`
}

type Story struct {
	ID        int `json:"id,omitempty"`
	ProjectID int `json:"project_id,omitempty"`

	URL string `json:"url,omitempty"`

	Name        string     `json:"name,omitempty"`
	Description string     `json:"description,omitempty"`
	Type        StoryType  `json:"story_type,omitempty"`
	State       StoryState `json:"current_state,omitempty"`

	Labels []Label `json:"labels,omitempty"`

	CreatedAt  *time.Time `json:"created_at,omitempty"`
	UpdatedAt  *time.Time `json:"updated_at,omitempty"`
	AcceptedAt *time.Time `json:"accepted_at,omitempty"`
}

type Task struct {
	ID      int `json:"id,omitempty"`
	StoryID int `json:"story_id,omitempty"`

	Description string `json:"description,omitempty"`
	IsComplete  bool   `json:"complete,omitempty"`
	Position    int    `json:"position,omitempty"`

	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

type Comment struct {
	Text string `json:"text,omitempty"`
}

type Label struct {
	ID        int `json:"id,omitempty"`
	ProjectID int `json:"project_id,omitempty"`

	Name string `json:"name"`
}

type StoryType string

const (
	StoryTypeFeature = "feature"
	StoryTypeBug     = "bug"
	StoryTypeChore   = "chore"
	StoryTypeRelease = "release"
)

type StoryState string

const (
	StoryStateUnscheduled = "unscheduled"
	StoryStatePlanned     = "planned"
	StoryStateStarted     = "started"
	StoryStateFinished    = "finished"
	StoryStateDelivered   = "delivered"
	StoryStateAccepted    = "accepted"
	StoryStateRejected    = "rejected"
)

type Activity struct {
	Kind             string        `json:"kind"`
	GUID             string        `json:"guid"`
	ProjectVersion   int           `json:"project_version"`
	Message          string        `json:"message"`
	Highlight        string        `json:"highlight"`
	Changes          []interface{} `json:"changes"`
	PrimaryResources []interface{} `json:"primary_resources"`
	Project          interface{}   `json:"project"`
	PerformedBy      interface{}   `json:"performed_by"`
	OccurredAt       time.Time     `json:"occurred_at"`
}

type ProjectMembership struct {
	ID     int `json:"id"`
	Person Person
}
