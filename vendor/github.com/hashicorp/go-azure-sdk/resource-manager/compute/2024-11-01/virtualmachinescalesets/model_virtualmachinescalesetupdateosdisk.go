package virtualmachinescalesets

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualMachineScaleSetUpdateOSDisk struct {
	Caching                 *CachingTypes                                `json:"caching,omitempty"`
	DeleteOption            *DiskDeleteOptionTypes                       `json:"deleteOption,omitempty"`
	DiffDiskSettings        *DiffDiskSettings                            `json:"diffDiskSettings,omitempty"`
	DiskSizeGB              *int64                                       `json:"diskSizeGB,omitempty"`
	Image                   *VirtualHardDisk                             `json:"image,omitempty"`
	ManagedDisk             *VirtualMachineScaleSetManagedDiskParameters `json:"managedDisk,omitempty"`
	VhdContainers           *[]string                                    `json:"vhdContainers,omitempty"`
	WriteAcceleratorEnabled *bool                                        `json:"writeAcceleratorEnabled,omitempty"`
}
