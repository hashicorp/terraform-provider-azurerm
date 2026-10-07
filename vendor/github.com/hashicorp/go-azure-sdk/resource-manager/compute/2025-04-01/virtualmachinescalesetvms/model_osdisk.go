package virtualmachinescalesetvms

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OSDisk struct {
	Caching                 *CachingTypes           `json:"caching,omitempty"`
	CreateOption            DiskCreateOptionTypes   `json:"createOption"`
	DeleteOption            *DiskDeleteOptionTypes  `json:"deleteOption,omitempty"`
	DiffDiskSettings        *DiffDiskSettings       `json:"diffDiskSettings,omitempty"`
	DiskSizeGB              *int64                  `json:"diskSizeGB,omitempty"`
	EncryptionSettings      *DiskEncryptionSettings `json:"encryptionSettings,omitempty"`
	Image                   *VirtualHardDisk        `json:"image,omitempty"`
	ManagedDisk             *ManagedDiskParameters  `json:"managedDisk,omitempty"`
	Name                    *string                 `json:"name,omitempty"`
	OsType                  *OperatingSystemTypes   `json:"osType,omitempty"`
	Vhd                     *VirtualHardDisk        `json:"vhd,omitempty"`
	WriteAcceleratorEnabled *bool                   `json:"writeAcceleratorEnabled,omitempty"`
}
