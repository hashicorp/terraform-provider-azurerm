package virtualmachineinstances

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualDiskUpdate struct {
	Bus              *int64                   `json:"bus,omitempty"`
	BusType          *string                  `json:"busType,omitempty"`
	DiskId           *string                  `json:"diskId,omitempty"`
	DiskSizeGB       *int64                   `json:"diskSizeGB,omitempty"`
	Lun              *int64                   `json:"lun,omitempty"`
	Name             *string                  `json:"name,omitempty"`
	StorageQoSPolicy *StorageQosPolicyDetails `json:"storageQoSPolicy,omitempty"`
	VhdType          *string                  `json:"vhdType,omitempty"`
}
