package webapps

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MSDeployLogEntry struct {
	Message *string               `json:"message,omitempty"`
	Time    *string               `json:"time,omitempty"`
	Type    *MSDeployLogEntryType `json:"type,omitempty"`
}

func (o *MSDeployLogEntry) GetTimeAsTime() (*time.Time, error) {
	if o.Time == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.Time, "2006-01-02T15:04:05Z07:00")
}

func (o *MSDeployLogEntry) SetTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.Time = &formatted
}
