package autonomousdatabasebackups

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AutonomousDatabaseBackupProperties struct {
	AutonomousDatabaseOcid *string                                 `json:"autonomousDatabaseOcid,omitempty"`
	BackupType             *AutonomousDatabaseBackupType           `json:"backupType,omitempty"`
	DatabaseSizeInTbs      *float64                                `json:"databaseSizeInTbs,omitempty"`
	DbVersion              *string                                 `json:"dbVersion,omitempty"`
	DisplayName            *string                                 `json:"displayName,omitempty"`
	IsAutomatic            *bool                                   `json:"isAutomatic,omitempty"`
	IsRestorable           *bool                                   `json:"isRestorable,omitempty"`
	LifecycleDetails       *string                                 `json:"lifecycleDetails,omitempty"`
	LifecycleState         *AutonomousDatabaseBackupLifecycleState `json:"lifecycleState,omitempty"`
	Ocid                   *string                                 `json:"ocid,omitempty"`
	ProvisioningState      *AzureResourceProvisioningState         `json:"provisioningState,omitempty"`
	RetentionPeriodInDays  *int64                                  `json:"retentionPeriodInDays,omitempty"`
	SizeInTbs              *float64                                `json:"sizeInTbs,omitempty"`
	TimeAvailableTil       *string                                 `json:"timeAvailableTil,omitempty"`
	TimeEnded              *string                                 `json:"timeEnded,omitempty"`
	TimeStarted            *string                                 `json:"timeStarted,omitempty"`
}

func (o *AutonomousDatabaseBackupProperties) GetTimeAvailableTilAsTime() (*time.Time, error) {
	if o.TimeAvailableTil == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.TimeAvailableTil, "2006-01-02T15:04:05Z07:00")
}

func (o *AutonomousDatabaseBackupProperties) SetTimeAvailableTilAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.TimeAvailableTil = &formatted
}
