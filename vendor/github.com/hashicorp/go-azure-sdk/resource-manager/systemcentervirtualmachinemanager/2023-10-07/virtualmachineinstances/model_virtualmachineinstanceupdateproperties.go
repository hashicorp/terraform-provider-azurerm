package virtualmachineinstances

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualMachineInstanceUpdateProperties struct {
	AvailabilitySets      *[]AvailabilitySetListItem   `json:"availabilitySets,omitempty"`
	HardwareProfile       *HardwareProfileUpdate       `json:"hardwareProfile,omitempty"`
	InfrastructureProfile *InfrastructureProfileUpdate `json:"infrastructureProfile,omitempty"`
	NetworkProfile        *NetworkProfileUpdate        `json:"networkProfile,omitempty"`
	StorageProfile        *StorageProfileUpdate        `json:"storageProfile,omitempty"`
}
