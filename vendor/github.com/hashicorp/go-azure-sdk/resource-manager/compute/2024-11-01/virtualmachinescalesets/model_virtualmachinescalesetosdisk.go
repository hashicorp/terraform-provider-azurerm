package virtualmachinescalesets

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualMachineScaleSetOSDisk struct {
	Caching                 *CachingTypes                                `json:"caching,omitempty"`
	CreateOption            DiskCreateOptionTypes                        `json:"createOption"`
	DeleteOption            *DiskDeleteOptionTypes                       `json:"deleteOption,omitempty"`
	DiffDiskSettings        *DiffDiskSettings                            `json:"diffDiskSettings,omitempty"`
	DiskSizeGB              *int64                                       `json:"diskSizeGB,omitempty"`
	Image                   *VirtualHardDisk                             `json:"image,omitempty"`
	ManagedDisk             *VirtualMachineScaleSetManagedDiskParameters `json:"managedDisk,omitempty"`
	Name                    *string                                      `json:"name,omitempty"`
	OsType                  *OperatingSystemTypes                        `json:"osType,omitempty"`
	VhdContainers           *[]string                                    `json:"vhdContainers,omitempty"`
	WriteAcceleratorEnabled *bool                                        `json:"writeAcceleratorEnabled,omitempty"`
}
