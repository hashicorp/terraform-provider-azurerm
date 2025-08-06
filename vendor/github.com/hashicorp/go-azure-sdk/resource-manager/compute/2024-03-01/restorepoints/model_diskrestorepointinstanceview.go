package restorepoints

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DiskRestorePointInstanceView struct {
	Id                *string                            `json:"id,omitempty"`
	ReplicationStatus *DiskRestorePointReplicationStatus `json:"replicationStatus,omitempty"`
}
