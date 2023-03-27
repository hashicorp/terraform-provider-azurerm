package sapdatabaseinstances

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DatabaseVMDetails struct {
	Status           *SAPVirtualInstanceStatus `json:"status,omitempty"`
	StorageDetails   *[]StorageInformation     `json:"storageDetails,omitempty"`
	VirtualMachineId *string                   `json:"virtualMachineId,omitempty"`
}
