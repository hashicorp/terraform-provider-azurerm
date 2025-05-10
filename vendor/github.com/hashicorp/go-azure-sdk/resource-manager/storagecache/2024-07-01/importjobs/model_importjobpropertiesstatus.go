package importjobs

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ImportJobPropertiesStatus struct {
	BlobsImportedPerSecond *int64            `json:"blobsImportedPerSecond,omitempty"`
	BlobsWalkedPerSecond   *int64            `json:"blobsWalkedPerSecond,omitempty"`
	ImportedDirectories    *int64            `json:"importedDirectories,omitempty"`
	ImportedFiles          *int64            `json:"importedFiles,omitempty"`
	ImportedSymlinks       *int64            `json:"importedSymlinks,omitempty"`
	LastCompletionTime     *string           `json:"lastCompletionTime,omitempty"`
	LastStartedTime        *string           `json:"lastStartedTime,omitempty"`
	PreexistingDirectories *int64            `json:"preexistingDirectories,omitempty"`
	PreexistingFiles       *int64            `json:"preexistingFiles,omitempty"`
	PreexistingSymlinks    *int64            `json:"preexistingSymlinks,omitempty"`
	State                  *ImportStatusType `json:"state,omitempty"`
	StatusMessage          *string           `json:"statusMessage,omitempty"`
	TotalBlobsImported     *int64            `json:"totalBlobsImported,omitempty"`
	TotalBlobsWalked       *int64            `json:"totalBlobsWalked,omitempty"`
	TotalConflicts         *int64            `json:"totalConflicts,omitempty"`
	TotalErrors            *int64            `json:"totalErrors,omitempty"`
}

func (o *ImportJobPropertiesStatus) GetLastCompletionTimeAsTime() (*time.Time, error) {
	if o.LastCompletionTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastCompletionTime, "2006-01-02T15:04:05Z07:00")
}

func (o *ImportJobPropertiesStatus) SetLastCompletionTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastCompletionTime = &formatted
}

func (o *ImportJobPropertiesStatus) GetLastStartedTimeAsTime() (*time.Time, error) {
	if o.LastStartedTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastStartedTime, "2006-01-02T15:04:05Z07:00")
}

func (o *ImportJobPropertiesStatus) SetLastStartedTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastStartedTime = &formatted
}
