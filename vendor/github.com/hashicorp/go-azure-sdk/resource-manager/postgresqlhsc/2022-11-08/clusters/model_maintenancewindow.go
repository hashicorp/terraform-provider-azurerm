package clusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MaintenanceWindow struct {
	CustomWindow *string `json:"customWindow,omitempty"`
	DayOfWeek    *int64  `json:"dayOfWeek,omitempty"`
	StartHour    *int64  `json:"startHour,omitempty"`
	StartMinute  *int64  `json:"startMinute,omitempty"`
}
