package autoimportjobs

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AutoImportJobPropertiesStatusBlobSyncEvents struct {
	Deletions                       *int64  `json:"deletions,omitempty"`
	ImportedDirectories             *int64  `json:"importedDirectories,omitempty"`
	ImportedFiles                   *int64  `json:"importedFiles,omitempty"`
	ImportedSymlinks                *int64  `json:"importedSymlinks,omitempty"`
	LastChangeFeedEventConsumedTime *string `json:"lastChangeFeedEventConsumedTime,omitempty"`
	LastTimeFullySynchronized       *string `json:"lastTimeFullySynchronized,omitempty"`
	PreexistingDirectories          *int64  `json:"preexistingDirectories,omitempty"`
	PreexistingFiles                *int64  `json:"preexistingFiles,omitempty"`
	PreexistingSymlinks             *int64  `json:"preexistingSymlinks,omitempty"`
	RateOfBlobImport                *int64  `json:"rateOfBlobImport,omitempty"`
	TotalBlobsImported              *int64  `json:"totalBlobsImported,omitempty"`
	TotalConflicts                  *int64  `json:"totalConflicts,omitempty"`
	TotalErrors                     *int64  `json:"totalErrors,omitempty"`
}

func (o *AutoImportJobPropertiesStatusBlobSyncEvents) GetLastChangeFeedEventConsumedTimeAsTime() (*time.Time, error) {
	if o.LastChangeFeedEventConsumedTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastChangeFeedEventConsumedTime, "2006-01-02T15:04:05Z07:00")
}

func (o *AutoImportJobPropertiesStatusBlobSyncEvents) SetLastChangeFeedEventConsumedTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastChangeFeedEventConsumedTime = &formatted
}

func (o *AutoImportJobPropertiesStatusBlobSyncEvents) GetLastTimeFullySynchronizedAsTime() (*time.Time, error) {
	if o.LastTimeFullySynchronized == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastTimeFullySynchronized, "2006-01-02T15:04:05Z07:00")
}

func (o *AutoImportJobPropertiesStatusBlobSyncEvents) SetLastTimeFullySynchronizedAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastTimeFullySynchronized = &formatted
}
