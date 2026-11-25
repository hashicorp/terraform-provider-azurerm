// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package backup_policy_file_share

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicesbackup/2024-10-01/protectionpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func flattenBackupProtectionPolicyFileShareSchedule(schedule protectionpolicies.SimpleSchedulePolicy) ([]any, error) {
	block := map[string]any{}

	block["frequency"] = pointer.FromEnum(schedule.ScheduleRunFrequency)

	if times := schedule.ScheduleRunTimes; times != nil && len(*times) > 0 {
		policyTime, _ := time.Parse(time.RFC3339, (*times)[0])
		block["time"] = policyTime.Format("15:04")
	}

	if hourly := schedule.HourlySchedule; hourly != nil {
		hourlyBlock := make(map[string]any, 0)
		if hourly.ScheduleWindowStartTime != nil {
			startTime, err := time.Parse(time.RFC3339, *hourly.ScheduleWindowStartTime)
			if err != nil {
				return nil, fmt.Errorf("error parsing schedule window start time: %s", err)
			}
			hourlyBlock["start_time"] = startTime.Format("15:04")
		}
		hourlyBlock["interval"] = pointer.From(hourly.Interval)
		hourlyBlock["window_duration"] = pointer.From(hourly.ScheduleWindowDuration)
		block["hourly"] = []any{hourlyBlock}
	}

	return []any{block}, nil
}

func flattenBackupProtectionPolicyFileShareRetentionDaily(daily *protectionpolicies.DailyRetentionSchedule) []any {
	if daily == nil {
		return []any{}
	}

	block := map[string]any{}

	if duration := daily.RetentionDuration; duration != nil {
		if v := duration.Count; v != nil {
			block["count"] = *v
		}
	}

	return []any{block}
}

func flattenBackupProtectionPolicyFileShareRetentionWeekly(weekly *protectionpolicies.WeeklyRetentionSchedule) []any {
	if weekly == nil {
		return []any{}
	}

	block := map[string]any{}

	if duration := weekly.RetentionDuration; duration != nil {
		if v := duration.Count; v != nil {
			block["count"] = *v
		}
	}

	if days := weekly.DaysOfTheWeek; days != nil {
		weekdays := make([]any, 0)
		for _, d := range *days {
			weekdays = append(weekdays, string(d))
		}
		block["weekdays"] = pluginsdk.NewSet(pluginsdk.HashString, weekdays)
	}

	return []any{block}
}

func flattenBackupProtectionPolicyFileShareRetentionMonthly(monthly *protectionpolicies.MonthlyRetentionSchedule) []any {
	if monthly == nil {
		return []any{}
	}

	block := map[string]any{}

	if duration := monthly.RetentionDuration; duration != nil {
		if v := duration.Count; v != nil {
			block["count"] = *v
		}
	}

	if weekly := monthly.RetentionScheduleWeekly; weekly != nil {
		block["weekdays"], block["weeks"] = flattenBackupProtectionPolicyFileShareRetentionWeeklyFormat(weekly)
	}

	if daily := monthly.RetentionScheduleDaily; daily != nil {
		block["days"], block["include_last_days"] = flattenBackupProtectionPolicyFileRetentionDailyFormat(daily)
	}

	return []any{block}
}

func flattenBackupProtectionPolicyFileShareRetentionYearly(yearly *protectionpolicies.YearlyRetentionSchedule) []any {
	if yearly == nil {
		return []any{}
	}

	block := map[string]any{}

	if duration := yearly.RetentionDuration; duration != nil {
		if v := duration.Count; v != nil {
			block["count"] = *v
		}
	}

	if weekly := yearly.RetentionScheduleWeekly; weekly != nil {
		block["weekdays"], block["weeks"] = flattenBackupProtectionPolicyFileShareRetentionWeeklyFormat(weekly)
	}

	if daily := yearly.RetentionScheduleDaily; daily != nil {
		block["days"], block["include_last_days"] = flattenBackupProtectionPolicyFileRetentionDailyFormat(daily)
	}

	if months := yearly.MonthsOfYear; months != nil {
		slice := make([]any, 0)
		for _, d := range *months {
			slice = append(slice, string(d))
		}
		block["months"] = pluginsdk.NewSet(pluginsdk.HashString, slice)
	}

	return []any{block}
}

func flattenBackupProtectionPolicyFileShareRetentionWeeklyFormat(retention *protectionpolicies.WeeklyRetentionFormat) (weekdays, weeks *pluginsdk.Set) {
	if days := retention.DaysOfTheWeek; days != nil {
		slice := make([]any, 0)
		for _, d := range *days {
			slice = append(slice, string(d))
		}
		weekdays = pluginsdk.NewSet(pluginsdk.HashString, slice)
	}

	if days := retention.WeeksOfTheMonth; days != nil {
		slice := make([]any, 0)
		for _, d := range *days {
			slice = append(slice, string(d))
		}
		weeks = pluginsdk.NewSet(pluginsdk.HashString, slice)
	}

	return weekdays, weeks
}

func flattenBackupProtectionPolicyFileRetentionDailyFormat(retention *protectionpolicies.DailyRetentionFormat) (days []any, includeLastDay bool) {
	if dotm := retention.DaysOfTheMonth; dotm != nil {
		for _, d := range *dotm {
			// for the last date, the service will return a record with date == 0.
			if d.Date != nil && (*d.Date != 0) {
				days = append(days, *d.Date)
			}
			if d.IsLast != nil && *d.IsLast {
				includeLastDay = true
			}
		}
	}
	return days, includeLastDay
}
