package cosmosdb

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BackupPolicyMigrationState struct {
	StartTime  *string                      `json:"startTime,omitempty"`
	Status     *BackupPolicyMigrationStatus `json:"status,omitempty"`
	TargetType *BackupPolicyType            `json:"targetType,omitempty"`
}

func (o *BackupPolicyMigrationState) GetStartTimeAsTime() (*time.Time, error) {
	if o.StartTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.StartTime, "2006-01-02T15:04:05Z07:00")
}

func (o *BackupPolicyMigrationState) SetStartTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.StartTime = &formatted
}
