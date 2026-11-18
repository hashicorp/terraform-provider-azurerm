package webapps

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BackupSchedule struct {
	FrequencyInterval     int64         `json:"frequencyInterval"`
	FrequencyUnit         FrequencyUnit `json:"frequencyUnit"`
	KeepAtLeastOneBackup  bool          `json:"keepAtLeastOneBackup"`
	LastExecutionTime     *string       `json:"lastExecutionTime,omitempty"`
	RetentionPeriodInDays int64         `json:"retentionPeriodInDays"`
	StartTime             *string       `json:"startTime,omitempty"`
}

func (o *BackupSchedule) GetLastExecutionTimeAsTime() (*time.Time, error) {
	if o.LastExecutionTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastExecutionTime, "2006-01-02T15:04:05Z07:00")
}

func (o *BackupSchedule) SetLastExecutionTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastExecutionTime = &formatted
}

func (o *BackupSchedule) GetStartTimeAsTime() (*time.Time, error) {
	if o.StartTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.StartTime, "2006-01-02T15:04:05Z07:00")
}

func (o *BackupSchedule) SetStartTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.StartTime = &formatted
}
