package protectionpolicies

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type HourlySchedule struct {
	Interval                *int64  `json:"interval,omitempty"`
	ScheduleWindowDuration  *int64  `json:"scheduleWindowDuration,omitempty"`
	ScheduleWindowStartTime *string `json:"scheduleWindowStartTime,omitempty"`
}

func (o *HourlySchedule) GetScheduleWindowStartTimeAsTime() (*time.Time, error) {
	if o.ScheduleWindowStartTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.ScheduleWindowStartTime, "2006-01-02T15:04:05Z07:00")
}

func (o *HourlySchedule) SetScheduleWindowStartTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.ScheduleWindowStartTime = &formatted
}
