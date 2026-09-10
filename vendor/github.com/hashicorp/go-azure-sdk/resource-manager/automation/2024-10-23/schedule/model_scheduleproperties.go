package schedule

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ScheduleProperties struct {
	AdvancedSchedule        *AdvancedSchedule  `json:"advancedSchedule,omitempty"`
	CreationTime            *string            `json:"creationTime,omitempty"`
	Description             *string            `json:"description,omitempty"`
	ExpiryTime              *string            `json:"expiryTime,omitempty"`
	ExpiryTimeOffsetMinutes *float64           `json:"expiryTimeOffsetMinutes,omitempty"`
	Frequency               *ScheduleFrequency `json:"frequency,omitempty"`
	Interval                *interface{}       `json:"interval,omitempty"`
	IsEnabled               *bool              `json:"isEnabled,omitempty"`
	LastModifiedTime        *string            `json:"lastModifiedTime,omitempty"`
	NextRun                 *string            `json:"nextRun,omitempty"`
	NextRunOffsetMinutes    *float64           `json:"nextRunOffsetMinutes,omitempty"`
	StartTime               *string            `json:"startTime,omitempty"`
	StartTimeOffsetMinutes  *float64           `json:"startTimeOffsetMinutes,omitempty"`
	TimeZone                *string            `json:"timeZone,omitempty"`
}

func (o *ScheduleProperties) GetCreationTimeAsTime() (*time.Time, error) {
	if o.CreationTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreationTime, "2006-01-02T15:04:05Z07:00")
}

func (o *ScheduleProperties) SetCreationTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreationTime = &formatted
}

func (o *ScheduleProperties) GetExpiryTimeAsTime() (*time.Time, error) {
	if o.ExpiryTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.ExpiryTime, "2006-01-02T15:04:05Z07:00")
}

func (o *ScheduleProperties) SetExpiryTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.ExpiryTime = &formatted
}

func (o *ScheduleProperties) GetLastModifiedTimeAsTime() (*time.Time, error) {
	if o.LastModifiedTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastModifiedTime, "2006-01-02T15:04:05Z07:00")
}

func (o *ScheduleProperties) SetLastModifiedTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastModifiedTime = &formatted
}

func (o *ScheduleProperties) GetNextRunAsTime() (*time.Time, error) {
	if o.NextRun == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.NextRun, "2006-01-02T15:04:05Z07:00")
}

func (o *ScheduleProperties) SetNextRunAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.NextRun = &formatted
}

func (o *ScheduleProperties) GetStartTimeAsTime() (*time.Time, error) {
	if o.StartTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.StartTime, "2006-01-02T15:04:05Z07:00")
}

func (o *ScheduleProperties) SetStartTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.StartTime = &formatted
}
