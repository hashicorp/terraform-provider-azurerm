package volumes

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VolumePropertiesDataProtection struct {
	Backup      *VolumeBackupProperties   `json:"backup"`
	Replication *ReplicationObject        `json:"replication"`
	Snapshot    *VolumeSnapshotProperties `json:"snapshot"`
}
