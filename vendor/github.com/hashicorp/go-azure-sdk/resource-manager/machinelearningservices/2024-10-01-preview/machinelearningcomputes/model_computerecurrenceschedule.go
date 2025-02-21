package machinelearningcomputes

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ComputeRecurrenceSchedule struct {
	Hours     []int64           `json:"hours"`
	Minutes   []int64           `json:"minutes"`
	MonthDays *[]int64          `json:"monthDays,omitempty"`
	WeekDays  *[]ComputeWeekDay `json:"weekDays,omitempty"`
}
