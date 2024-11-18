package pool

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OSDisk struct {
	Caching                 *CachingType      `json:"caching,omitempty"`
	DiskSizeGB              *int64            `json:"diskSizeGB,omitempty"`
	EphemeralOSDiskSettings *DiffDiskSettings `json:"ephemeralOSDiskSettings,omitempty"`
	ManagedDisk             *ManagedDisk      `json:"managedDisk,omitempty"`
	WriteAcceleratorEnabled *bool             `json:"writeAcceleratorEnabled,omitempty"`
}
