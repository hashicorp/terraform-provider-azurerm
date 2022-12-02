package replicationfabrics

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AgentDiskDetails struct {
	CapacityInBytes *int64  `json:"capacityInBytes,omitempty"`
	DiskId          *string `json:"diskId,omitempty"`
	DiskName        *string `json:"diskName,omitempty"`
	IsOSDisk        *string `json:"isOSDisk,omitempty"`
	LunId           *int64  `json:"lunId,omitempty"`
}
