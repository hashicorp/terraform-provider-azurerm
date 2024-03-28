package virtualmachines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualMachineSize struct {
	MaxDataDiskCount     *int64  `json:"maxDataDiskCount,omitempty"`
	MemoryInMB           *int64  `json:"memoryInMB,omitempty"`
	Name                 *string `json:"name,omitempty"`
	NumberOfCores        *int64  `json:"numberOfCores,omitempty"`
	OsDiskSizeInMB       *int64  `json:"osDiskSizeInMB,omitempty"`
	ResourceDiskSizeInMB *int64  `json:"resourceDiskSizeInMB,omitempty"`
}
