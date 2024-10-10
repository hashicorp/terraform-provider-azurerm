package serviceresource

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DatabaseBackupInfo struct {
	BackupFiles      *[]string   `json:"backupFiles,omitempty"`
	BackupFinishDate *string     `json:"backupFinishDate,omitempty"`
	BackupType       *BackupType `json:"backupType,omitempty"`
	DatabaseName     *string     `json:"databaseName,omitempty"`
	FamilyCount      *int64      `json:"familyCount,omitempty"`
	IsCompressed     *bool       `json:"isCompressed,omitempty"`
	IsDamaged        *bool       `json:"isDamaged,omitempty"`
	Position         *int64      `json:"position,omitempty"`
}

func (o *DatabaseBackupInfo) GetBackupFinishDateAsTime() (*time.Time, error) {
	if o.BackupFinishDate == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.BackupFinishDate, "2006-01-02T15:04:05Z07:00")
}

func (o *DatabaseBackupInfo) SetBackupFinishDateAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.BackupFinishDate = &formatted
}
