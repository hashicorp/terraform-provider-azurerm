package schedules

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ScheduleUpdateProperties struct {
	Frequency *ScheduledFrequency   `json:"frequency,omitempty"`
	State     *ScheduleEnableStatus `json:"state,omitempty"`
	Time      *string               `json:"time,omitempty"`
	TimeZone  *string               `json:"timeZone,omitempty"`
	Type      *ScheduledType        `json:"type,omitempty"`
}
