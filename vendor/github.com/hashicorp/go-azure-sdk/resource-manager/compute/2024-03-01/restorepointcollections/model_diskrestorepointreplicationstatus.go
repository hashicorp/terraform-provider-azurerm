package restorepointcollections

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DiskRestorePointReplicationStatus struct {
	CompletionPercent *int64              `json:"completionPercent,omitempty"`
	Status            *InstanceViewStatus `json:"status,omitempty"`
}
