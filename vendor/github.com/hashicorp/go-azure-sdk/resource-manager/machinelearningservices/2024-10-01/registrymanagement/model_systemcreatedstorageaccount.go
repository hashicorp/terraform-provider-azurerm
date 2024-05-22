package registrymanagement

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SystemCreatedStorageAccount struct {
	AllowBlobPublicAccess    *bool          `json:"allowBlobPublicAccess,omitempty"`
	ArmResourceId            *ArmResourceId `json:"armResourceId,omitempty"`
	StorageAccountHnsEnabled *bool          `json:"storageAccountHnsEnabled,omitempty"`
	StorageAccountName       *string        `json:"storageAccountName,omitempty"`
	StorageAccountType       *string        `json:"storageAccountType,omitempty"`
}
