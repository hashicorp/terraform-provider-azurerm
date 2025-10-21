package softwareupdateconfiguration

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AdvancedSchedule struct {
	MonthDays          *[]int64                             `json:"monthDays,omitempty"`
	MonthlyOccurrences *[]AdvancedScheduleMonthlyOccurrence `json:"monthlyOccurrences,omitempty"`
	WeekDays           *[]string                            `json:"weekDays,omitempty"`
}
