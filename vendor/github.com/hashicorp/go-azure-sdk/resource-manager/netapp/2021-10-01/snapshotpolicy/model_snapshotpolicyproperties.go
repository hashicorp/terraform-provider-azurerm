package snapshotpolicy

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SnapshotPolicyProperties struct {
	DailySchedule     *DailySchedule   `json:"dailySchedule"`
	Enabled           *bool            `json:"enabled,omitempty"`
	HourlySchedule    *HourlySchedule  `json:"hourlySchedule"`
	MonthlySchedule   *MonthlySchedule `json:"monthlySchedule"`
	ProvisioningState *string          `json:"provisioningState,omitempty"`
	WeeklySchedule    *WeeklySchedule  `json:"weeklySchedule"`
}
