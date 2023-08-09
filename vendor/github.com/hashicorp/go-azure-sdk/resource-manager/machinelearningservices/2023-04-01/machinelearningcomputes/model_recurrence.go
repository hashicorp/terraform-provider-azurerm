package machinelearningcomputes

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Recurrence struct {
	Frequency *RecurrenceFrequency `json:"frequency,omitempty"`
	Interval  *int64               `json:"interval,omitempty"`
	Schedule  *RecurrenceSchedule  `json:"schedule,omitempty"`
	StartTime *string              `json:"startTime,omitempty"`
	TimeZone  *string              `json:"timeZone,omitempty"`
}
