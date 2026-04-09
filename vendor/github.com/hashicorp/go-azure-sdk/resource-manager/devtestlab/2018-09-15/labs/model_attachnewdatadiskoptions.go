package labs

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AttachNewDataDiskOptions struct {
	DiskName    *string      `json:"diskName,omitempty"`
	DiskSizeGiB *int64       `json:"diskSizeGiB,omitempty"`
	DiskType    *StorageType `json:"diskType,omitempty"`
}
