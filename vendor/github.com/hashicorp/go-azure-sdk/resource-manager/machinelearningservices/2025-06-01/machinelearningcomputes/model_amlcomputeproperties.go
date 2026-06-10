package machinelearningcomputes

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AmlComputeProperties struct {
	AllocationState               *AllocationState             `json:"allocationState,omitempty"`
	AllocationStateTransitionTime *string                      `json:"allocationStateTransitionTime,omitempty"`
	CurrentNodeCount              *int64                       `json:"currentNodeCount,omitempty"`
	EnableNodePublicIP            *bool                        `json:"enableNodePublicIp,omitempty"`
	Errors                        *[]ErrorResponse             `json:"errors,omitempty"`
	IsolatedNetwork               *bool                        `json:"isolatedNetwork,omitempty"`
	NodeStateCounts               *NodeStateCounts             `json:"nodeStateCounts,omitempty"`
	OsType                        *OsType                      `json:"osType,omitempty"`
	PropertyBag                   *interface{}                 `json:"propertyBag,omitempty"`
	RemoteLoginPortPublicAccess   *RemoteLoginPortPublicAccess `json:"remoteLoginPortPublicAccess,omitempty"`
	ScaleSettings                 *ScaleSettings               `json:"scaleSettings,omitempty"`
	Subnet                        *ResourceId                  `json:"subnet,omitempty"`
	TargetNodeCount               *int64                       `json:"targetNodeCount,omitempty"`
	UserAccountCredentials        *UserAccountCredentials      `json:"userAccountCredentials,omitempty"`
	VMPriority                    *VMPriority                  `json:"vmPriority,omitempty"`
	VMSize                        *string                      `json:"vmSize,omitempty"`
	VirtualMachineImage           *VirtualMachineImage         `json:"virtualMachineImage,omitempty"`
}

func (o *AmlComputeProperties) GetAllocationStateTransitionTimeAsTime() (*time.Time, error) {
	if o.AllocationStateTransitionTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.AllocationStateTransitionTime, "2006-01-02T15:04:05Z07:00")
}

func (o *AmlComputeProperties) SetAllocationStateTransitionTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.AllocationStateTransitionTime = &formatted
}
