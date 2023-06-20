package sqlvirtualmachines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AutoBackupSettings struct {
	BackupScheduleType    *BackupScheduleType      `json:"backupScheduleType,omitempty"`
	BackupSystemDbs       *bool                    `json:"backupSystemDbs,omitempty"`
	DaysOfWeek            *[]AutoBackupDaysOfWeek  `json:"daysOfWeek,omitempty"`
	Enable                *bool                    `json:"enable,omitempty"`
	EnableEncryption      *bool                    `json:"enableEncryption,omitempty"`
	FullBackupFrequency   *FullBackupFrequencyType `json:"fullBackupFrequency,omitempty"`
	FullBackupStartTime   *int64                   `json:"fullBackupStartTime,omitempty"`
	FullBackupWindowHours *int64                   `json:"fullBackupWindowHours,omitempty"`
	LogBackupFrequency    *int64                   `json:"logBackupFrequency,omitempty"`
	Password              *string                  `json:"password,omitempty"`
	RetentionPeriod       *int64                   `json:"retentionPeriod,omitempty"`
	StorageAccessKey      *string                  `json:"storageAccessKey,omitempty"`
	StorageAccountUrl     *string                  `json:"storageAccountUrl,omitempty"`
	StorageContainerName  *string                  `json:"storageContainerName,omitempty"`
}
