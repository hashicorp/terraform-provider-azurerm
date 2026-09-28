package backups

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BackupProperties struct {
	BackupId               *string     `json:"backupId,omitempty"`
	BackupPolicyResourceId *string     `json:"backupPolicyResourceId,omitempty"`
	BackupType             *BackupType `json:"backupType,omitempty"`
	CompletionDate         *string     `json:"completionDate,omitempty"`
	CreationDate           *string     `json:"creationDate,omitempty"`
	FailureReason          *string     `json:"failureReason,omitempty"`
	IsLargeVolume          *bool       `json:"isLargeVolume,omitempty"`
	Label                  *string     `json:"label,omitempty"`
	ProvisioningState      *string     `json:"provisioningState,omitempty"`
	Size                   *int64      `json:"size,omitempty"`
	SnapshotCreationDate   *string     `json:"snapshotCreationDate,omitempty"`
	SnapshotName           *string     `json:"snapshotName,omitempty"`
	UseExistingSnapshot    *bool       `json:"useExistingSnapshot,omitempty"`
	VolumeResourceId       string      `json:"volumeResourceId"`
}

func (o *BackupProperties) GetCompletionDateAsTime() (*time.Time, error) {
	if o.CompletionDate == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CompletionDate, "2006-01-02T15:04:05Z07:00")
}

func (o *BackupProperties) SetCompletionDateAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CompletionDate = &formatted
}

func (o *BackupProperties) GetCreationDateAsTime() (*time.Time, error) {
	if o.CreationDate == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreationDate, "2006-01-02T15:04:05Z07:00")
}

func (o *BackupProperties) SetCreationDateAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreationDate = &formatted
}

func (o *BackupProperties) GetSnapshotCreationDateAsTime() (*time.Time, error) {
	if o.SnapshotCreationDate == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.SnapshotCreationDate, "2006-01-02T15:04:05Z07:00")
}

func (o *BackupProperties) SetSnapshotCreationDateAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.SnapshotCreationDate = &formatted
}
