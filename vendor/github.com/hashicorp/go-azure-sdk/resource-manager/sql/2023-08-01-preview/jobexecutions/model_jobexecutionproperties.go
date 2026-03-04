package jobexecutions

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type JobExecutionProperties struct {
	CreateTime              *string                `json:"createTime,omitempty"`
	CurrentAttemptStartTime *string                `json:"currentAttemptStartTime,omitempty"`
	CurrentAttempts         *int64                 `json:"currentAttempts,omitempty"`
	EndTime                 *string                `json:"endTime,omitempty"`
	JobExecutionId          *string                `json:"jobExecutionId,omitempty"`
	JobVersion              *int64                 `json:"jobVersion,omitempty"`
	LastMessage             *string                `json:"lastMessage,omitempty"`
	Lifecycle               *JobExecutionLifecycle `json:"lifecycle,omitempty"`
	ProvisioningState       *ProvisioningState     `json:"provisioningState,omitempty"`
	StartTime               *string                `json:"startTime,omitempty"`
	StepId                  *int64                 `json:"stepId,omitempty"`
	StepName                *string                `json:"stepName,omitempty"`
	Target                  *JobExecutionTarget    `json:"target,omitempty"`
}

func (o *JobExecutionProperties) GetCreateTimeAsTime() (*time.Time, error) {
	if o.CreateTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreateTime, "2006-01-02T15:04:05Z07:00")
}

func (o *JobExecutionProperties) SetCreateTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreateTime = &formatted
}

func (o *JobExecutionProperties) GetCurrentAttemptStartTimeAsTime() (*time.Time, error) {
	if o.CurrentAttemptStartTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CurrentAttemptStartTime, "2006-01-02T15:04:05Z07:00")
}

func (o *JobExecutionProperties) SetCurrentAttemptStartTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CurrentAttemptStartTime = &formatted
}

func (o *JobExecutionProperties) GetEndTimeAsTime() (*time.Time, error) {
	if o.EndTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.EndTime, "2006-01-02T15:04:05Z07:00")
}

func (o *JobExecutionProperties) SetEndTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.EndTime = &formatted
}

func (o *JobExecutionProperties) GetStartTimeAsTime() (*time.Time, error) {
	if o.StartTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.StartTime, "2006-01-02T15:04:05Z07:00")
}

func (o *JobExecutionProperties) SetStartTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.StartTime = &formatted
}
