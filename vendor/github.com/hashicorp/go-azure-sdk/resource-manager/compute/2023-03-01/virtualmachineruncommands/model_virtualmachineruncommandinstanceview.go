package virtualmachineruncommands

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualMachineRunCommandInstanceView struct {
	EndTime          *string               `json:"endTime,omitempty"`
	Error            *string               `json:"error,omitempty"`
	ExecutionMessage *string               `json:"executionMessage,omitempty"`
	ExecutionState   *ExecutionState       `json:"executionState,omitempty"`
	ExitCode         *int64                `json:"exitCode,omitempty"`
	Output           *string               `json:"output,omitempty"`
	StartTime        *string               `json:"startTime,omitempty"`
	Statuses         *[]InstanceViewStatus `json:"statuses,omitempty"`
}

func (o *VirtualMachineRunCommandInstanceView) GetEndTimeAsTime() (*time.Time, error) {
	if o.EndTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.EndTime, "2006-01-02T15:04:05Z07:00")
}

func (o *VirtualMachineRunCommandInstanceView) SetEndTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.EndTime = &formatted
}

func (o *VirtualMachineRunCommandInstanceView) GetStartTimeAsTime() (*time.Time, error) {
	if o.StartTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.StartTime, "2006-01-02T15:04:05Z07:00")
}

func (o *VirtualMachineRunCommandInstanceView) SetStartTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.StartTime = &formatted
}
