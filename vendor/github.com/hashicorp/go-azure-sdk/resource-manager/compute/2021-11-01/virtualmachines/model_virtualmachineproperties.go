package virtualmachines

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualMachineProperties struct {
	AdditionalCapabilities  *AdditionalCapabilities            `json:"additionalCapabilities"`
	ApplicationProfile      *ApplicationProfile                `json:"applicationProfile"`
	AvailabilitySet         *SubResource                       `json:"availabilitySet"`
	BillingProfile          *BillingProfile                    `json:"billingProfile"`
	CapacityReservation     *CapacityReservationProfile        `json:"capacityReservation"`
	DiagnosticsProfile      *DiagnosticsProfile                `json:"diagnosticsProfile"`
	EvictionPolicy          *VirtualMachineEvictionPolicyTypes `json:"evictionPolicy,omitempty"`
	ExtensionsTimeBudget    *string                            `json:"extensionsTimeBudget,omitempty"`
	HardwareProfile         *HardwareProfile                   `json:"hardwareProfile"`
	Host                    *SubResource                       `json:"host"`
	HostGroup               *SubResource                       `json:"hostGroup"`
	InstanceView            *VirtualMachineInstanceView        `json:"instanceView"`
	LicenseType             *string                            `json:"licenseType,omitempty"`
	NetworkProfile          *NetworkProfile                    `json:"networkProfile"`
	OsProfile               *OSProfile                         `json:"osProfile"`
	PlatformFaultDomain     *int64                             `json:"platformFaultDomain,omitempty"`
	Priority                *VirtualMachinePriorityTypes       `json:"priority,omitempty"`
	ProvisioningState       *string                            `json:"provisioningState,omitempty"`
	ProximityPlacementGroup *SubResource                       `json:"proximityPlacementGroup"`
	ScheduledEventsProfile  *ScheduledEventsProfile            `json:"scheduledEventsProfile"`
	SecurityProfile         *SecurityProfile                   `json:"securityProfile"`
	StorageProfile          *StorageProfile                    `json:"storageProfile"`
	TimeCreated             *string                            `json:"timeCreated,omitempty"`
	UserData                *string                            `json:"userData,omitempty"`
	VirtualMachineScaleSet  *SubResource                       `json:"virtualMachineScaleSet"`
	VmId                    *string                            `json:"vmId,omitempty"`
}

func (o *VirtualMachineProperties) GetTimeCreatedAsTime() (*time.Time, error) {
	if o.TimeCreated == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.TimeCreated, "2006-01-02T15:04:05Z07:00")
}

func (o *VirtualMachineProperties) SetTimeCreatedAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.TimeCreated = &formatted
}
