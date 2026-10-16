package cosmosdb

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PeriodicModeProperties struct {
	BackupIntervalInMinutes        *int64                   `json:"backupIntervalInMinutes,omitempty"`
	BackupRetentionIntervalInHours *int64                   `json:"backupRetentionIntervalInHours,omitempty"`
	BackupStorageRedundancy        *BackupStorageRedundancy `json:"backupStorageRedundancy,omitempty"`
}
