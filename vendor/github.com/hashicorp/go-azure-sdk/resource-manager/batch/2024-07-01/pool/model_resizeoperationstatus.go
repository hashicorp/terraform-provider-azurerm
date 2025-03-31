package pool

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ResizeOperationStatus struct {
	Errors                 *[]ResizeError                 `json:"errors,omitempty"`
	NodeDeallocationOption *ComputeNodeDeallocationOption `json:"nodeDeallocationOption,omitempty"`
	ResizeTimeout          *string                        `json:"resizeTimeout,omitempty"`
	StartTime              *string                        `json:"startTime,omitempty"`
	TargetDedicatedNodes   *int64                         `json:"targetDedicatedNodes,omitempty"`
	TargetLowPriorityNodes *int64                         `json:"targetLowPriorityNodes,omitempty"`
}

func (o *ResizeOperationStatus) GetStartTimeAsTime() (*time.Time, error) {
	if o.StartTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.StartTime, "2006-01-02T15:04:05Z07:00")
}

func (o *ResizeOperationStatus) SetStartTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.StartTime = &formatted
}
