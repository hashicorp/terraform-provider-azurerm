package integrationaccounts

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TrackingEvent struct {
	Error      *TrackingEventErrorInfo `json:"error,omitempty"`
	EventLevel EventLevel              `json:"eventLevel"`
	EventTime  string                  `json:"eventTime"`
	Record     *interface{}            `json:"record,omitempty"`
	RecordType TrackingRecordType      `json:"recordType"`
}

func (o *TrackingEvent) GetEventTimeAsTime() (*time.Time, error) {
	return dates.ParseAsFormat(&o.EventTime, "2006-01-02T15:04:05Z07:00")
}

func (o *TrackingEvent) SetEventTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.EventTime = formatted
}
