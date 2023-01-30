package workflowtriggers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RecurrenceScheduleOccurrence struct {
	Day        *DayOfWeek `json:"day,omitempty"`
	Occurrence *int64     `json:"occurrence,omitempty"`
}
