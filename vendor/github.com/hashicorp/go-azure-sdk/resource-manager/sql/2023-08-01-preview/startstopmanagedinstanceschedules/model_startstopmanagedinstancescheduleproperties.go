package startstopmanagedinstanceschedules

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StartStopManagedInstanceScheduleProperties struct {
	Description       *string        `json:"description,omitempty"`
	NextExecutionTime *string        `json:"nextExecutionTime,omitempty"`
	NextRunAction     *string        `json:"nextRunAction,omitempty"`
	ScheduleList      []ScheduleItem `json:"scheduleList"`
	TimeZoneId        *string        `json:"timeZoneId,omitempty"`
}
