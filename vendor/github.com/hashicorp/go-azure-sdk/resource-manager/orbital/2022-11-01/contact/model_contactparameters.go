package contact

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ContactParameters struct {
	ContactProfile    ResourceReference `json:"contactProfile"`
	EndTime           string            `json:"endTime"`
	GroundStationName string            `json:"groundStationName"`
	StartTime         string            `json:"startTime"`
}

func (o *ContactParameters) GetEndTimeAsTime() (*time.Time, error) {
	return dates.ParseAsFormat(&o.EndTime, "2006-01-02T15:04:05Z07:00")
}

func (o *ContactParameters) SetEndTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.EndTime = formatted
}

func (o *ContactParameters) GetStartTimeAsTime() (*time.Time, error) {
	return dates.ParseAsFormat(&o.StartTime, "2006-01-02T15:04:05Z07:00")
}

func (o *ContactParameters) SetStartTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.StartTime = formatted
}
