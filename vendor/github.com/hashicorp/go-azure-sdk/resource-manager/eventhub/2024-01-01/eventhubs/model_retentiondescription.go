package eventhubs

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RetentionDescription struct {
	CleanupPolicy                 *CleanupPolicyRetentionDescription `json:"cleanupPolicy,omitempty"`
	RetentionTimeInHours          *int64                             `json:"retentionTimeInHours,omitempty"`
	TombstoneRetentionTimeInHours *int64                             `json:"tombstoneRetentionTimeInHours,omitempty"`
}
