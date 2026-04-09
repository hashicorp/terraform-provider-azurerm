package tables

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RestoredLogs struct {
	AzureAsyncOperationId *string `json:"azureAsyncOperationId,omitempty"`
	EndRestoreTime        *string `json:"endRestoreTime,omitempty"`
	SourceTable           *string `json:"sourceTable,omitempty"`
	StartRestoreTime      *string `json:"startRestoreTime,omitempty"`
}

func (o *RestoredLogs) GetEndRestoreTimeAsTime() (*time.Time, error) {
	if o.EndRestoreTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.EndRestoreTime, "2006-01-02T15:04:05Z07:00")
}

func (o *RestoredLogs) SetEndRestoreTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.EndRestoreTime = &formatted
}

func (o *RestoredLogs) GetStartRestoreTimeAsTime() (*time.Time, error) {
	if o.StartRestoreTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.StartRestoreTime, "2006-01-02T15:04:05Z07:00")
}

func (o *RestoredLogs) SetStartRestoreTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.StartRestoreTime = &formatted
}
