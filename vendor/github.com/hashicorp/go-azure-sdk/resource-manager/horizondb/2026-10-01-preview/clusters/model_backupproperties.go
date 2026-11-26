package clusters

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BackupProperties struct {
	EarliestRestorePointUtc *string            `json:"earliestRestorePointUtc,omitempty"`
	GeoRedundancy           GeoRedundancyState `json:"geoRedundancy"`
	RetentionDays           *int64             `json:"retentionDays,omitempty"`
}

func (o *BackupProperties) GetEarliestRestorePointUtcAsTime() (*time.Time, error) {
	if o.EarliestRestorePointUtc == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.EarliestRestorePointUtc, "2006-01-02T15:04:05Z07:00")
}

func (o *BackupProperties) SetEarliestRestorePointUtcAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.EarliestRestorePointUtc = &formatted
}
