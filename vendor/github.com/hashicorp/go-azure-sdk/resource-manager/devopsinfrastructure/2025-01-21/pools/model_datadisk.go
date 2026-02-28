package pools

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DataDisk struct {
	Caching            *CachingType        `json:"caching,omitempty"`
	DiskSizeGiB        *int64              `json:"diskSizeGiB,omitempty"`
	DriveLetter        *string             `json:"driveLetter,omitempty"`
	StorageAccountType *StorageAccountType `json:"storageAccountType,omitempty"`
}
