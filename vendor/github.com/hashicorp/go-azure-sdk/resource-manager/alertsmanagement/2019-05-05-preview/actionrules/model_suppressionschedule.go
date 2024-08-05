package actionrules

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SuppressionSchedule struct {
	EndDate          *string  `json:"endDate,omitempty"`
	EndTime          *string  `json:"endTime,omitempty"`
	RecurrenceValues *[]int64 `json:"recurrenceValues,omitempty"`
	StartDate        *string  `json:"startDate,omitempty"`
	StartTime        *string  `json:"startTime,omitempty"`
}
