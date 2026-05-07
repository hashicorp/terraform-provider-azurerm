package startstopmanagedinstanceschedules

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ScheduleItem struct {
	StartDay  DayOfWeek `json:"startDay"`
	StartTime string    `json:"startTime"`
	StopDay   DayOfWeek `json:"stopDay"`
	StopTime  string    `json:"stopTime"`
}
