package serviceresource

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DatabaseSummaryResult struct {
	EndedOn             *string         `json:"endedOn,omitempty"`
	ErrorPrefix         *string         `json:"errorPrefix,omitempty"`
	ItemsCompletedCount *int64          `json:"itemsCompletedCount,omitempty"`
	ItemsCount          *int64          `json:"itemsCount,omitempty"`
	Name                *string         `json:"name,omitempty"`
	ResultPrefix        *string         `json:"resultPrefix,omitempty"`
	SizeMB              *float64        `json:"sizeMB,omitempty"`
	StartedOn           *string         `json:"startedOn,omitempty"`
	State               *MigrationState `json:"state,omitempty"`
	StatusMessage       *string         `json:"statusMessage,omitempty"`
}

func (o *DatabaseSummaryResult) GetEndedOnAsTime() (*time.Time, error) {
	if o.EndedOn == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.EndedOn, "2006-01-02T15:04:05Z07:00")
}

func (o *DatabaseSummaryResult) SetEndedOnAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.EndedOn = &formatted
}

func (o *DatabaseSummaryResult) GetStartedOnAsTime() (*time.Time, error) {
	if o.StartedOn == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.StartedOn, "2006-01-02T15:04:05Z07:00")
}

func (o *DatabaseSummaryResult) SetStartedOnAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.StartedOn = &formatted
}
