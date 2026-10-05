package virtualharddisks

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualHardDiskProperties struct {
	BlockSizeBytes      *int64                 `json:"blockSizeBytes,omitempty"`
	ContainerId         *string                `json:"containerId,omitempty"`
	DiskFileFormat      *DiskFileFormat        `json:"diskFileFormat,omitempty"`
	DiskSizeGB          *int64                 `json:"diskSizeGB,omitempty"`
	Dynamic             *bool                  `json:"dynamic,omitempty"`
	HyperVGeneration    *HyperVGeneration      `json:"hyperVGeneration,omitempty"`
	LogicalSectorBytes  *int64                 `json:"logicalSectorBytes,omitempty"`
	PhysicalSectorBytes *int64                 `json:"physicalSectorBytes,omitempty"`
	ProvisioningState   *ProvisioningStateEnum `json:"provisioningState,omitempty"`
	Status              *VirtualHardDiskStatus `json:"status,omitempty"`
}
