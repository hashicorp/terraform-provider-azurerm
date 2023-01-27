package managedcassandras

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CassandraClusterPublicStatusDataCentersInlinedNodesInlined struct {
	Address                  *string    `json:"address,omitempty"`
	CpuUsage                 *float64   `json:"cpuUsage,omitempty"`
	DiskFreeKB               *int64     `json:"diskFreeKB,omitempty"`
	DiskUsedKB               *int64     `json:"diskUsedKB,omitempty"`
	HostID                   *string    `json:"hostID,omitempty"`
	Load                     *string    `json:"load,omitempty"`
	MemoryBuffersAndCachedKB *int64     `json:"memoryBuffersAndCachedKB,omitempty"`
	MemoryFreeKB             *int64     `json:"memoryFreeKB,omitempty"`
	MemoryTotalKB            *int64     `json:"memoryTotalKB,omitempty"`
	MemoryUsedKB             *int64     `json:"memoryUsedKB,omitempty"`
	Rack                     *string    `json:"rack,omitempty"`
	Size                     *int64     `json:"size,omitempty"`
	State                    *NodeState `json:"state,omitempty"`
	Status                   *string    `json:"status,omitempty"`
	Timestamp                *string    `json:"timestamp,omitempty"`
	Tokens                   *[]string  `json:"tokens,omitempty"`
}
