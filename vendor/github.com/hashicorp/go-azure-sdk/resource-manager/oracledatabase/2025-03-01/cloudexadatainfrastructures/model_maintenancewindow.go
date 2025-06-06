package cloudexadatainfrastructures

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MaintenanceWindow struct {
	CustomActionTimeoutInMins    *int64        `json:"customActionTimeoutInMins,omitempty"`
	DaysOfWeek                   *[]DayOfWeek  `json:"daysOfWeek,omitempty"`
	HoursOfDay                   *[]int64      `json:"hoursOfDay,omitempty"`
	IsCustomActionTimeoutEnabled *bool         `json:"isCustomActionTimeoutEnabled,omitempty"`
	IsMonthlyPatchingEnabled     *bool         `json:"isMonthlyPatchingEnabled,omitempty"`
	LeadTimeInWeeks              *int64        `json:"leadTimeInWeeks,omitempty"`
	Months                       *[]Month      `json:"months,omitempty"`
	PatchingMode                 *PatchingMode `json:"patchingMode,omitempty"`
	Preference                   *Preference   `json:"preference,omitempty"`
	WeeksOfMonth                 *[]int64      `json:"weeksOfMonth,omitempty"`
}
