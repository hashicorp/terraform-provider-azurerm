package amlfilesystems

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AmlFilesystemSubnetInfo struct {
	FilesystemSubnet   *string  `json:"filesystemSubnet,omitempty"`
	Location           *string  `json:"location,omitempty"`
	Sku                *SkuName `json:"sku,omitempty"`
	StorageCapacityTiB *float64 `json:"storageCapacityTiB,omitempty"`
}
