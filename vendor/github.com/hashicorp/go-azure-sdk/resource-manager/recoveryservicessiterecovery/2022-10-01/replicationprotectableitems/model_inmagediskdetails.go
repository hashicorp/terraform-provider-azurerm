package replicationprotectableitems

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type InMageDiskDetails struct {
	DiskConfiguration *string              `json:"diskConfiguration,omitempty"`
	DiskId            *string              `json:"diskId,omitempty"`
	DiskName          *string              `json:"diskName,omitempty"`
	DiskSizeInMB      *string              `json:"diskSizeInMB,omitempty"`
	DiskType          *string              `json:"diskType,omitempty"`
	VolumeList        *[]DiskVolumeDetails `json:"volumeList,omitempty"`
}
