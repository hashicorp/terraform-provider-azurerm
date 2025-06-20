package maintenanceconfigurations

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TimeInWeek struct {
	Day       *WeekDay `json:"day,omitempty"`
	HourSlots *[]int64 `json:"hourSlots,omitempty"`
}
