package machinelearningcomputes

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ComputeInstanceDataMount struct {
	CreatedBy   *string      `json:"createdBy,omitempty"`
	Error       *string      `json:"error,omitempty"`
	MountAction *MountAction `json:"mountAction,omitempty"`
	MountName   *string      `json:"mountName,omitempty"`
	MountPath   *string      `json:"mountPath,omitempty"`
	MountState  *MountState  `json:"mountState,omitempty"`
	MountedOn   *string      `json:"mountedOn,omitempty"`
	Source      *string      `json:"source,omitempty"`
	SourceType  *SourceType  `json:"sourceType,omitempty"`
}

func (o *ComputeInstanceDataMount) GetMountedOnAsTime() (*time.Time, error) {
	if o.MountedOn == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.MountedOn, "2006-01-02T15:04:05Z07:00")
}

func (o *ComputeInstanceDataMount) SetMountedOnAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.MountedOn = &formatted
}
