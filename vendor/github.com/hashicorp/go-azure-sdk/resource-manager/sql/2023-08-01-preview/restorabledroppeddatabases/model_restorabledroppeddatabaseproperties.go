package restorabledroppeddatabases

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RestorableDroppedDatabaseProperties struct {
	BackupStorageRedundancy *BackupStorageRedundancy `json:"backupStorageRedundancy,omitempty"`
	CreationDate            *string                  `json:"creationDate,omitempty"`
	DatabaseName            *string                  `json:"databaseName,omitempty"`
	DeletionDate            *string                  `json:"deletionDate,omitempty"`
	EarliestRestoreDate     *string                  `json:"earliestRestoreDate,omitempty"`
	Keys                    *map[string]DatabaseKey  `json:"keys,omitempty"`
	MaxSizeBytes            *int64                   `json:"maxSizeBytes,omitempty"`
}

func (o *RestorableDroppedDatabaseProperties) GetCreationDateAsTime() (*time.Time, error) {
	if o.CreationDate == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreationDate, "2006-01-02T15:04:05Z07:00")
}

func (o *RestorableDroppedDatabaseProperties) SetCreationDateAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreationDate = &formatted
}

func (o *RestorableDroppedDatabaseProperties) GetDeletionDateAsTime() (*time.Time, error) {
	if o.DeletionDate == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.DeletionDate, "2006-01-02T15:04:05Z07:00")
}

func (o *RestorableDroppedDatabaseProperties) SetDeletionDateAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.DeletionDate = &formatted
}

func (o *RestorableDroppedDatabaseProperties) GetEarliestRestoreDateAsTime() (*time.Time, error) {
	if o.EarliestRestoreDate == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.EarliestRestoreDate, "2006-01-02T15:04:05Z07:00")
}

func (o *RestorableDroppedDatabaseProperties) SetEarliestRestoreDateAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.EarliestRestoreDate = &formatted
}
