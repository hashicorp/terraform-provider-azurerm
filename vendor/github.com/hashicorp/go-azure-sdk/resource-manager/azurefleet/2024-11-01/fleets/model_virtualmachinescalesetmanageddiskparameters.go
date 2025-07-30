package fleets

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualMachineScaleSetManagedDiskParameters struct {
	DiskEncryptionSet  *DiskEncryptionSetParameters `json:"diskEncryptionSet,omitempty"`
	SecurityProfile    *VMDiskSecurityProfile       `json:"securityProfile,omitempty"`
	StorageAccountType *StorageAccountTypes         `json:"storageAccountType,omitempty"`
}
