package machines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VMwareDisk struct {
	DiskMode               *VirtualDiskMode `json:"diskMode,omitempty"`
	DiskProvisioningPolicy *string          `json:"diskProvisioningPolicy,omitempty"`
	DiskScrubbingPolicy    *string          `json:"diskScrubbingPolicy,omitempty"`
	DiskType               *string          `json:"diskType,omitempty"`
	Label                  *string          `json:"label,omitempty"`
	Lun                    *int64           `json:"lun,omitempty"`
	MaxSizeInBytes         *int64           `json:"maxSizeInBytes,omitempty"`
	Name                   *string          `json:"name,omitempty"`
	Path                   *string          `json:"path,omitempty"`
	Uuid                   *string          `json:"uuid,omitempty"`
}
