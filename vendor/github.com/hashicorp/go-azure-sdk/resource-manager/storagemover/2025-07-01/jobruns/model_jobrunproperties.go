package jobruns

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type JobRunProperties struct {
	AgentName               *string            `json:"agentName,omitempty"`
	AgentResourceId         *string            `json:"agentResourceId,omitempty"`
	BytesExcluded           *int64             `json:"bytesExcluded,omitempty"`
	BytesFailed             *int64             `json:"bytesFailed,omitempty"`
	BytesNoTransferNeeded   *int64             `json:"bytesNoTransferNeeded,omitempty"`
	BytesScanned            *int64             `json:"bytesScanned,omitempty"`
	BytesTransferred        *int64             `json:"bytesTransferred,omitempty"`
	BytesUnsupported        *int64             `json:"bytesUnsupported,omitempty"`
	Error                   *JobRunError       `json:"error,omitempty"`
	ExecutionEndTime        *string            `json:"executionEndTime,omitempty"`
	ExecutionStartTime      *string            `json:"executionStartTime,omitempty"`
	ItemsExcluded           *int64             `json:"itemsExcluded,omitempty"`
	ItemsFailed             *int64             `json:"itemsFailed,omitempty"`
	ItemsNoTransferNeeded   *int64             `json:"itemsNoTransferNeeded,omitempty"`
	ItemsScanned            *int64             `json:"itemsScanned,omitempty"`
	ItemsTransferred        *int64             `json:"itemsTransferred,omitempty"`
	ItemsUnsupported        *int64             `json:"itemsUnsupported,omitempty"`
	JobDefinitionProperties *interface{}       `json:"jobDefinitionProperties,omitempty"`
	LastStatusUpdate        *string            `json:"lastStatusUpdate,omitempty"`
	ProvisioningState       *ProvisioningState `json:"provisioningState,omitempty"`
	ScanStatus              *JobRunScanStatus  `json:"scanStatus,omitempty"`
	SourceName              *string            `json:"sourceName,omitempty"`
	SourceProperties        *interface{}       `json:"sourceProperties,omitempty"`
	SourceResourceId        *string            `json:"sourceResourceId,omitempty"`
	Status                  *JobRunStatus      `json:"status,omitempty"`
	TargetName              *string            `json:"targetName,omitempty"`
	TargetProperties        *interface{}       `json:"targetProperties,omitempty"`
	TargetResourceId        *string            `json:"targetResourceId,omitempty"`
}

func (o *JobRunProperties) GetExecutionEndTimeAsTime() (*time.Time, error) {
	if o.ExecutionEndTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.ExecutionEndTime, "2006-01-02T15:04:05Z07:00")
}

func (o *JobRunProperties) SetExecutionEndTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.ExecutionEndTime = &formatted
}

func (o *JobRunProperties) GetExecutionStartTimeAsTime() (*time.Time, error) {
	if o.ExecutionStartTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.ExecutionStartTime, "2006-01-02T15:04:05Z07:00")
}

func (o *JobRunProperties) SetExecutionStartTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.ExecutionStartTime = &formatted
}

func (o *JobRunProperties) GetLastStatusUpdateAsTime() (*time.Time, error) {
	if o.LastStatusUpdate == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastStatusUpdate, "2006-01-02T15:04:05Z07:00")
}

func (o *JobRunProperties) SetLastStatusUpdateAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastStatusUpdate = &formatted
}
