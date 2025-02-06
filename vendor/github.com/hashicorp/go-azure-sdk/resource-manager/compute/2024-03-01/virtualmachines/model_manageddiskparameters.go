package virtualmachines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagedDiskParameters struct {
	DiskEncryptionSet  *SubResource           `json:"diskEncryptionSet,omitempty"`
	Id                 *string                `json:"id,omitempty"`
	SecurityProfile    *VMDiskSecurityProfile `json:"securityProfile,omitempty"`
	StorageAccountType *StorageAccountTypes   `json:"storageAccountType,omitempty"`
}
