package locationbasedcapabilities

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StorageEditionCapability struct {
	MaxBackupIntervalHours *int64  `json:"maxBackupIntervalHours,omitempty"`
	MaxBackupRetentionDays *int64  `json:"maxBackupRetentionDays,omitempty"`
	MaxStorageSize         *int64  `json:"maxStorageSize,omitempty"`
	MinBackupIntervalHours *int64  `json:"minBackupIntervalHours,omitempty"`
	MinBackupRetentionDays *int64  `json:"minBackupRetentionDays,omitempty"`
	MinStorageSize         *int64  `json:"minStorageSize,omitempty"`
	Name                   *string `json:"name,omitempty"`
}
