package virtualmachines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ScheduleCreationParameterProperties struct {
	DailyRecurrence      *DayDetails           `json:"dailyRecurrence,omitempty"`
	HourlyRecurrence     *HourDetails          `json:"hourlyRecurrence,omitempty"`
	NotificationSettings *NotificationSettings `json:"notificationSettings,omitempty"`
	Status               *EnableStatus         `json:"status,omitempty"`
	TargetResourceId     *string               `json:"targetResourceId,omitempty"`
	TaskType             *string               `json:"taskType,omitempty"`
	TimeZoneId           *string               `json:"timeZoneId,omitempty"`
	WeeklyRecurrence     *WeekDetails          `json:"weeklyRecurrence,omitempty"`
}
