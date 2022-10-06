package machinelearningcomputes

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ComputeInstanceLastOperation struct {
	OperationName    *OperationName    `json:"operationName,omitempty"`
	OperationStatus  *OperationStatus  `json:"operationStatus,omitempty"`
	OperationTime    *string           `json:"operationTime,omitempty"`
	OperationTrigger *OperationTrigger `json:"operationTrigger,omitempty"`
}

func (o *ComputeInstanceLastOperation) GetOperationTimeAsTime() (*time.Time, error) {
	if o.OperationTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.OperationTime, "2006-01-02T15:04:05Z07:00")
}

func (o *ComputeInstanceLastOperation) SetOperationTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.OperationTime = &formatted
}
