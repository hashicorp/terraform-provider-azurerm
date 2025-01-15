package hostpool

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MaintenanceWindowProperties struct {
	DayOfWeek *DayOfWeek `json:"dayOfWeek,omitempty"`
	Hour      *int64     `json:"hour,omitempty"`
}
