package servers

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Backup struct {
	BackupRetentionDays *int64            `json:"backupRetentionDays,omitempty"`
	EarliestRestoreDate *string           `json:"earliestRestoreDate,omitempty"`
	GeoRedundantBackup  *EnableStatusEnum `json:"geoRedundantBackup,omitempty"`
}

func (o *Backup) GetEarliestRestoreDateAsTime() (*time.Time, error) {
	if o.EarliestRestoreDate == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.EarliestRestoreDate, "2006-01-02T15:04:05Z07:00")
}

func (o *Backup) SetEarliestRestoreDateAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.EarliestRestoreDate = &formatted
}
