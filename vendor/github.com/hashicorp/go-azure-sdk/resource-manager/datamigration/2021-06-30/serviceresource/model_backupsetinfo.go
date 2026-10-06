package serviceresource

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BackupSetInfo struct {
	BackupFinishedDate *string           `json:"backupFinishedDate,omitempty"`
	BackupSetId        *string           `json:"backupSetId,omitempty"`
	BackupStartDate    *string           `json:"backupStartDate,omitempty"`
	BackupType         *BackupType       `json:"backupType,omitempty"`
	DatabaseName       *string           `json:"databaseName,omitempty"`
	FirstLsn           *string           `json:"firstLsn,omitempty"`
	IsBackupRestored   *bool             `json:"isBackupRestored,omitempty"`
	LastLsn            *string           `json:"lastLsn,omitempty"`
	LastModifiedTime   *string           `json:"lastModifiedTime,omitempty"`
	ListOfBackupFiles  *[]BackupFileInfo `json:"listOfBackupFiles,omitempty"`
}

func (o *BackupSetInfo) GetBackupFinishedDateAsTime() (*time.Time, error) {
	if o.BackupFinishedDate == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.BackupFinishedDate, "2006-01-02T15:04:05Z07:00")
}

func (o *BackupSetInfo) SetBackupFinishedDateAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.BackupFinishedDate = &formatted
}

func (o *BackupSetInfo) GetBackupStartDateAsTime() (*time.Time, error) {
	if o.BackupStartDate == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.BackupStartDate, "2006-01-02T15:04:05Z07:00")
}

func (o *BackupSetInfo) SetBackupStartDateAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.BackupStartDate = &formatted
}

func (o *BackupSetInfo) GetLastModifiedTimeAsTime() (*time.Time, error) {
	if o.LastModifiedTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastModifiedTime, "2006-01-02T15:04:05Z07:00")
}

func (o *BackupSetInfo) SetLastModifiedTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastModifiedTime = &formatted
}
