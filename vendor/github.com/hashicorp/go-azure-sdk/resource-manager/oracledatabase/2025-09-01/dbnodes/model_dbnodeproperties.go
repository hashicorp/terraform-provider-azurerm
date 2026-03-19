package dbnodes

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DbNodeProperties struct {
	AdditionalDetails          *string                    `json:"additionalDetails,omitempty"`
	BackupIPId                 *string                    `json:"backupIpId,omitempty"`
	BackupVnic2Id              *string                    `json:"backupVnic2Id,omitempty"`
	BackupVnicId               *string                    `json:"backupVnicId,omitempty"`
	CpuCoreCount               *int64                     `json:"cpuCoreCount,omitempty"`
	DbNodeStorageSizeInGbs     *int64                     `json:"dbNodeStorageSizeInGbs,omitempty"`
	DbServerId                 *string                    `json:"dbServerId,omitempty"`
	DbSystemId                 string                     `json:"dbSystemId"`
	FaultDomain                *string                    `json:"faultDomain,omitempty"`
	HostIPId                   *string                    `json:"hostIpId,omitempty"`
	Hostname                   *string                    `json:"hostname,omitempty"`
	LifecycleDetails           *string                    `json:"lifecycleDetails,omitempty"`
	LifecycleState             DbNodeProvisioningState    `json:"lifecycleState"`
	MaintenanceType            *DbNodeMaintenanceType     `json:"maintenanceType,omitempty"`
	MemorySizeInGbs            *int64                     `json:"memorySizeInGbs,omitempty"`
	Ocid                       string                     `json:"ocid"`
	ProvisioningState          *ResourceProvisioningState `json:"provisioningState,omitempty"`
	SoftwareStorageSizeInGb    *int64                     `json:"softwareStorageSizeInGb,omitempty"`
	TimeCreated                string                     `json:"timeCreated"`
	TimeMaintenanceWindowEnd   *string                    `json:"timeMaintenanceWindowEnd,omitempty"`
	TimeMaintenanceWindowStart *string                    `json:"timeMaintenanceWindowStart,omitempty"`
	Vnic2Id                    *string                    `json:"vnic2Id,omitempty"`
	VnicId                     string                     `json:"vnicId"`
}

func (o *DbNodeProperties) GetTimeCreatedAsTime() (*time.Time, error) {
	return dates.ParseAsFormat(&o.TimeCreated, "2006-01-02T15:04:05Z07:00")
}

func (o *DbNodeProperties) SetTimeCreatedAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.TimeCreated = formatted
}

func (o *DbNodeProperties) GetTimeMaintenanceWindowEndAsTime() (*time.Time, error) {
	if o.TimeMaintenanceWindowEnd == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.TimeMaintenanceWindowEnd, "2006-01-02T15:04:05Z07:00")
}

func (o *DbNodeProperties) SetTimeMaintenanceWindowEndAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.TimeMaintenanceWindowEnd = &formatted
}

func (o *DbNodeProperties) GetTimeMaintenanceWindowStartAsTime() (*time.Time, error) {
	if o.TimeMaintenanceWindowStart == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.TimeMaintenanceWindowStart, "2006-01-02T15:04:05Z07:00")
}

func (o *DbNodeProperties) SetTimeMaintenanceWindowStartAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.TimeMaintenanceWindowStart = &formatted
}
