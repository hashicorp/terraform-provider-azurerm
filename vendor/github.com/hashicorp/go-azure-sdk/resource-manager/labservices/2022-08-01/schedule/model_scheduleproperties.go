package schedule

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ScheduleProperties struct {
	Notes             *string            `json:"notes,omitempty"`
	ProvisioningState *ProvisioningState `json:"provisioningState,omitempty"`
	RecurrencePattern *RecurrencePattern `json:"recurrencePattern,omitempty"`
	StartAt           *string            `json:"startAt,omitempty"`
	StopAt            *string            `json:"stopAt,omitempty"`
	TimeZoneId        *string            `json:"timeZoneId,omitempty"`
}

func (o *ScheduleProperties) GetStartAtAsTime() (*time.Time, error) {
	if o.StartAt == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.StartAt, "2006-01-02T15:04:05Z07:00")
}

func (o *ScheduleProperties) SetStartAtAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.StartAt = &formatted
}

func (o *ScheduleProperties) GetStopAtAsTime() (*time.Time, error) {
	if o.StopAt == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.StopAt, "2006-01-02T15:04:05Z07:00")
}

func (o *ScheduleProperties) SetStopAtAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.StopAt = &formatted
}
