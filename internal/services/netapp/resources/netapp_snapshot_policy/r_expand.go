// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package netapp_snapshot_policy

import (
	"strconv"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2026-05-01/snapshotpolicies"
)

func expandNetAppSnapshotPolicyHourlySchedule(input []HourlyScheduleModel) *snapshotpolicies.HourlySchedule {
	if len(input) == 0 {
		return &snapshotpolicies.HourlySchedule{}
	}

	return &snapshotpolicies.HourlySchedule{
		SnapshotsToKeep: pointer.To(input[0].SnapshotsToKeep),
		Minute:          pointer.To(input[0].Minute),
	}
}

func expandNetAppSnapshotPolicyDailySchedule(input []DailyScheduleModel) *snapshotpolicies.DailySchedule {
	if len(input) == 0 {
		return &snapshotpolicies.DailySchedule{}
	}

	return &snapshotpolicies.DailySchedule{
		SnapshotsToKeep: pointer.To(input[0].SnapshotsToKeep),
		Hour:            pointer.To(input[0].Hour),
		Minute:          pointer.To(input[0].Minute),
	}
}

func expandNetAppSnapshotPolicyWeeklySchedule(input []WeeklyScheduleModel) *snapshotpolicies.WeeklySchedule {
	if len(input) == 0 {
		return &snapshotpolicies.WeeklySchedule{}
	}

	return &snapshotpolicies.WeeklySchedule{
		SnapshotsToKeep: pointer.To(input[0].SnapshotsToKeep),
		Day:             pointer.To(strings.Join(input[0].DaysOfWeek, ",")),
		Hour:            pointer.To(input[0].Hour),
		Minute:          pointer.To(input[0].Minute),
	}
}

func expandNetAppSnapshotPolicyMonthlySchedule(input []MonthlyScheduleModel) *snapshotpolicies.MonthlySchedule {
	if len(input) == 0 {
		return &snapshotpolicies.MonthlySchedule{}
	}

	daysOfMonth := make([]string, 0, len(input[0].DaysOfMonth))
	for _, day := range input[0].DaysOfMonth {
		daysOfMonth = append(daysOfMonth, strconv.FormatInt(day, 10))
	}

	return &snapshotpolicies.MonthlySchedule{
		SnapshotsToKeep: pointer.To(input[0].SnapshotsToKeep),
		DaysOfMonth:     pointer.To(strings.Join(daysOfMonth, ",")),
		Hour:            pointer.To(input[0].Hour),
		Minute:          pointer.To(input[0].Minute),
	}
}
