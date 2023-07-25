package schedule

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RecurrencePattern struct {
	ExpirationDate string              `json:"expirationDate"`
	Frequency      RecurrenceFrequency `json:"frequency"`
	Interval       *int64              `json:"interval,omitempty"`
	WeekDays       *[]WeekDay          `json:"weekDays,omitempty"`
}

func (o *RecurrencePattern) GetExpirationDateAsTime() (*time.Time, error) {
	return dates.ParseAsFormat(&o.ExpirationDate, "2006-01-02T15:04:05Z07:00")
}

func (o *RecurrencePattern) SetExpirationDateAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.ExpirationDate = formatted
}
