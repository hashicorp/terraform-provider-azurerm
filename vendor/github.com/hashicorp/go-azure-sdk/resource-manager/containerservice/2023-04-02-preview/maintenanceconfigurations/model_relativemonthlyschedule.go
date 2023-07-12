package maintenanceconfigurations

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RelativeMonthlySchedule struct {
	DayOfWeek      WeekDay `json:"dayOfWeek"`
	IntervalMonths int64   `json:"intervalMonths"`
	WeekIndex      Type    `json:"weekIndex"`
}
