package exports

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExportRecurrencePeriod struct {
	From string  `json:"from"`
	To   *string `json:"to,omitempty"`
}

func (o *ExportRecurrencePeriod) GetFromAsTime() (*time.Time, error) {
	return dates.ParseAsFormat(&o.From, "2006-01-02T15:04:05Z07:00")
}

func (o *ExportRecurrencePeriod) SetFromAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.From = formatted
}

func (o *ExportRecurrencePeriod) GetToAsTime() (*time.Time, error) {
	if o.To == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.To, "2006-01-02T15:04:05Z07:00")
}

func (o *ExportRecurrencePeriod) SetToAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.To = &formatted
}
