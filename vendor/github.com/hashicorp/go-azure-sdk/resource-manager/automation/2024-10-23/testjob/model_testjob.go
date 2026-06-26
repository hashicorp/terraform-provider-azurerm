package testjob

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TestJob struct {
	CreationTime           *string            `json:"creationTime,omitempty"`
	EndTime                *string            `json:"endTime,omitempty"`
	Exception              *string            `json:"exception,omitempty"`
	LastModifiedTime       *string            `json:"lastModifiedTime,omitempty"`
	LastStatusModifiedTime *string            `json:"lastStatusModifiedTime,omitempty"`
	LogActivityTrace       *int64             `json:"logActivityTrace,omitempty"`
	Parameters             *map[string]string `json:"parameters,omitempty"`
	RunOn                  *string            `json:"runOn,omitempty"`
	StartTime              *string            `json:"startTime,omitempty"`
	Status                 *string            `json:"status,omitempty"`
	StatusDetails          *string            `json:"statusDetails,omitempty"`
}

func (o *TestJob) GetCreationTimeAsTime() (*time.Time, error) {
	if o.CreationTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreationTime, "2006-01-02T15:04:05Z07:00")
}

func (o *TestJob) SetCreationTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreationTime = &formatted
}

func (o *TestJob) GetEndTimeAsTime() (*time.Time, error) {
	if o.EndTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.EndTime, "2006-01-02T15:04:05Z07:00")
}

func (o *TestJob) SetEndTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.EndTime = &formatted
}

func (o *TestJob) GetLastModifiedTimeAsTime() (*time.Time, error) {
	if o.LastModifiedTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastModifiedTime, "2006-01-02T15:04:05Z07:00")
}

func (o *TestJob) SetLastModifiedTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastModifiedTime = &formatted
}

func (o *TestJob) GetLastStatusModifiedTimeAsTime() (*time.Time, error) {
	if o.LastStatusModifiedTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastStatusModifiedTime, "2006-01-02T15:04:05Z07:00")
}

func (o *TestJob) SetLastStatusModifiedTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastStatusModifiedTime = &formatted
}

func (o *TestJob) GetStartTimeAsTime() (*time.Time, error) {
	if o.StartTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.StartTime, "2006-01-02T15:04:05Z07:00")
}

func (o *TestJob) SetStartTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.StartTime = &formatted
}
