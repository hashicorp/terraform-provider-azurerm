package virtualmachinescalesets

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualMachineScaleSetVMProfile struct {
	ApplicationProfile       *ApplicationProfile                     `json:"applicationProfile,omitempty"`
	BillingProfile           *BillingProfile                         `json:"billingProfile,omitempty"`
	CapacityReservation      *CapacityReservationProfile             `json:"capacityReservation,omitempty"`
	DiagnosticsProfile       *DiagnosticsProfile                     `json:"diagnosticsProfile,omitempty"`
	EvictionPolicy           *VirtualMachineEvictionPolicyTypes      `json:"evictionPolicy,omitempty"`
	ExtensionProfile         *VirtualMachineScaleSetExtensionProfile `json:"extensionProfile,omitempty"`
	HardwareProfile          *VirtualMachineScaleSetHardwareProfile  `json:"hardwareProfile,omitempty"`
	LicenseType              *string                                 `json:"licenseType,omitempty"`
	NetworkProfile           *VirtualMachineScaleSetNetworkProfile   `json:"networkProfile,omitempty"`
	OsProfile                *VirtualMachineScaleSetOSProfile        `json:"osProfile,omitempty"`
	Priority                 *VirtualMachinePriorityTypes            `json:"priority,omitempty"`
	ScheduledEventsProfile   *ScheduledEventsProfile                 `json:"scheduledEventsProfile,omitempty"`
	SecurityPostureReference *SecurityPostureReference               `json:"securityPostureReference,omitempty"`
	SecurityProfile          *SecurityProfile                        `json:"securityProfile,omitempty"`
	ServiceArtifactReference *ServiceArtifactReference               `json:"serviceArtifactReference,omitempty"`
	StorageProfile           *VirtualMachineScaleSetStorageProfile   `json:"storageProfile,omitempty"`
	TimeCreated              *string                                 `json:"timeCreated,omitempty"`
	UserData                 *string                                 `json:"userData,omitempty"`
}

func (o *VirtualMachineScaleSetVMProfile) GetTimeCreatedAsTime() (*time.Time, error) {
	if o.TimeCreated == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.TimeCreated, "2006-01-02T15:04:05Z07:00")
}

func (o *VirtualMachineScaleSetVMProfile) SetTimeCreatedAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.TimeCreated = &formatted
}
