package backups

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServerBackupProperties struct {
	BackupType    *string `json:"backupType,omitempty"`
	CompletedTime *string `json:"completedTime,omitempty"`
	Source        *string `json:"source,omitempty"`
}

func (o *ServerBackupProperties) GetCompletedTimeAsTime() (*time.Time, error) {
	if o.CompletedTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CompletedTime, "2006-01-02T15:04:05Z07:00")
}

func (o *ServerBackupProperties) SetCompletedTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CompletedTime = &formatted
}
