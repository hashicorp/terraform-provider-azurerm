package workflowtriggers

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WorkflowTriggerProperties struct {
	ChangedTime       *string                           `json:"changedTime,omitempty"`
	CreatedTime       *string                           `json:"createdTime,omitempty"`
	LastExecutionTime *string                           `json:"lastExecutionTime,omitempty"`
	NextExecutionTime *string                           `json:"nextExecutionTime,omitempty"`
	ProvisioningState *WorkflowTriggerProvisioningState `json:"provisioningState,omitempty"`
	Recurrence        *WorkflowTriggerRecurrence        `json:"recurrence,omitempty"`
	State             *WorkflowState                    `json:"state,omitempty"`
	Status            *WorkflowStatus                   `json:"status,omitempty"`
	Workflow          *ResourceReference                `json:"workflow,omitempty"`
}

func (o *WorkflowTriggerProperties) GetChangedTimeAsTime() (*time.Time, error) {
	if o.ChangedTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.ChangedTime, "2006-01-02T15:04:05Z07:00")
}

func (o *WorkflowTriggerProperties) SetChangedTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.ChangedTime = &formatted
}

func (o *WorkflowTriggerProperties) GetCreatedTimeAsTime() (*time.Time, error) {
	if o.CreatedTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreatedTime, "2006-01-02T15:04:05Z07:00")
}

func (o *WorkflowTriggerProperties) SetCreatedTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreatedTime = &formatted
}

func (o *WorkflowTriggerProperties) GetLastExecutionTimeAsTime() (*time.Time, error) {
	if o.LastExecutionTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastExecutionTime, "2006-01-02T15:04:05Z07:00")
}

func (o *WorkflowTriggerProperties) SetLastExecutionTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastExecutionTime = &formatted
}

func (o *WorkflowTriggerProperties) GetNextExecutionTimeAsTime() (*time.Time, error) {
	if o.NextExecutionTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.NextExecutionTime, "2006-01-02T15:04:05Z07:00")
}

func (o *WorkflowTriggerProperties) SetNextExecutionTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.NextExecutionTime = &formatted
}
