package job

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type JobProperties struct {
	CreationTime           *string                     `json:"creationTime,omitempty"`
	EndTime                *string                     `json:"endTime,omitempty"`
	Exception              *string                     `json:"exception,omitempty"`
	JobId                  *string                     `json:"jobId,omitempty"`
	LastModifiedTime       *string                     `json:"lastModifiedTime,omitempty"`
	LastStatusModifiedTime *string                     `json:"lastStatusModifiedTime,omitempty"`
	Parameters             *map[string]string          `json:"parameters,omitempty"`
	ProvisioningState      *JobProvisioningState       `json:"provisioningState,omitempty"`
	RunOn                  *string                     `json:"runOn,omitempty"`
	Runbook                *RunbookAssociationProperty `json:"runbook,omitempty"`
	StartTime              *string                     `json:"startTime,omitempty"`
	StartedBy              *string                     `json:"startedBy,omitempty"`
	Status                 *JobStatus                  `json:"status,omitempty"`
	StatusDetails          *string                     `json:"statusDetails,omitempty"`
}

func (o *JobProperties) GetCreationTimeAsTime() (*time.Time, error) {
	if o.CreationTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreationTime, "2006-01-02T15:04:05Z07:00")
}

func (o *JobProperties) SetCreationTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreationTime = &formatted
}

func (o *JobProperties) GetEndTimeAsTime() (*time.Time, error) {
	if o.EndTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.EndTime, "2006-01-02T15:04:05Z07:00")
}

func (o *JobProperties) SetEndTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.EndTime = &formatted
}

func (o *JobProperties) GetLastModifiedTimeAsTime() (*time.Time, error) {
	if o.LastModifiedTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastModifiedTime, "2006-01-02T15:04:05Z07:00")
}

func (o *JobProperties) SetLastModifiedTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastModifiedTime = &formatted
}

func (o *JobProperties) GetLastStatusModifiedTimeAsTime() (*time.Time, error) {
	if o.LastStatusModifiedTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastStatusModifiedTime, "2006-01-02T15:04:05Z07:00")
}

func (o *JobProperties) SetLastStatusModifiedTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastStatusModifiedTime = &formatted
}

func (o *JobProperties) GetStartTimeAsTime() (*time.Time, error) {
	if o.StartTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.StartTime, "2006-01-02T15:04:05Z07:00")
}

func (o *JobProperties) SetStartTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.StartTime = &formatted
}
