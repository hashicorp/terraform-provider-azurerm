package containerinstance

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureFileVolume struct {
	ReadOnly           *bool   `json:"readOnly,omitempty"`
	ShareName          string  `json:"shareName"`
	StorageAccountKey  *string `json:"storageAccountKey,omitempty"`
	StorageAccountName string  `json:"storageAccountName"`
}
