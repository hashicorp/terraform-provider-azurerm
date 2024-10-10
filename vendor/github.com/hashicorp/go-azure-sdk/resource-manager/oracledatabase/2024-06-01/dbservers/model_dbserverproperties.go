package dbservers

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DbServerProperties struct {
	AutonomousVMClusterIds      *[]string                  `json:"autonomousVmClusterIds,omitempty"`
	AutonomousVirtualMachineIds *[]string                  `json:"autonomousVirtualMachineIds,omitempty"`
	CompartmentId               *string                    `json:"compartmentId,omitempty"`
	CpuCoreCount                *int64                     `json:"cpuCoreCount,omitempty"`
	DbNodeIds                   *[]string                  `json:"dbNodeIds,omitempty"`
	DbNodeStorageSizeInGbs      *int64                     `json:"dbNodeStorageSizeInGbs,omitempty"`
	DbServerPatchingDetails     *DbServerPatchingDetails   `json:"dbServerPatchingDetails,omitempty"`
	DisplayName                 *string                    `json:"displayName,omitempty"`
	ExadataInfrastructureId     *string                    `json:"exadataInfrastructureId,omitempty"`
	LifecycleDetails            *string                    `json:"lifecycleDetails,omitempty"`
	LifecycleState              *DbServerProvisioningState `json:"lifecycleState,omitempty"`
	MaxCPUCount                 *int64                     `json:"maxCpuCount,omitempty"`
	MaxDbNodeStorageInGbs       *int64                     `json:"maxDbNodeStorageInGbs,omitempty"`
	MaxMemoryInGbs              *int64                     `json:"maxMemoryInGbs,omitempty"`
	MemorySizeInGbs             *int64                     `json:"memorySizeInGbs,omitempty"`
	Ocid                        *string                    `json:"ocid,omitempty"`
	ProvisioningState           *ResourceProvisioningState `json:"provisioningState,omitempty"`
	Shape                       *string                    `json:"shape,omitempty"`
	TimeCreated                 *string                    `json:"timeCreated,omitempty"`
	VMClusterIds                *[]string                  `json:"vmClusterIds,omitempty"`
}

func (o *DbServerProperties) GetTimeCreatedAsTime() (*time.Time, error) {
	if o.TimeCreated == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.TimeCreated, "2006-01-02T15:04:05Z07:00")
}

func (o *DbServerProperties) SetTimeCreatedAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.TimeCreated = &formatted
}
