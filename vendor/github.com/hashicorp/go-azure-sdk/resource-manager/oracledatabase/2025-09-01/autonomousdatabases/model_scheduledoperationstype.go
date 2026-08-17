package autonomousdatabases

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ScheduledOperationsType struct {
	DayOfWeek          DayOfWeek `json:"dayOfWeek"`
	ScheduledStartTime *string   `json:"scheduledStartTime,omitempty"`
	ScheduledStopTime  *string   `json:"scheduledStopTime,omitempty"`
}
