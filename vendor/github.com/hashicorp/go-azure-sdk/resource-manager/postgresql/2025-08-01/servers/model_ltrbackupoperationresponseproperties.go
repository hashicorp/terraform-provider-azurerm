package servers

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LtrBackupOperationResponseProperties struct {
	BackupMetadata         *string         `json:"backupMetadata,omitempty"`
	BackupName             *string         `json:"backupName,omitempty"`
	DataTransferredInBytes *int64          `json:"dataTransferredInBytes,omitempty"`
	DatasourceSizeInBytes  *int64          `json:"datasourceSizeInBytes,omitempty"`
	EndTime                *string         `json:"endTime,omitempty"`
	ErrorCode              *string         `json:"errorCode,omitempty"`
	ErrorMessage           *string         `json:"errorMessage,omitempty"`
	PercentComplete        *float64        `json:"percentComplete,omitempty"`
	StartTime              string          `json:"startTime"`
	Status                 ExecutionStatus `json:"status"`
}

func (o *LtrBackupOperationResponseProperties) GetEndTimeAsTime() (*time.Time, error) {
	if o.EndTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.EndTime, "2006-01-02T15:04:05Z07:00")
}

func (o *LtrBackupOperationResponseProperties) SetEndTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.EndTime = &formatted
}

func (o *LtrBackupOperationResponseProperties) GetStartTimeAsTime() (*time.Time, error) {
	return dates.ParseAsFormat(&o.StartTime, "2006-01-02T15:04:05Z07:00")
}

func (o *LtrBackupOperationResponseProperties) SetStartTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.StartTime = formatted
}
