package autonomousdatabases

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LongTermBackUpScheduleDetails struct {
	IsDisabled            *bool              `json:"isDisabled,omitempty"`
	RepeatCadence         *RepeatCadenceType `json:"repeatCadence,omitempty"`
	RetentionPeriodInDays *int64             `json:"retentionPeriodInDays,omitempty"`
	TimeOfBackup          *string            `json:"timeOfBackup,omitempty"`
}

func (o *LongTermBackUpScheduleDetails) GetTimeOfBackupAsTime() (*time.Time, error) {
	if o.TimeOfBackup == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.TimeOfBackup, "2006-01-02T15:04:05Z07:00")
}

func (o *LongTermBackUpScheduleDetails) SetTimeOfBackupAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.TimeOfBackup = &formatted
}
