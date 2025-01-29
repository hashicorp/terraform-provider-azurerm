package maintenanceconfigurations

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MaintenanceWindow struct {
	DurationHours   int64       `json:"durationHours"`
	NotAllowedDates *[]DateSpan `json:"notAllowedDates,omitempty"`
	Schedule        Schedule    `json:"schedule"`
	StartDate       *string     `json:"startDate,omitempty"`
	StartTime       string      `json:"startTime"`
	UtcOffset       *string     `json:"utcOffset,omitempty"`
}
