package snapshotpolicy

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SnapshotPolicyProperties struct {
	DailySchedule     *DailySchedule   `json:"dailySchedule,omitempty"`
	Enabled           *bool            `json:"enabled,omitempty"`
	HourlySchedule    *HourlySchedule  `json:"hourlySchedule,omitempty"`
	MonthlySchedule   *MonthlySchedule `json:"monthlySchedule,omitempty"`
	ProvisioningState *string          `json:"provisioningState,omitempty"`
	WeeklySchedule    *WeeklySchedule  `json:"weeklySchedule,omitempty"`
}
