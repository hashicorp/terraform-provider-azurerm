package virtualmachineimagetemplate

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ImageTemplateLastRunStatus struct {
	EndTime     *string      `json:"endTime,omitempty"`
	Message     *string      `json:"message,omitempty"`
	RunState    *RunState    `json:"runState,omitempty"`
	RunSubState *RunSubState `json:"runSubState,omitempty"`
	StartTime   *string      `json:"startTime,omitempty"`
}

func (o *ImageTemplateLastRunStatus) GetEndTimeAsTime() (*time.Time, error) {
	if o.EndTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.EndTime, "2006-01-02T15:04:05Z07:00")
}

func (o *ImageTemplateLastRunStatus) SetEndTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.EndTime = &formatted
}

func (o *ImageTemplateLastRunStatus) GetStartTimeAsTime() (*time.Time, error) {
	if o.StartTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.StartTime, "2006-01-02T15:04:05Z07:00")
}

func (o *ImageTemplateLastRunStatus) SetStartTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.StartTime = &formatted
}
