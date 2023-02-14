package virtualmachines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ComputeDataDisk struct {
	DiskSizeGiB   *int64  `json:"diskSizeGiB,omitempty"`
	DiskUri       *string `json:"diskUri,omitempty"`
	ManagedDiskId *string `json:"managedDiskId,omitempty"`
	Name          *string `json:"name,omitempty"`
}
