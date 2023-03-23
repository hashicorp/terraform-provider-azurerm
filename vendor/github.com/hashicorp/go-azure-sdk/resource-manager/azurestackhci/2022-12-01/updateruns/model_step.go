package updateruns

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Step struct {
	Description        *string `json:"description,omitempty"`
	EndTimeUtc         *string `json:"endTimeUtc,omitempty"`
	ErrorMessage       *string `json:"errorMessage,omitempty"`
	LastUpdatedTimeUtc *string `json:"lastUpdatedTimeUtc,omitempty"`
	Name               *string `json:"name,omitempty"`
	StartTimeUtc       *string `json:"startTimeUtc,omitempty"`
	Status             *string `json:"status,omitempty"`
	Steps              *[]Step `json:"steps,omitempty"`
}

func (o *Step) GetEndTimeUtcAsTime() (*time.Time, error) {
	if o.EndTimeUtc == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.EndTimeUtc, "2006-01-02T15:04:05Z07:00")
}

func (o *Step) SetEndTimeUtcAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.EndTimeUtc = &formatted
}

func (o *Step) GetLastUpdatedTimeUtcAsTime() (*time.Time, error) {
	if o.LastUpdatedTimeUtc == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastUpdatedTimeUtc, "2006-01-02T15:04:05Z07:00")
}

func (o *Step) SetLastUpdatedTimeUtcAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastUpdatedTimeUtc = &formatted
}

func (o *Step) GetStartTimeUtcAsTime() (*time.Time, error) {
	if o.StartTimeUtc == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.StartTimeUtc, "2006-01-02T15:04:05Z07:00")
}

func (o *Step) SetStartTimeUtcAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.StartTimeUtc = &formatted
}
