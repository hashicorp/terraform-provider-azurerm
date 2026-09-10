package schedule

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ScheduleCreateOrUpdateProperties struct {
	AdvancedSchedule *AdvancedSchedule `json:"advancedSchedule,omitempty"`
	Description      *string           `json:"description,omitempty"`
	ExpiryTime       *string           `json:"expiryTime,omitempty"`
	Frequency        ScheduleFrequency `json:"frequency"`
	Interval         *interface{}      `json:"interval,omitempty"`
	StartTime        string            `json:"startTime"`
	TimeZone         *string           `json:"timeZone,omitempty"`
}

func (o *ScheduleCreateOrUpdateProperties) GetExpiryTimeAsTime() (*time.Time, error) {
	if o.ExpiryTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.ExpiryTime, "2006-01-02T15:04:05Z07:00")
}

func (o *ScheduleCreateOrUpdateProperties) SetExpiryTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.ExpiryTime = &formatted
}

func (o *ScheduleCreateOrUpdateProperties) GetStartTimeAsTime() (*time.Time, error) {
	return dates.ParseAsFormat(&o.StartTime, "2006-01-02T15:04:05Z07:00")
}

func (o *ScheduleCreateOrUpdateProperties) SetStartTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.StartTime = formatted
}
