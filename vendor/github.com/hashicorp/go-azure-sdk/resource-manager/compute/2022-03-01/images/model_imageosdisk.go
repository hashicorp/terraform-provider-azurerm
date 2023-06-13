package images

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ImageOSDisk struct {
	BlobUri            *string                   `json:"blobUri,omitempty"`
	Caching            *CachingTypes             `json:"caching,omitempty"`
	DiskEncryptionSet  *SubResource              `json:"diskEncryptionSet,omitempty"`
	DiskSizeGB         *int64                    `json:"diskSizeGB,omitempty"`
	ManagedDisk        *SubResource              `json:"managedDisk,omitempty"`
	OsState            OperatingSystemStateTypes `json:"osState"`
	OsType             OperatingSystemTypes      `json:"osType"`
	Snapshot           *SubResource              `json:"snapshot,omitempty"`
	StorageAccountType *StorageAccountTypes      `json:"storageAccountType,omitempty"`
}
