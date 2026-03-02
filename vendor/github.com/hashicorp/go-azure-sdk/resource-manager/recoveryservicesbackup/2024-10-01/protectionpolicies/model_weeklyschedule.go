package protectionpolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WeeklySchedule struct {
	ScheduleRunDays  *[]DayOfWeek `json:"scheduleRunDays,omitempty"`
	ScheduleRunTimes *[]string    `json:"scheduleRunTimes,omitempty"`
}
