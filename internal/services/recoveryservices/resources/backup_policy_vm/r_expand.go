// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package backup_policy_vm

import (
	"errors"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicesbackup/2024-10-01/protectionpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func expandBackupProtectionPolicyVMSchedule(d *pluginsdk.ResourceData, times []string) (protectionpolicies.SchedulePolicy, error) {
	if bb, ok := d.Get("backup").([]any); ok && len(bb) > 0 {
		block := bb[0].(map[string]any)

		policyType := d.Get("policy_type").(string)
		if policyType == string(protectionpolicies.IAASVMPolicyTypeVOne) {
			schedule := protectionpolicies.SimpleSchedulePolicy{ // LongTermSchedulePolicy has no properties
				ScheduleRunTimes: &times,
			}

			if v, ok := block["frequency"].(string); ok {
				schedule.ScheduleRunFrequency = pointer.ToEnum[protectionpolicies.ScheduleRunType](v)
			}

			if v, ok := block["weekdays"].(*pluginsdk.Set); ok {
				days := make([]protectionpolicies.DayOfWeek, 0)
				for _, day := range v.List() {
					days = append(days, protectionpolicies.DayOfWeek(day.(string)))
				}
				schedule.ScheduleRunDays = &days
			}

			return schedule, nil
		} else {
			frequency := block["frequency"].(string)
			schedule := protectionpolicies.SimpleSchedulePolicyV2{
				ScheduleRunFrequency: pointer.ToEnum[protectionpolicies.ScheduleRunType](frequency),
			}

			switch frequency {
			case string(protectionpolicies.ScheduleRunTypeHourly):
				interval, ok := block["hour_interval"].(int)
				if !ok {
					return nil, fmt.Errorf("`hour_interval` must be specified when `backup.0.frequency` is `Hourly`")
				}

				duration, ok := block["hour_duration"].(int)
				if !ok {
					return nil, fmt.Errorf("`hour_duration` must be specified when `backup.0.frequency` is `Hourly`")
				}

				if interval == 0 && duration == 0 {
					return nil, fmt.Errorf("`hour_interval` and `hour_duration` must be specified when `backup.0.frequency` is `Hourly`")
				}
				if interval == 0 {
					return nil, fmt.Errorf("`hour_interval` must be specified when `backup.0.frequency` is `Hourly`")
				}
				if duration == 0 {
					return nil, fmt.Errorf("`hour_duration` must be specified when `backup.0.frequency` is `Hourly`")
				}

				if duration%interval != 0 {
					return nil, fmt.Errorf("`hour_duration` must be multiplier of `hour_interval`")
				}

				schedule.HourlySchedule = &protectionpolicies.HourlySchedule{
					Interval:                pointer.To(int64(interval)),
					ScheduleWindowStartTime: &times[0],
					ScheduleWindowDuration:  pointer.To(int64(duration)),
				}
			case string(protectionpolicies.ScheduleRunTypeDaily):
				schedule.DailySchedule = &protectionpolicies.DailySchedule{
					ScheduleRunTimes: &times,
				}
			case string(protectionpolicies.ScheduleRunTypeWeekly):
				weekDays, ok := block["weekdays"].(*pluginsdk.Set)
				if !ok {
					return nil, fmt.Errorf("`weekdays` must be specified when `backup.0.frequency` is `Weekly`")
				}

				days := make([]protectionpolicies.DayOfWeek, 0)
				for _, day := range weekDays.List() {
					days = append(days, protectionpolicies.DayOfWeek(day.(string)))
				}

				schedule.WeeklySchedule = &protectionpolicies.WeeklySchedule{
					ScheduleRunDays:  &days,
					ScheduleRunTimes: &times,
				}
			default:
				return nil, errors.New("unrecognized value for backup.0.frequency")
			}

			return schedule, nil
		}
	}

	return nil, nil
}

func expandBackupProtectionPolicyVMResourceGroup(d *pluginsdk.ResourceData) *protectionpolicies.InstantRPAdditionalDetails {
	if v, ok := d.Get("instant_restore_resource_group").([]any); ok && len(v) > 0 {
		rgRaw := v[0].(map[string]any)
		output := &protectionpolicies.InstantRPAdditionalDetails{
			AzureBackupRGNamePrefix: pointer.To(rgRaw["prefix"].(string)),
		}
		if suffix, ok := rgRaw["suffix"]; ok && suffix != "" {
			output.AzureBackupRGNameSuffix = pointer.To(suffix.(string))
		}
		return output
	}

	return nil
}

func expandBackupProtectionPolicyVMRetentionDaily(d *pluginsdk.ResourceData, times []string) *protectionpolicies.DailyRetentionSchedule {
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

func expandBackupProtectionPolicyVMRetentionWeekly(d *pluginsdk.ResourceData, times []string) *protectionpolicies.WeeklyRetentionSchedule {
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

func expandBackupProtectionPolicyVMRetentionMonthly(d *pluginsdk.ResourceData, times []string) *protectionpolicies.MonthlyRetentionSchedule {
	if rb, ok := d.Get("retention_monthly").([]any); ok && len(rb) > 0 {
		block := rb[0].(map[string]any)

		scheduleFormat := protectionpolicies.RetentionScheduleFormatWeekly
		var weekly *protectionpolicies.WeeklyRetentionFormat = nil
		var daily *protectionpolicies.DailyRetentionFormat = nil
		if v, ok := block["days"]; ok && v.(*pluginsdk.Set).Len() > 0 || block["include_last_days"].(bool) {
			scheduleFormat = protectionpolicies.RetentionScheduleFormatDaily
			daily = expandBackupProtectionPolicyVMRetentionDailyFormat(block)
		} else {
			weekly = expandBackupProtectionPolicyVMRetentionWeeklyFormat(block)
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

func expandBackupProtectionPolicyVMRetentionYearly(d *pluginsdk.ResourceData, times []string) *protectionpolicies.YearlyRetentionSchedule {
	if rb, ok := d.Get("retention_yearly").([]any); ok && len(rb) > 0 {
		block := rb[0].(map[string]any)

		scheduleFormat := protectionpolicies.RetentionScheduleFormatWeekly
		var weekly *protectionpolicies.WeeklyRetentionFormat = nil
		var daily *protectionpolicies.DailyRetentionFormat = nil
		if v, ok := block["days"]; ok && v.(*pluginsdk.Set).Len() > 0 || block["include_last_days"].(bool) {
			scheduleFormat = protectionpolicies.RetentionScheduleFormatDaily
			daily = expandBackupProtectionPolicyVMRetentionDailyFormat(block)
		} else {
			weekly = expandBackupProtectionPolicyVMRetentionWeeklyFormat(block)
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

func expandBackupProtectionPolicyVMRetentionWeeklyFormat(block map[string]any) *protectionpolicies.WeeklyRetentionFormat {
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

func expandBackupProtectionPolicyVMRetentionDailyFormat(block map[string]any) *protectionpolicies.DailyRetentionFormat {
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

func expandBackupProtectionPolicyVMTieringPolicy(input []any) *map[string]protectionpolicies.TieringPolicy {
	result := make(map[string]protectionpolicies.TieringPolicy)
	if len(input) == 0 {
		result["ArchivedRP"] = protectionpolicies.TieringPolicy{
			TieringMode:  pointer.To(protectionpolicies.TieringModeDoNotTier),
			DurationType: pointer.To(protectionpolicies.RetentionDurationTypeInvalid),
			Duration:     pointer.To(int64(0)),
		}

		return &result
	}

	tieringPolicy := input[0].(map[string]any)
	archivedRP := tieringPolicy["archived_restore_point"].([]any)
	result["ArchivedRP"] = expandBackupProtectionPolicyVMArchivedRP(archivedRP)

	return &result
}

func expandBackupProtectionPolicyVMArchivedRP(input []any) protectionpolicies.TieringPolicy {
	if len(input) == 0 {
		return protectionpolicies.TieringPolicy{}
	}

	archivedRP := input[0].(map[string]any)

	result := protectionpolicies.TieringPolicy{
		TieringMode: pointer.ToEnum[protectionpolicies.TieringMode](archivedRP["mode"].(string)),
	}

	if v := archivedRP["duration_type"].(string); v != "" {
		result.DurationType = pointer.ToEnum[protectionpolicies.RetentionDurationType](v)
	}

	if v := archivedRP["duration"].(int); v != 0 {
		result.Duration = pointer.To(int64(v))
	}

	return result
}
