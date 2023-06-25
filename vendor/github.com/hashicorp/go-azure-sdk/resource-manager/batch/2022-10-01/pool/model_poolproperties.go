package pool

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PoolProperties struct {
	AllocationState                 *AllocationState               `json:"allocationState,omitempty"`
	AllocationStateTransitionTime   *string                        `json:"allocationStateTransitionTime,omitempty"`
	ApplicationLicenses             *[]string                      `json:"applicationLicenses,omitempty"`
	ApplicationPackages             *[]ApplicationPackageReference `json:"applicationPackages,omitempty"`
	AutoScaleRun                    *AutoScaleRun                  `json:"autoScaleRun,omitempty"`
	Certificates                    *[]CertificateReference        `json:"certificates,omitempty"`
	CreationTime                    *string                        `json:"creationTime,omitempty"`
	CurrentDedicatedNodes           *int64                         `json:"currentDedicatedNodes,omitempty"`
	CurrentLowPriorityNodes         *int64                         `json:"currentLowPriorityNodes,omitempty"`
	CurrentNodeCommunicationMode    *NodeCommunicationMode         `json:"currentNodeCommunicationMode,omitempty"`
	DeploymentConfiguration         *DeploymentConfiguration       `json:"deploymentConfiguration,omitempty"`
	DisplayName                     *string                        `json:"displayName,omitempty"`
	InterNodeCommunication          *InterNodeCommunicationState   `json:"interNodeCommunication,omitempty"`
	LastModified                    *string                        `json:"lastModified,omitempty"`
	Metadata                        *[]MetadataItem                `json:"metadata,omitempty"`
	MountConfiguration              *[]MountConfiguration          `json:"mountConfiguration,omitempty"`
	NetworkConfiguration            *NetworkConfiguration          `json:"networkConfiguration,omitempty"`
	ProvisioningState               *PoolProvisioningState         `json:"provisioningState,omitempty"`
	ProvisioningStateTransitionTime *string                        `json:"provisioningStateTransitionTime,omitempty"`
	ResizeOperationStatus           *ResizeOperationStatus         `json:"resizeOperationStatus,omitempty"`
	ScaleSettings                   *ScaleSettings                 `json:"scaleSettings,omitempty"`
	StartTask                       *StartTask                     `json:"startTask,omitempty"`
	TargetNodeCommunicationMode     *NodeCommunicationMode         `json:"targetNodeCommunicationMode,omitempty"`
	TaskSchedulingPolicy            *TaskSchedulingPolicy          `json:"taskSchedulingPolicy,omitempty"`
	TaskSlotsPerNode                *int64                         `json:"taskSlotsPerNode,omitempty"`
	UserAccounts                    *[]UserAccount                 `json:"userAccounts,omitempty"`
	VMSize                          *string                        `json:"vmSize,omitempty"`
}

func (o *PoolProperties) GetAllocationStateTransitionTimeAsTime() (*time.Time, error) {
	if o.AllocationStateTransitionTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.AllocationStateTransitionTime, "2006-01-02T15:04:05Z07:00")
}

func (o *PoolProperties) SetAllocationStateTransitionTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.AllocationStateTransitionTime = &formatted
}

func (o *PoolProperties) GetCreationTimeAsTime() (*time.Time, error) {
	if o.CreationTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreationTime, "2006-01-02T15:04:05Z07:00")
}

func (o *PoolProperties) SetCreationTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreationTime = &formatted
}

func (o *PoolProperties) GetLastModifiedAsTime() (*time.Time, error) {
	if o.LastModified == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastModified, "2006-01-02T15:04:05Z07:00")
}

func (o *PoolProperties) SetLastModifiedAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastModified = &formatted
}

func (o *PoolProperties) GetProvisioningStateTransitionTimeAsTime() (*time.Time, error) {
	if o.ProvisioningStateTransitionTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.ProvisioningStateTransitionTime, "2006-01-02T15:04:05Z07:00")
}

func (o *PoolProperties) SetProvisioningStateTransitionTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.ProvisioningStateTransitionTime = &formatted
}
