package containerinstance

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Event struct {
	Count          *int64  `json:"count,omitempty"`
	FirstTimestamp *string `json:"firstTimestamp,omitempty"`
	LastTimestamp  *string `json:"lastTimestamp,omitempty"`
	Message        *string `json:"message,omitempty"`
	Name           *string `json:"name,omitempty"`
	Type           *string `json:"type,omitempty"`
}

func (o *Event) GetFirstTimestampAsTime() (*time.Time, error) {
	if o.FirstTimestamp == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.FirstTimestamp, "2006-01-02T15:04:05Z07:00")
}

func (o *Event) SetFirstTimestampAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.FirstTimestamp = &formatted
}

func (o *Event) GetLastTimestampAsTime() (*time.Time, error) {
	if o.LastTimestamp == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastTimestamp, "2006-01-02T15:04:05Z07:00")
}

func (o *Event) SetLastTimestampAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastTimestamp = &formatted
}
