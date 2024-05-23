package regions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VMSizeProperty struct {
	Cores                              *int64  `json:"cores,omitempty"`
	DataDiskStorageTier                *string `json:"dataDiskStorageTier,omitempty"`
	Label                              *string `json:"label,omitempty"`
	MaxDataDiskCount                   *int64  `json:"maxDataDiskCount,omitempty"`
	MemoryInMb                         *int64  `json:"memoryInMb,omitempty"`
	Name                               *string `json:"name,omitempty"`
	SupportedByVirtualMachines         *bool   `json:"supportedByVirtualMachines,omitempty"`
	SupportedByWebWorkerRoles          *bool   `json:"supportedByWebWorkerRoles,omitempty"`
	VirtualMachineResourceDiskSizeInMb *int64  `json:"virtualMachineResourceDiskSizeInMb,omitempty"`
	WebWorkerResourceDiskSizeInMb      *int64  `json:"webWorkerResourceDiskSizeInMb,omitempty"`
}
