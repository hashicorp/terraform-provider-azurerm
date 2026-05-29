package virtualmachinescalesets

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualMachineScaleSetDataDisk struct {
	Caching                 *CachingTypes                                `json:"caching,omitempty"`
	CreateOption            DiskCreateOptionTypes                        `json:"createOption"`
	DeleteOption            *DiskDeleteOptionTypes                       `json:"deleteOption,omitempty"`
	DiskIOPSReadWrite       *int64                                       `json:"diskIOPSReadWrite,omitempty"`
	DiskMBpsReadWrite       *int64                                       `json:"diskMBpsReadWrite,omitempty"`
	DiskSizeGB              *int64                                       `json:"diskSizeGB,omitempty"`
	Lun                     int64                                        `json:"lun"`
	ManagedDisk             *VirtualMachineScaleSetManagedDiskParameters `json:"managedDisk,omitempty"`
	Name                    *string                                      `json:"name,omitempty"`
	WriteAcceleratorEnabled *bool                                        `json:"writeAcceleratorEnabled,omitempty"`
}
