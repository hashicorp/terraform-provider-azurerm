package dbsystemshapes

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DbSystemShapeProperties struct {
	AvailableCoreCount                 int64    `json:"availableCoreCount"`
	AvailableCoreCountPerNode          *int64   `json:"availableCoreCountPerNode,omitempty"`
	AvailableDataStorageInTbs          *int64   `json:"availableDataStorageInTbs,omitempty"`
	AvailableDataStoragePerServerInTbs *float64 `json:"availableDataStoragePerServerInTbs,omitempty"`
	AvailableDbNodePerNodeInGbs        *int64   `json:"availableDbNodePerNodeInGbs,omitempty"`
	AvailableDbNodeStorageInGbs        *int64   `json:"availableDbNodeStorageInGbs,omitempty"`
	AvailableMemoryInGbs               *int64   `json:"availableMemoryInGbs,omitempty"`
	AvailableMemoryPerNodeInGbs        *int64   `json:"availableMemoryPerNodeInGbs,omitempty"`
	CoreCountIncrement                 *int64   `json:"coreCountIncrement,omitempty"`
	MaxStorageCount                    *int64   `json:"maxStorageCount,omitempty"`
	MaximumNodeCount                   *int64   `json:"maximumNodeCount,omitempty"`
	MinCoreCountPerNode                *int64   `json:"minCoreCountPerNode,omitempty"`
	MinDataStorageInTbs                *int64   `json:"minDataStorageInTbs,omitempty"`
	MinDbNodeStoragePerNodeInGbs       *int64   `json:"minDbNodeStoragePerNodeInGbs,omitempty"`
	MinMemoryPerNodeInGbs              *int64   `json:"minMemoryPerNodeInGbs,omitempty"`
	MinStorageCount                    *int64   `json:"minStorageCount,omitempty"`
	MinimumCoreCount                   *int64   `json:"minimumCoreCount,omitempty"`
	MinimumNodeCount                   *int64   `json:"minimumNodeCount,omitempty"`
	RuntimeMinimumCoreCount            *int64   `json:"runtimeMinimumCoreCount,omitempty"`
	ShapeFamily                        *string  `json:"shapeFamily,omitempty"`
}
