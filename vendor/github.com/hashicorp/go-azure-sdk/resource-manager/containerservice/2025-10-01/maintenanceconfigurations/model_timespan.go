package maintenanceconfigurations

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TimeSpan struct {
	End   *string `json:"end,omitempty"`
	Start *string `json:"start,omitempty"`
}

func (o *TimeSpan) GetEndAsTime() (*time.Time, error) {
	if o.End == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.End, "2006-01-02T15:04:05Z07:00")
}

func (o *TimeSpan) SetEndAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.End = &formatted
}

func (o *TimeSpan) GetStartAsTime() (*time.Time, error) {
	if o.Start == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.Start, "2006-01-02T15:04:05Z07:00")
}

func (o *TimeSpan) SetStartAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.Start = &formatted
}
