package longtermretentionbackups

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LongTermRetentionBackupProperties struct {
	BackupExpirationTime             *string                  `json:"backupExpirationTime,omitempty"`
	BackupStorageRedundancy          *BackupStorageRedundancy `json:"backupStorageRedundancy,omitempty"`
	BackupTime                       *string                  `json:"backupTime,omitempty"`
	DatabaseDeletionTime             *string                  `json:"databaseDeletionTime,omitempty"`
	DatabaseName                     *string                  `json:"databaseName,omitempty"`
	IsBackupImmutable                *bool                    `json:"isBackupImmutable,omitempty"`
	RequestedBackupStorageRedundancy *BackupStorageRedundancy `json:"requestedBackupStorageRedundancy,omitempty"`
	ServerCreateTime                 *string                  `json:"serverCreateTime,omitempty"`
	ServerName                       *string                  `json:"serverName,omitempty"`
}

func (o *LongTermRetentionBackupProperties) GetBackupExpirationTimeAsTime() (*time.Time, error) {
	if o.BackupExpirationTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.BackupExpirationTime, "2006-01-02T15:04:05Z07:00")
}

func (o *LongTermRetentionBackupProperties) SetBackupExpirationTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.BackupExpirationTime = &formatted
}

func (o *LongTermRetentionBackupProperties) GetBackupTimeAsTime() (*time.Time, error) {
	if o.BackupTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.BackupTime, "2006-01-02T15:04:05Z07:00")
}

func (o *LongTermRetentionBackupProperties) SetBackupTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.BackupTime = &formatted
}

func (o *LongTermRetentionBackupProperties) GetDatabaseDeletionTimeAsTime() (*time.Time, error) {
	if o.DatabaseDeletionTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.DatabaseDeletionTime, "2006-01-02T15:04:05Z07:00")
}

func (o *LongTermRetentionBackupProperties) SetDatabaseDeletionTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.DatabaseDeletionTime = &formatted
}

func (o *LongTermRetentionBackupProperties) GetServerCreateTimeAsTime() (*time.Time, error) {
	if o.ServerCreateTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.ServerCreateTime, "2006-01-02T15:04:05Z07:00")
}

func (o *LongTermRetentionBackupProperties) SetServerCreateTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.ServerCreateTime = &formatted
}
