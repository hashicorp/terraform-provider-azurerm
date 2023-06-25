package share

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ShareSynchronization struct {
	ConsumerEmail       *string              `json:"consumerEmail,omitempty"`
	ConsumerName        *string              `json:"consumerName,omitempty"`
	ConsumerTenantName  *string              `json:"consumerTenantName,omitempty"`
	DurationMs          *int64               `json:"durationMs,omitempty"`
	EndTime             *string              `json:"endTime,omitempty"`
	Message             *string              `json:"message,omitempty"`
	StartTime           *string              `json:"startTime,omitempty"`
	Status              *string              `json:"status,omitempty"`
	SynchronizationId   *string              `json:"synchronizationId,omitempty"`
	SynchronizationMode *SynchronizationMode `json:"synchronizationMode,omitempty"`
}

func (o *ShareSynchronization) GetEndTimeAsTime() (*time.Time, error) {
	if o.EndTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.EndTime, "2006-01-02T15:04:05Z07:00")
}

func (o *ShareSynchronization) SetEndTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.EndTime = &formatted
}

func (o *ShareSynchronization) GetStartTimeAsTime() (*time.Time, error) {
	if o.StartTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.StartTime, "2006-01-02T15:04:05Z07:00")
}

func (o *ShareSynchronization) SetStartTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.StartTime = &formatted
}
