package machinelearningcomputes

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ComputeInstanceDataDisk struct {
	Caching            *Caching            `json:"caching,omitempty"`
	DiskSizeGB         *int64              `json:"diskSizeGB,omitempty"`
	Lun                *int64              `json:"lun,omitempty"`
	StorageAccountType *StorageAccountType `json:"storageAccountType,omitempty"`
}
