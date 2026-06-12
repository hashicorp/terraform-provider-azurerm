package webhooks

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EventContent struct {
	Action    *string  `json:"action,omitempty"`
	Actor     *Actor   `json:"actor,omitempty"`
	Id        *string  `json:"id,omitempty"`
	Request   *Request `json:"request,omitempty"`
	Source    *Source  `json:"source,omitempty"`
	Target    *Target  `json:"target,omitempty"`
	Timestamp *string  `json:"timestamp,omitempty"`
}

func (o *EventContent) GetTimestampAsTime() (*time.Time, error) {
	if o.Timestamp == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.Timestamp, "2006-01-02T15:04:05Z07:00")
}

func (o *EventContent) SetTimestampAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.Timestamp = &formatted
}
