package virtualmachineinstances

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualMachineInstanceUpdateProperties struct {
	HardwareProfile *HardwareProfileUpdate `json:"hardwareProfile,omitempty"`
	NetworkProfile  *NetworkProfileUpdate  `json:"networkProfile,omitempty"`
	OsProfile       *OsProfileUpdate       `json:"osProfile,omitempty"`
	StorageProfile  *StorageProfileUpdate  `json:"storageProfile,omitempty"`
}
