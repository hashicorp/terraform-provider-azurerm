package protectionpolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MonthlyRetentionSchedule struct {
	RetentionDuration           *RetentionDuration       `json:"retentionDuration,omitempty"`
	RetentionScheduleDaily      *DailyRetentionFormat    `json:"retentionScheduleDaily,omitempty"`
	RetentionScheduleFormatType *RetentionScheduleFormat `json:"retentionScheduleFormatType,omitempty"`
	RetentionScheduleWeekly     *WeeklyRetentionFormat   `json:"retentionScheduleWeekly,omitempty"`
	RetentionTimes              *[]string                `json:"retentionTimes,omitempty"`
}
