package managedclustersnapshots

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagedClusterSnapshotProperties struct {
	CreationData                     *CreationData                        `json:"creationData,omitempty"`
	ManagedClusterPropertiesReadOnly *ManagedClusterPropertiesForSnapshot `json:"managedClusterPropertiesReadOnly,omitempty"`
	SnapshotType                     *SnapshotType                        `json:"snapshotType,omitempty"`
}
