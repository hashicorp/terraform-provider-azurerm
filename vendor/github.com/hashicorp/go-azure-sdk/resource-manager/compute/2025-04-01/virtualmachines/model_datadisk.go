package virtualmachines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DataDisk struct {
	Caching                 *CachingTypes          `json:"caching,omitempty"`
	CreateOption            DiskCreateOptionTypes  `json:"createOption"`
	DeleteOption            *DiskDeleteOptionTypes `json:"deleteOption,omitempty"`
	DetachOption            *DiskDetachOptionTypes `json:"detachOption,omitempty"`
	DiskIOPSReadWrite       *int64                 `json:"diskIOPSReadWrite,omitempty"`
	DiskMBpsReadWrite       *int64                 `json:"diskMBpsReadWrite,omitempty"`
	DiskSizeGB              *int64                 `json:"diskSizeGB,omitempty"`
	Image                   *VirtualHardDisk       `json:"image,omitempty"`
	Lun                     int64                  `json:"lun"`
	ManagedDisk             *ManagedDiskParameters `json:"managedDisk,omitempty"`
	Name                    *string                `json:"name,omitempty"`
	SourceResource          *ApiEntityReference    `json:"sourceResource,omitempty"`
	ToBeDetached            *bool                  `json:"toBeDetached,omitempty"`
	Vhd                     *VirtualHardDisk       `json:"vhd,omitempty"`
	WriteAcceleratorEnabled *bool                  `json:"writeAcceleratorEnabled,omitempty"`
}
