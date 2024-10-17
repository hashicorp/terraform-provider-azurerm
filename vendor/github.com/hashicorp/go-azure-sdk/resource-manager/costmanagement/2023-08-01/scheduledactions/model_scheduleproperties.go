package scheduledactions

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ScheduleProperties struct {
	DayOfMonth   *int64            `json:"dayOfMonth,omitempty"`
	DaysOfWeek   *[]DaysOfWeek     `json:"daysOfWeek,omitempty"`
	EndDate      string            `json:"endDate"`
	Frequency    ScheduleFrequency `json:"frequency"`
	HourOfDay    *int64            `json:"hourOfDay,omitempty"`
	StartDate    string            `json:"startDate"`
	WeeksOfMonth *[]WeeksOfMonth   `json:"weeksOfMonth,omitempty"`
}

func (o *ScheduleProperties) GetEndDateAsTime() (*time.Time, error) {
	return dates.ParseAsFormat(&o.EndDate, "2006-01-02T15:04:05Z07:00")
}

func (o *ScheduleProperties) SetEndDateAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.EndDate = formatted
}

func (o *ScheduleProperties) GetStartDateAsTime() (*time.Time, error) {
	return dates.ParseAsFormat(&o.StartDate, "2006-01-02T15:04:05Z07:00")
}

func (o *ScheduleProperties) SetStartDateAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.StartDate = formatted
}
