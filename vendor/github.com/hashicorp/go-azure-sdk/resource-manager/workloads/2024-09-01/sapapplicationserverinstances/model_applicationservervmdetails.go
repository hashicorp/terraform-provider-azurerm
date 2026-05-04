package sapapplicationserverinstances

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApplicationServerVMDetails struct {
	StorageDetails   *[]StorageInformation                `json:"storageDetails,omitempty"`
	Type             *ApplicationServerVirtualMachineType `json:"type,omitempty"`
	VirtualMachineId *string                              `json:"virtualMachineId,omitempty"`
}
