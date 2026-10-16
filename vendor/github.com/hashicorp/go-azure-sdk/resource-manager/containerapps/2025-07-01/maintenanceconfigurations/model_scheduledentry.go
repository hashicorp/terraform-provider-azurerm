package maintenanceconfigurations

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ScheduledEntry struct {
	DurationHours int64   `json:"durationHours"`
	StartHourUtc  int64   `json:"startHourUtc"`
	WeekDay       WeekDay `json:"weekDay"`
}
