package virtualmachinetemplates

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualDisk struct {
	Bus              *int64                   `json:"bus,omitempty"`
	BusType          *string                  `json:"busType,omitempty"`
	CreateDiffDisk   *CreateDiffDisk          `json:"createDiffDisk,omitempty"`
	DiskId           *string                  `json:"diskId,omitempty"`
	DiskSizeGB       *int64                   `json:"diskSizeGB,omitempty"`
	DisplayName      *string                  `json:"displayName,omitempty"`
	Lun              *int64                   `json:"lun,omitempty"`
	MaxDiskSizeGB    *int64                   `json:"maxDiskSizeGB,omitempty"`
	Name             *string                  `json:"name,omitempty"`
	StorageQoSPolicy *StorageQosPolicyDetails `json:"storageQoSPolicy,omitempty"`
	TemplateDiskId   *string                  `json:"templateDiskId,omitempty"`
	VhdFormatType    *string                  `json:"vhdFormatType,omitempty"`
	VhdType          *string                  `json:"vhdType,omitempty"`
	VolumeType       *string                  `json:"volumeType,omitempty"`
}
