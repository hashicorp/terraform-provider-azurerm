package sqlvirtualmachines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Schedule struct {
	DayOfWeek         *AssessmentDayOfWeek `json:"dayOfWeek,omitempty"`
	Enable            *bool                `json:"enable,omitempty"`
	MonthlyOccurrence *int64               `json:"monthlyOccurrence,omitempty"`
	StartTime         *string              `json:"startTime,omitempty"`
	WeeklyInterval    *int64               `json:"weeklyInterval,omitempty"`
}
