package sapcentralinstances

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CentralServerVMDetails struct {
	StorageDetails   *[]StorageInformation            `json:"storageDetails,omitempty"`
	Type             *CentralServerVirtualMachineType `json:"type,omitempty"`
	VirtualMachineId *string                          `json:"virtualMachineId,omitempty"`
}
