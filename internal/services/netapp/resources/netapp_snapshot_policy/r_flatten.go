// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package netapp_snapshot_policy

import (
	"strconv"
	"strings"

	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2026-05-01/snapshotpolicies"
)

func flattenNetAppVolumeSnapshotPolicyHourlySchedule(input *snapshotpolicies.HourlySchedule) []any {
	if input == nil || (input.Minute == nil && input.SnapshotsToKeep == nil) {
		return []any{}
	}

	return []any{
		map[string]any{
			"snapshots_to_keep": input.SnapshotsToKeep,
			"minute":            input.Minute,
		},
	}
}

func flattenNetAppVolumeSnapshotPolicyDailySchedule(input *snapshotpolicies.DailySchedule) []any {
	if input == nil || (input.SnapshotsToKeep == nil && input.Hour == nil && input.Minute == nil) {
		return []any{}
	}

	return []any{
		map[string]any{
			"snapshots_to_keep": input.SnapshotsToKeep,
			"hour":              input.Hour,
			"minute":            input.Minute,
		},
	}
}

func flattenNetAppVolumeSnapshotPolicyWeeklySchedule(input *snapshotpolicies.WeeklySchedule) []any {
	if input == nil || (input.SnapshotsToKeep == nil && input.Day == nil && input.Hour == nil && input.Minute == nil) {
		return []any{}
	}

	weekDays := make([]any, 0)
	if input.Day != nil {
		for day := range strings.SplitSeq(*input.Day, ",") {
			weekDays = append(weekDays, day)
		}
	}

	return []any{
		map[string]any{
			"snapshots_to_keep": input.SnapshotsToKeep,
			"days_of_week":      weekDays,
			"hour":              input.Hour,
			"minute":            input.Minute,
		},
	}
}

func flattenNetAppVolumeSnapshotPolicyMonthlySchedule(input *snapshotpolicies.MonthlySchedule) []any {
	if input == nil || (input.SnapshotsToKeep == nil && input.DaysOfMonth == nil && input.Hour == nil && input.Minute == nil) {
		return []any{}
	}

	daysOfMonth := make([]any, 0)
	if input.DaysOfMonth != nil {
		for day := range strings.SplitSeq(*input.DaysOfMonth, ",") {
			intDay, _ := strconv.Atoi(day)
			daysOfMonth = append(daysOfMonth, intDay)
		}
	}

	return []any{
		map[string]any{
			"snapshots_to_keep": input.SnapshotsToKeep,
			"days_of_month":     daysOfMonth,
			"hour":              input.Hour,
			"minute":            input.Minute,
		},
	}
}
