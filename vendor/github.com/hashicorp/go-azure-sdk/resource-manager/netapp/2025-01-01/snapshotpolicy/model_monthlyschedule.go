package snapshotpolicy

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MonthlySchedule struct {
	DaysOfMonth     *string `json:"daysOfMonth,omitempty"`
	Hour            *int64  `json:"hour,omitempty"`
	Minute          *int64  `json:"minute,omitempty"`
	SnapshotsToKeep *int64  `json:"snapshotsToKeep,omitempty"`
	UsedBytes       *int64  `json:"usedBytes,omitempty"`
}
