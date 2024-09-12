package virtualmachineinstances

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualMachineInstanceProperties struct {
	AvailabilitySets      *[]AvailabilitySetListItem `json:"availabilitySets,omitempty"`
	HardwareProfile       *HardwareProfile           `json:"hardwareProfile,omitempty"`
	InfrastructureProfile *InfrastructureProfile     `json:"infrastructureProfile,omitempty"`
	NetworkProfile        *NetworkProfile            `json:"networkProfile,omitempty"`
	OsProfile             *OsProfileForVMInstance    `json:"osProfile,omitempty"`
	PowerState            *string                    `json:"powerState,omitempty"`
	ProvisioningState     *ProvisioningState         `json:"provisioningState,omitempty"`
	StorageProfile        *StorageProfile            `json:"storageProfile,omitempty"`
}
