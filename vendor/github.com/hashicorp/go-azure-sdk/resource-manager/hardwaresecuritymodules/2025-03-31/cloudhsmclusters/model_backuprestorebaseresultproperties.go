package cloudhsmclusters

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BackupRestoreBaseResultProperties struct {
	EndTime       *string                       `json:"endTime,omitempty"`
	Error         *ErrorDetail                  `json:"error,omitempty"`
	JobId         *string                       `json:"jobId,omitempty"`
	StartTime     *string                       `json:"startTime,omitempty"`
	Status        *BackupRestoreOperationStatus `json:"status,omitempty"`
	StatusDetails *string                       `json:"statusDetails,omitempty"`
}

func (o *BackupRestoreBaseResultProperties) GetEndTimeAsTime() (*time.Time, error) {
	if o.EndTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.EndTime, "2006-01-02T15:04:05Z07:00")
}

func (o *BackupRestoreBaseResultProperties) SetEndTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.EndTime = &formatted
}

func (o *BackupRestoreBaseResultProperties) GetStartTimeAsTime() (*time.Time, error) {
	if o.StartTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.StartTime, "2006-01-02T15:04:05Z07:00")
}

func (o *BackupRestoreBaseResultProperties) SetStartTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.StartTime = &formatted
}
