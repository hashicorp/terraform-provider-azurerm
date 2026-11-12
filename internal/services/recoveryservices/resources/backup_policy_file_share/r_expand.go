// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package backup_policy_file_share

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicesbackup/2024-10-01/protectionpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func expandBackupProtectionPolicyFileShareSchedule(d *pluginsdk.ResourceData, times []string) *protectionpolicies.SimpleSchedulePolicy {
	if bb, ok := d.Get("backup").([]any); ok && len(bb) > 0 {
		block := bb[0].(map[string]any)

		schedule := protectionpolicies.SimpleSchedulePolicy{ // LongTermSchedulePolicy has no properties
			ScheduleRunTimes: &times,
		}

		if v, ok := block["frequency"].(string); ok {
			schedule.ScheduleRunFrequency = pointer.ToEnum[protectionpolicies.ScheduleRunType](v)
		}

		if v, ok := block["hourly"].([]any); ok && len(v) > 0 {
			hourlyBlock := v[0].(map[string]any)

			if schedule.ScheduleRunFrequency != nil && *schedule.ScheduleRunFrequency == protectionpolicies.ScheduleRunTypeHourly {
				schedule.HourlySchedule = &protectionpolicies.HourlySchedule{}
				if v, ok := hourlyBlock["interval"].(int); ok {
					schedule.HourlySchedule.Interval = pointer.To(int64(v))
				}
				if t, ok := hourlyBlock["start_time"].(string); ok {
					schedule.HourlySchedule.ScheduleWindowStartTime = pointer.To(fmt.Sprintf("2018-07-30T%s:00.000Z", t))
				}
				if v, ok := hourlyBlock["window_duration"].(int); ok {
					schedule.HourlySchedule.ScheduleWindowDuration = pointer.To(int64(v))
				}
			}
		}

		return &schedule
	}

	return nil
}

func expandBackupProtectionPolicyFileShareRetentionDaily(d *pluginsdk.ResourceData, times []string) *protectionpolicies.DailyRetentionSchedule {
	if rb, ok := d.Get("retention_daily").([]any); ok && len(rb) > 0 {
		block := rb[0].(map[string]any)

		return &protectionpolicies.DailyRetentionSchedule{
			RetentionTimes: &times,
			RetentionDuration: &protectionpolicies.RetentionDuration{
				Count:        pointer.To(int64(block["count"].(int))),
				DurationType: pointer.To(protectionpolicies.RetentionDurationTypeDays),
			},
		}
	}

	return nil
}

func expandBackupProtectionPolicyFileShareRetentionWeekly(d *pluginsdk.ResourceData, times []string) *protectionpolicies.WeeklyRetentionSchedule {
	if rb, ok := d.Get("retention_weekly").([]any); ok && len(rb) > 0 {
		block := rb[0].(map[string]any)

		retention := protectionpolicies.WeeklyRetentionSchedule{
			RetentionTimes: &times,
			RetentionDuration: &protectionpolicies.RetentionDuration{
				Count:        pointer.To(int64(block["count"].(int))),
				DurationType: pointer.To(protectionpolicies.RetentionDurationTypeWeeks),
			},
		}

		if v, ok := block["weekdays"].(*pluginsdk.Set); ok {
			days := make([]protectionpolicies.DayOfWeek, 0)
			for _, day := range v.List() {
				days = append(days, protectionpolicies.DayOfWeek(day.(string)))
			}
			retention.DaysOfTheWeek = &days
		}

		return &retention
	}

	return nil
}

func expandBackupProtectionPolicyFileShareRetentionMonthly(d *pluginsdk.ResourceData, times []string) *protectionpolicies.MonthlyRetentionSchedule {
	if rb, ok := d.Get("retention_monthly").([]any); ok && len(rb) > 0 {
		block := rb[0].(map[string]any)

		scheduleFormat := protectionpolicies.RetentionScheduleFormatWeekly
		var weekly *protectionpolicies.WeeklyRetentionFormat = nil
		var daily *protectionpolicies.DailyRetentionFormat = nil
		if v, ok := block["days"]; ok && v.(*pluginsdk.Set).Len() > 0 || block["include_last_days"].(bool) {
			scheduleFormat = protectionpolicies.RetentionScheduleFormatDaily
			daily = expandBackupProtectionPolicyFileShareRetentionDailyFormat(block)
		} else {
			weekly = expandBackupProtectionPolicyFileShareRetentionWeeklyFormat(block)
		}

		retention := protectionpolicies.MonthlyRetentionSchedule{
			RetentionScheduleFormatType: &scheduleFormat,
			RetentionScheduleDaily:      daily,
			RetentionScheduleWeekly:     weekly,
			RetentionTimes:              &times,
			RetentionDuration: &protectionpolicies.RetentionDuration{
				Count:        pointer.To(int64(block["count"].(int))),
				DurationType: pointer.To(protectionpolicies.RetentionDurationTypeMonths),
			},
		}

		return &retention
	}

	return nil
}

func expandBackupProtectionPolicyFileShareRetentionYearly(d *pluginsdk.ResourceData, times []string) *protectionpolicies.YearlyRetentionSchedule {
	if rb, ok := d.Get("retention_yearly").([]any); ok && len(rb) > 0 {
		block := rb[0].(map[string]any)

		scheduleFormat := protectionpolicies.RetentionScheduleFormatWeekly
		var weekly *protectionpolicies.WeeklyRetentionFormat = nil
		var daily *protectionpolicies.DailyRetentionFormat = nil
		if v, ok := block["days"]; ok && v.(*pluginsdk.Set).Len() > 0 || block["include_last_days"].(bool) {
			scheduleFormat = protectionpolicies.RetentionScheduleFormatDaily
			daily = expandBackupProtectionPolicyFileShareRetentionDailyFormat(block)
		} else {
			weekly = expandBackupProtectionPolicyFileShareRetentionWeeklyFormat(block)
		}

		retention := protectionpolicies.YearlyRetentionSchedule{
			RetentionScheduleFormatType: &scheduleFormat,
			RetentionScheduleDaily:      daily,
			RetentionScheduleWeekly:     weekly,
			RetentionTimes:              &times,
			RetentionDuration: &protectionpolicies.RetentionDuration{
				Count:        pointer.To(int64(block["count"].(int))),
				DurationType: pointer.To(protectionpolicies.RetentionDurationTypeYears),
			},
		}

		if v, ok := block["months"].(*pluginsdk.Set); ok {
			months := make([]protectionpolicies.MonthOfYear, 0)
			for _, month := range v.List() {
				months = append(months, protectionpolicies.MonthOfYear(month.(string)))
			}
			retention.MonthsOfYear = &months
		}

		return &retention
	}

	return nil
}

func expandBackupProtectionPolicyFileShareRetentionWeeklyFormat(block map[string]any) *protectionpolicies.WeeklyRetentionFormat {
	weekly := protectionpolicies.WeeklyRetentionFormat{}

	if v, ok := block["weekdays"].(*pluginsdk.Set); ok {
		days := make([]protectionpolicies.DayOfWeek, 0)
		for _, day := range v.List() {
			days = append(days, protectionpolicies.DayOfWeek(day.(string)))
		}
		weekly.DaysOfTheWeek = &days
	}

	if v, ok := block["weeks"].(*pluginsdk.Set); ok {
		weeks := make([]protectionpolicies.WeekOfMonth, 0)
		for _, week := range v.List() {
			weeks = append(weeks, protectionpolicies.WeekOfMonth(week.(string)))
		}
		weekly.WeeksOfTheMonth = &weeks
	}

	return &weekly
}

func expandBackupProtectionPolicyFileShareRetentionDailyFormat(block map[string]any) *protectionpolicies.DailyRetentionFormat {
	days := make([]protectionpolicies.Day, 0)

	if block["include_last_days"].(bool) {
		days = append(days, protectionpolicies.Day{
			Date:   pointer.To(int64(0)),
			IsLast: pointer.To(true),
		})
	}

	if v, ok := block["days"].(*pluginsdk.Set); ok {
		for _, day := range v.List() {
			days = append(days, protectionpolicies.Day{
				Date:   pointer.To(int64(day.(int))),
				IsLast: pointer.To(false),
			})
		}
	}

	daily := protectionpolicies.DailyRetentionFormat{
		DaysOfTheMonth: &days,
	}

	return &daily
}
