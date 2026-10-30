// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package netapp_snapshot_policy

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2026-05-01/snapshotpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func expandNetAppSnapshotPolicyHourlySchedule(input []any) *snapshotpolicies.HourlySchedule {
	if len(input) == 0 || input[0] == nil {
		return &snapshotpolicies.HourlySchedule{}
	}

	hourlyScheduleObject := snapshotpolicies.HourlySchedule{}

	hourlyScheduleRaw := input[0].(map[string]any)

	if v, ok := hourlyScheduleRaw["snapshots_to_keep"]; ok {
		hourlyScheduleObject.SnapshotsToKeep = pointer.To(int64(v.(int)))
	}
	if v, ok := hourlyScheduleRaw["minute"]; ok {
		hourlyScheduleObject.Minute = pointer.To(int64(v.(int)))
	}

	return &hourlyScheduleObject
}

func expandNetAppSnapshotPolicyDailySchedule(input []any) *snapshotpolicies.DailySchedule {
	if len(input) == 0 || input[0] == nil {
		return &snapshotpolicies.DailySchedule{}
	}

	dailyScheduleObject := snapshotpolicies.DailySchedule{}

	dailyScheduleRaw := input[0].(map[string]any)

	if v, ok := dailyScheduleRaw["snapshots_to_keep"]; ok {
		dailyScheduleObject.SnapshotsToKeep = pointer.To(int64(v.(int)))
	}
	if v, ok := dailyScheduleRaw["hour"]; ok {
		dailyScheduleObject.Hour = pointer.To(int64(v.(int)))
	}
	if v, ok := dailyScheduleRaw["minute"]; ok {
		dailyScheduleObject.Minute = pointer.To(int64(v.(int)))
	}

	return &dailyScheduleObject
}

func expandNetAppSnapshotPolicyWeeklySchedule(input []any) *snapshotpolicies.WeeklySchedule {
	if len(input) == 0 || input[0] == nil {
		return &snapshotpolicies.WeeklySchedule{}
	}

	weeklyScheduleObject := snapshotpolicies.WeeklySchedule{}

	weeklyScheduleRaw := input[0].(map[string]any)

	if v, ok := weeklyScheduleRaw["snapshots_to_keep"]; ok {
		weeklyScheduleObject.SnapshotsToKeep = pointer.To(int64(v.(int)))
	}
	if _, ok := weeklyScheduleRaw["days_of_week"]; ok {
		weeklyScheduleObject.Day = pluginsdk.ExpandStringSliceWithDelimiter(weeklyScheduleRaw["days_of_week"].(*pluginsdk.Set).List(), ",")
	}
	if v, ok := weeklyScheduleRaw["hour"]; ok {
		weeklyScheduleObject.Hour = pointer.To(int64(v.(int)))
	}
	if v, ok := weeklyScheduleRaw["minute"]; ok {
		weeklyScheduleObject.Minute = pointer.To(int64(v.(int)))
	}

	return &weeklyScheduleObject
}

func expandNetAppSnapshotPolicyMonthlySchedule(input []any) *snapshotpolicies.MonthlySchedule {
	if len(input) == 0 || input[0] == nil {
		return &snapshotpolicies.MonthlySchedule{}
	}

	monthlyScheduleObject := snapshotpolicies.MonthlySchedule{}

	monthlyScheduleRaw := input[0].(map[string]any)

	if v, ok := monthlyScheduleRaw["snapshots_to_keep"]; ok {
		monthlyScheduleObject.SnapshotsToKeep = pointer.To(int64(v.(int)))
	}
	if _, ok := monthlyScheduleRaw["days_of_month"]; ok {
		monthlyScheduleObject.DaysOfMonth = pluginsdk.ExpandIntSliceWithDelimiter(monthlyScheduleRaw["days_of_month"].(*pluginsdk.Set).List(), ",")
	}
	if v, ok := monthlyScheduleRaw["hour"]; ok {
		monthlyScheduleObject.Hour = pointer.To(int64(v.(int)))
	}
	if v, ok := monthlyScheduleRaw["minute"]; ok {
		monthlyScheduleObject.Minute = pointer.To(int64(v.(int)))
	}

	return &monthlyScheduleObject
}
