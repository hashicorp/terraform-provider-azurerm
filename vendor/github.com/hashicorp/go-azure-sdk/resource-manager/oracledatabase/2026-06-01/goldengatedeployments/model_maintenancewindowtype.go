package goldengatedeployments

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MaintenanceWindowType struct {
	Day       *DayOfWeekName `json:"day,omitempty"`
	StartHour *int64         `json:"startHour,omitempty"`
}
