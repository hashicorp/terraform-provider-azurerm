package maintenanceconfigurations

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Schedule struct {
	AbsoluteMonthly *AbsoluteMonthlySchedule `json:"absoluteMonthly,omitempty"`
	Daily           *DailySchedule           `json:"daily,omitempty"`
	RelativeMonthly *RelativeMonthlySchedule `json:"relativeMonthly,omitempty"`
	Weekly          *WeeklySchedule          `json:"weekly,omitempty"`
}
