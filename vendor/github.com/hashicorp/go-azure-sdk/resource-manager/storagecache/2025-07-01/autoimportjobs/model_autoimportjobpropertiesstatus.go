package autoimportjobs

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AutoImportJobPropertiesStatus struct {
	BlobSyncEvents         *AutoImportJobPropertiesStatusBlobSyncEvents `json:"blobSyncEvents,omitempty"`
	ImportedDirectories    *int64                                       `json:"importedDirectories,omitempty"`
	ImportedFiles          *int64                                       `json:"importedFiles,omitempty"`
	ImportedSymlinks       *int64                                       `json:"importedSymlinks,omitempty"`
	LastCompletionTimeUTC  *string                                      `json:"lastCompletionTimeUTC,omitempty"`
	LastStartedTimeUTC     *string                                      `json:"lastStartedTimeUTC,omitempty"`
	PreexistingDirectories *int64                                       `json:"preexistingDirectories,omitempty"`
	PreexistingFiles       *int64                                       `json:"preexistingFiles,omitempty"`
	PreexistingSymlinks    *int64                                       `json:"preexistingSymlinks,omitempty"`
	RateOfBlobImport       *int64                                       `json:"rateOfBlobImport,omitempty"`
	RateOfBlobWalk         *int64                                       `json:"rateOfBlobWalk,omitempty"`
	ScanEndTime            *string                                      `json:"scanEndTime,omitempty"`
	ScanStartTime          *string                                      `json:"scanStartTime,omitempty"`
	State                  *AutoImportJobState                          `json:"state,omitempty"`
	StatusCode             *string                                      `json:"statusCode,omitempty"`
	StatusMessage          *string                                      `json:"statusMessage,omitempty"`
	TotalBlobsImported     *int64                                       `json:"totalBlobsImported,omitempty"`
	TotalBlobsWalked       *int64                                       `json:"totalBlobsWalked,omitempty"`
	TotalConflicts         *int64                                       `json:"totalConflicts,omitempty"`
	TotalErrors            *int64                                       `json:"totalErrors,omitempty"`
}

func (o *AutoImportJobPropertiesStatus) GetLastCompletionTimeUTCAsTime() (*time.Time, error) {
	if o.LastCompletionTimeUTC == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastCompletionTimeUTC, "2006-01-02T15:04:05Z07:00")
}

func (o *AutoImportJobPropertiesStatus) SetLastCompletionTimeUTCAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastCompletionTimeUTC = &formatted
}

func (o *AutoImportJobPropertiesStatus) GetLastStartedTimeUTCAsTime() (*time.Time, error) {
	if o.LastStartedTimeUTC == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastStartedTimeUTC, "2006-01-02T15:04:05Z07:00")
}

func (o *AutoImportJobPropertiesStatus) SetLastStartedTimeUTCAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastStartedTimeUTC = &formatted
}

func (o *AutoImportJobPropertiesStatus) GetScanEndTimeAsTime() (*time.Time, error) {
	if o.ScanEndTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.ScanEndTime, "2006-01-02T15:04:05Z07:00")
}

func (o *AutoImportJobPropertiesStatus) SetScanEndTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.ScanEndTime = &formatted
}

func (o *AutoImportJobPropertiesStatus) GetScanStartTimeAsTime() (*time.Time, error) {
	if o.ScanStartTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.ScanStartTime, "2006-01-02T15:04:05Z07:00")
}

func (o *AutoImportJobPropertiesStatus) SetScanStartTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.ScanStartTime = &formatted
}
