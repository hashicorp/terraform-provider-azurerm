package workflowrunactions

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RetryHistory struct {
	ClientRequestId  *string        `json:"clientRequestId,omitempty"`
	Code             *string        `json:"code,omitempty"`
	EndTime          *string        `json:"endTime,omitempty"`
	Error            *ErrorResponse `json:"error,omitempty"`
	ServiceRequestId *string        `json:"serviceRequestId,omitempty"`
	StartTime        *string        `json:"startTime,omitempty"`
}

func (o *RetryHistory) GetEndTimeAsTime() (*time.Time, error) {
	if o.EndTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.EndTime, "2006-01-02T15:04:05Z07:00")
}

func (o *RetryHistory) SetEndTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.EndTime = &formatted
}

func (o *RetryHistory) GetStartTimeAsTime() (*time.Time, error) {
	if o.StartTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.StartTime, "2006-01-02T15:04:05Z07:00")
}

func (o *RetryHistory) SetStartTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.StartTime = &formatted
}
