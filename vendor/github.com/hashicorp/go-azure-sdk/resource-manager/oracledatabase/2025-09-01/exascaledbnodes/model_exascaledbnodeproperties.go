package exascaledbnodes

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExascaleDbNodeProperties struct {
	AdditionalDetails          *string                  `json:"additionalDetails,omitempty"`
	CpuCoreCount               *int64                   `json:"cpuCoreCount,omitempty"`
	DbNodeStorageSizeInGbs     *int64                   `json:"dbNodeStorageSizeInGbs,omitempty"`
	FaultDomain                *string                  `json:"faultDomain,omitempty"`
	Hostname                   *string                  `json:"hostname,omitempty"`
	LifecycleState             *DbNodeProvisioningState `json:"lifecycleState,omitempty"`
	MaintenanceType            *string                  `json:"maintenanceType,omitempty"`
	MemorySizeInGbs            *int64                   `json:"memorySizeInGbs,omitempty"`
	Ocid                       string                   `json:"ocid"`
	SoftwareStorageSizeInGb    *int64                   `json:"softwareStorageSizeInGb,omitempty"`
	TimeMaintenanceWindowEnd   *string                  `json:"timeMaintenanceWindowEnd,omitempty"`
	TimeMaintenanceWindowStart *string                  `json:"timeMaintenanceWindowStart,omitempty"`
	TotalCPUCoreCount          *int64                   `json:"totalCpuCoreCount,omitempty"`
}

func (o *ExascaleDbNodeProperties) GetTimeMaintenanceWindowEndAsTime() (*time.Time, error) {
	if o.TimeMaintenanceWindowEnd == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.TimeMaintenanceWindowEnd, "2006-01-02T15:04:05Z07:00")
}

func (o *ExascaleDbNodeProperties) SetTimeMaintenanceWindowEndAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.TimeMaintenanceWindowEnd = &formatted
}

func (o *ExascaleDbNodeProperties) GetTimeMaintenanceWindowStartAsTime() (*time.Time, error) {
	if o.TimeMaintenanceWindowStart == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.TimeMaintenanceWindowStart, "2006-01-02T15:04:05Z07:00")
}

func (o *ExascaleDbNodeProperties) SetTimeMaintenanceWindowStartAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.TimeMaintenanceWindowStart = &formatted
}
