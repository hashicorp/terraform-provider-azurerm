package regions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DataDisksGroups struct {
	DiskSizeGB         *int64  `json:"diskSizeGB,omitempty"`
	DisksPerNode       *int64  `json:"disksPerNode,omitempty"`
	StorageAccountType *string `json:"storageAccountType,omitempty"`
}
