package availabilitysets

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AvailabilitySetProperties struct {
	PlatformFaultDomainCount  *int64                `json:"platformFaultDomainCount,omitempty"`
	PlatformUpdateDomainCount *int64                `json:"platformUpdateDomainCount,omitempty"`
	ProximityPlacementGroup   *SubResource          `json:"proximityPlacementGroup,omitempty"`
	Statuses                  *[]InstanceViewStatus `json:"statuses,omitempty"`
	VirtualMachines           *[]SubResource        `json:"virtualMachines,omitempty"`
}
