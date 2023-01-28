package workflowrunactions

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WorkflowRunActionRepetitionProperties struct {
	Code              *string               `json:"code,omitempty"`
	Correlation       *RunActionCorrelation `json:"correlation,omitempty"`
	EndTime           *string               `json:"endTime,omitempty"`
	Error             *interface{}          `json:"error,omitempty"`
	Inputs            *interface{}          `json:"inputs,omitempty"`
	InputsLink        *ContentLink          `json:"inputsLink,omitempty"`
	IterationCount    *int64                `json:"iterationCount,omitempty"`
	Outputs           *interface{}          `json:"outputs,omitempty"`
	OutputsLink       *ContentLink          `json:"outputsLink,omitempty"`
	RepetitionIndexes *[]RepetitionIndex    `json:"repetitionIndexes,omitempty"`
	RetryHistory      *[]RetryHistory       `json:"retryHistory,omitempty"`
	StartTime         *string               `json:"startTime,omitempty"`
	Status            *WorkflowStatus       `json:"status,omitempty"`
	TrackedProperties *interface{}          `json:"trackedProperties,omitempty"`
	TrackingId        *string               `json:"trackingId,omitempty"`
}

func (o *WorkflowRunActionRepetitionProperties) GetEndTimeAsTime() (*time.Time, error) {
	if o.EndTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.EndTime, "2006-01-02T15:04:05Z07:00")
}

func (o *WorkflowRunActionRepetitionProperties) SetEndTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.EndTime = &formatted
}

func (o *WorkflowRunActionRepetitionProperties) GetStartTimeAsTime() (*time.Time, error) {
	if o.StartTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.StartTime, "2006-01-02T15:04:05Z07:00")
}

func (o *WorkflowRunActionRepetitionProperties) SetStartTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.StartTime = &formatted
}
