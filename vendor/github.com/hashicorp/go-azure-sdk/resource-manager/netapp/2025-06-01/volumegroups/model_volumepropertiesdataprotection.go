package volumegroups

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VolumePropertiesDataProtection struct {
	Backup           *VolumeBackupProperties     `json:"backup,omitempty"`
	Replication      *ReplicationObject          `json:"replication,omitempty"`
	Snapshot         *VolumeSnapshotProperties   `json:"snapshot,omitempty"`
	VolumeRelocation *VolumeRelocationProperties `json:"volumeRelocation,omitempty"`
}
