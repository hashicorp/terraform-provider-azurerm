package protectionpolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WeeklyRetentionSchedule struct {
	DaysOfTheWeek     *[]DayOfWeek       `json:"daysOfTheWeek,omitempty"`
	RetentionDuration *RetentionDuration `json:"retentionDuration,omitempty"`
	RetentionTimes    *[]string          `json:"retentionTimes,omitempty"`
}
