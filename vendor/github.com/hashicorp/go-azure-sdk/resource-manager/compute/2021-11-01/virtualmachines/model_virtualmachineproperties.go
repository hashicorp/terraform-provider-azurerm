package virtualmachines

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualMachineProperties struct {
	AdditionalCapabilities  *AdditionalCapabilities            `json:"additionalCapabilities,omitempty"`
	ApplicationProfile      *ApplicationProfile                `json:"applicationProfile,omitempty"`
	AvailabilitySet         *SubResource                       `json:"availabilitySet,omitempty"`
	BillingProfile          *BillingProfile                    `json:"billingProfile,omitempty"`
	CapacityReservation     *CapacityReservationProfile        `json:"capacityReservation,omitempty"`
	DiagnosticsProfile      *DiagnosticsProfile                `json:"diagnosticsProfile,omitempty"`
	EvictionPolicy          *VirtualMachineEvictionPolicyTypes `json:"evictionPolicy,omitempty"`
	ExtensionsTimeBudget    *string                            `json:"extensionsTimeBudget,omitempty"`
	HardwareProfile         *HardwareProfile                   `json:"hardwareProfile,omitempty"`
	Host                    *SubResource                       `json:"host,omitempty"`
	HostGroup               *SubResource                       `json:"hostGroup,omitempty"`
	InstanceView            *VirtualMachineInstanceView        `json:"instanceView,omitempty"`
	LicenseType             *string                            `json:"licenseType,omitempty"`
	NetworkProfile          *NetworkProfile                    `json:"networkProfile,omitempty"`
	OsProfile               *OSProfile                         `json:"osProfile,omitempty"`
	PlatformFaultDomain     *int64                             `json:"platformFaultDomain,omitempty"`
	Priority                *VirtualMachinePriorityTypes       `json:"priority,omitempty"`
	ProvisioningState       *string                            `json:"provisioningState,omitempty"`
	ProximityPlacementGroup *SubResource                       `json:"proximityPlacementGroup,omitempty"`
	ScheduledEventsProfile  *ScheduledEventsProfile            `json:"scheduledEventsProfile,omitempty"`
	SecurityProfile         *SecurityProfile                   `json:"securityProfile,omitempty"`
	StorageProfile          *StorageProfile                    `json:"storageProfile,omitempty"`
	TimeCreated             *string                            `json:"timeCreated,omitempty"`
	UserData                *string                            `json:"userData,omitempty"`
	VMId                    *string                            `json:"vmId,omitempty"`
	VirtualMachineScaleSet  *SubResource                       `json:"virtualMachineScaleSet,omitempty"`
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
