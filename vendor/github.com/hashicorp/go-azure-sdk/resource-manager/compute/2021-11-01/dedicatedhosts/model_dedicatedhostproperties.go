package dedicatedhosts

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DedicatedHostProperties struct {
	AutoReplaceOnFailure *bool                      `json:"autoReplaceOnFailure,omitempty"`
	HostId               *string                    `json:"hostId,omitempty"`
	InstanceView         *DedicatedHostInstanceView `json:"instanceView,omitempty"`
	LicenseType          *DedicatedHostLicenseTypes `json:"licenseType,omitempty"`
	PlatformFaultDomain  *int64                     `json:"platformFaultDomain,omitempty"`
	ProvisioningState    *string                    `json:"provisioningState,omitempty"`
	ProvisioningTime     *string                    `json:"provisioningTime,omitempty"`
	TimeCreated          *string                    `json:"timeCreated,omitempty"`
	VirtualMachines      *[]SubResourceReadOnly     `json:"virtualMachines,omitempty"`
}

func (o *DedicatedHostProperties) GetProvisioningTimeAsTime() (*time.Time, error) {
	if o.ProvisioningTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.ProvisioningTime, "2006-01-02T15:04:05Z07:00")
}

func (o *DedicatedHostProperties) SetProvisioningTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.ProvisioningTime = &formatted
}

func (o *DedicatedHostProperties) GetTimeCreatedAsTime() (*time.Time, error) {
	if o.TimeCreated == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.TimeCreated, "2006-01-02T15:04:05Z07:00")
}

func (o *DedicatedHostProperties) SetTimeCreatedAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.TimeCreated = &formatted
}
