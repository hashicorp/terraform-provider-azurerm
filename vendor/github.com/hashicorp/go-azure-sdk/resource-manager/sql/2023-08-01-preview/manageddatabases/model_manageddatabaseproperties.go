package manageddatabases

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagedDatabaseProperties struct {
	AutoCompleteRestore                          *bool                      `json:"autoCompleteRestore,omitempty"`
	CatalogCollation                             *CatalogCollationType      `json:"catalogCollation,omitempty"`
	Collation                                    *string                    `json:"collation,omitempty"`
	CreateMode                                   *ManagedDatabaseCreateMode `json:"createMode,omitempty"`
	CreationDate                                 *string                    `json:"creationDate,omitempty"`
	CrossSubscriptionRestorableDroppedDatabaseId *string                    `json:"crossSubscriptionRestorableDroppedDatabaseId,omitempty"`
	CrossSubscriptionSourceDatabaseId            *string                    `json:"crossSubscriptionSourceDatabaseId,omitempty"`
	CrossSubscriptionTargetManagedInstanceId     *string                    `json:"crossSubscriptionTargetManagedInstanceId,omitempty"`
	DefaultSecondaryLocation                     *string                    `json:"defaultSecondaryLocation,omitempty"`
	EarliestRestorePoint                         *string                    `json:"earliestRestorePoint,omitempty"`
	FailoverGroupId                              *string                    `json:"failoverGroupId,omitempty"`
	IsLedgerOn                                   *bool                      `json:"isLedgerOn,omitempty"`
	LastBackupName                               *string                    `json:"lastBackupName,omitempty"`
	LongTermRetentionBackupResourceId            *string                    `json:"longTermRetentionBackupResourceId,omitempty"`
	RecoverableDatabaseId                        *string                    `json:"recoverableDatabaseId,omitempty"`
	RestorableDroppedDatabaseId                  *string                    `json:"restorableDroppedDatabaseId,omitempty"`
	RestorePointInTime                           *string                    `json:"restorePointInTime,omitempty"`
	SourceDatabaseId                             *string                    `json:"sourceDatabaseId,omitempty"`
	Status                                       *ManagedDatabaseStatus     `json:"status,omitempty"`
	StorageContainerIdentity                     *string                    `json:"storageContainerIdentity,omitempty"`
	StorageContainerSasToken                     *string                    `json:"storageContainerSasToken,omitempty"`
	StorageContainerUri                          *string                    `json:"storageContainerUri,omitempty"`
}

func (o *ManagedDatabaseProperties) GetCreationDateAsTime() (*time.Time, error) {
	if o.CreationDate == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreationDate, "2006-01-02T15:04:05Z07:00")
}

func (o *ManagedDatabaseProperties) SetCreationDateAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreationDate = &formatted
}

func (o *ManagedDatabaseProperties) GetEarliestRestorePointAsTime() (*time.Time, error) {
	if o.EarliestRestorePoint == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.EarliestRestorePoint, "2006-01-02T15:04:05Z07:00")
}

func (o *ManagedDatabaseProperties) SetEarliestRestorePointAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.EarliestRestorePoint = &formatted
}

func (o *ManagedDatabaseProperties) GetRestorePointInTimeAsTime() (*time.Time, error) {
	if o.RestorePointInTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.RestorePointInTime, "2006-01-02T15:04:05Z07:00")
}

func (o *ManagedDatabaseProperties) SetRestorePointInTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.RestorePointInTime = &formatted
}
