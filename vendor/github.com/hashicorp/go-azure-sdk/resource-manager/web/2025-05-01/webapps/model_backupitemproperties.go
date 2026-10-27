package webapps

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BackupItemProperties struct {
	BlobName             *string                  `json:"blobName,omitempty"`
	CorrelationId        *string                  `json:"correlationId,omitempty"`
	Created              *string                  `json:"created,omitempty"`
	Databases            *[]DatabaseBackupSetting `json:"databases,omitempty"`
	FinishedTimeStamp    *string                  `json:"finishedTimeStamp,omitempty"`
	Id                   *int64                   `json:"id,omitempty"`
	LastRestoreTimeStamp *string                  `json:"lastRestoreTimeStamp,omitempty"`
	Log                  *string                  `json:"log,omitempty"`
	Name                 *string                  `json:"name,omitempty"`
	Scheduled            *bool                    `json:"scheduled,omitempty"`
	SizeInBytes          *int64                   `json:"sizeInBytes,omitempty"`
	Status               *BackupItemStatus        `json:"status,omitempty"`
	StorageAccountURL    *string                  `json:"storageAccountUrl,omitempty"`
	WebsiteSizeInBytes   *int64                   `json:"websiteSizeInBytes,omitempty"`
}

func (o *BackupItemProperties) GetCreatedAsTime() (*time.Time, error) {
	if o.Created == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.Created, "2006-01-02T15:04:05Z07:00")
}

func (o *BackupItemProperties) SetCreatedAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.Created = &formatted
}

func (o *BackupItemProperties) GetFinishedTimeStampAsTime() (*time.Time, error) {
	if o.FinishedTimeStamp == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.FinishedTimeStamp, "2006-01-02T15:04:05Z07:00")
}

func (o *BackupItemProperties) SetFinishedTimeStampAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.FinishedTimeStamp = &formatted
}

func (o *BackupItemProperties) GetLastRestoreTimeStampAsTime() (*time.Time, error) {
	if o.LastRestoreTimeStamp == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastRestoreTimeStamp, "2006-01-02T15:04:05Z07:00")
}

func (o *BackupItemProperties) SetLastRestoreTimeStampAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastRestoreTimeStamp = &formatted
}
