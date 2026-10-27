// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package backup_policy_vm

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicesbackup/2024-10-01/protectionpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func resourceBackupProtectionPolicyVMFlatten(d *pluginsdk.ResourceData, id *protectionpolicies.BackupPolicyId, model *protectionpolicies.ProtectionPolicyResource) error {
	d.Set("name", id.BackupPolicyName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("recovery_vault_name", id.VaultName)

	if model != nil {
		if properties, ok := model.Properties.(protectionpolicies.AzureIaaSVMProtectionPolicy); ok {
			d.Set("timezone", properties.TimeZone)
			d.Set("instant_restore_retention_days", properties.InstantRpRetentionRangeInDays)
			d.Set("tiering_policy", flattenBackupProtectionPolicyVMTieringPolicy(properties.TieringPolicy))

			if schedule, ok := properties.SchedulePolicy.(protectionpolicies.SimpleSchedulePolicy); ok {
				if err := d.Set("backup", flattenBackupProtectionPolicyVMSchedule(schedule)); err != nil {
					return fmt.Errorf("setting `backup`: %+v", err)
				}
			}

			if schedule, ok := properties.SchedulePolicy.(protectionpolicies.SimpleSchedulePolicyV2); ok {
				if err := d.Set("backup", flattenBackupProtectionPolicyVMScheduleV2(schedule)); err != nil {
					return fmt.Errorf("setting `backup`: %+v", err)
				}
			}

			policyType := string(protectionpolicies.IAASVMPolicyTypeVOne)
			if pointer.From(properties.PolicyType) != "" {
				policyType = pointer.FromEnum(properties.PolicyType)
			}
			d.Set("policy_type", policyType)

			d.Set("consistency_type", pointer.FromEnum(properties.SnapshotConsistencyType))

			if retention, ok := properties.RetentionPolicy.(protectionpolicies.LongTermRetentionPolicy); ok {
				if s := retention.DailySchedule; s != nil {
					if err := d.Set("retention_daily", flattenBackupProtectionPolicyVMRetentionDaily(s)); err != nil {
						return fmt.Errorf("setting `retention_daily`: %+v", err)
					}
				} else {
					d.Set("retention_daily", nil)
				}

				if s := retention.WeeklySchedule; s != nil {
					if err := d.Set("retention_weekly", flattenBackupProtectionPolicyVMRetentionWeekly(s)); err != nil {
						return fmt.Errorf("setting `retention_weekly`: %+v", err)
					}
				} else {
					d.Set("retention_weekly", nil)
				}

				if s := retention.MonthlySchedule; s != nil {
					if err := d.Set("retention_monthly", flattenBackupProtectionPolicyVMRetentionMonthly(s)); err != nil {
						return fmt.Errorf("setting `retention_monthly`: %+v", err)
					}
				} else {
					d.Set("retention_monthly", nil)
				}

				if s := retention.YearlySchedule; s != nil {
					if err := d.Set("retention_yearly", flattenBackupProtectionPolicyVMRetentionYearly(s)); err != nil {
						return fmt.Errorf("setting `retention_yearly`: %+v", err)
					}
				} else {
					d.Set("retention_yearly", nil)
				}
			}

			if instantRPDetail := properties.InstantRPDetails; instantRPDetail != nil {
				d.Set("instant_restore_resource_group", flattenBackupProtectionPolicyVMResourceGroup(*instantRPDetail))
			}
		}
	}

	return pluginsdk.SetResourceIdentityData(d, id)
}

func flattenBackupProtectionPolicyVMResourceGroup(rpDetail protectionpolicies.InstantRPAdditionalDetails) []any {
	if rpDetail.AzureBackupRGNamePrefix == nil {
		return []any{}
	}

	block := map[string]any{}

	block["prefix"] = pointer.From(rpDetail.AzureBackupRGNamePrefix)

	block["suffix"] = pointer.From(rpDetail.AzureBackupRGNameSuffix)

	return []any{block}
}

func flattenBackupProtectionPolicyVMSchedule(schedule protectionpolicies.SimpleSchedulePolicy) []any {
	block := map[string]any{}

	block["frequency"] = pointer.FromEnum(schedule.ScheduleRunFrequency)

	if times := schedule.ScheduleRunTimes; times != nil && len(*times) > 0 {
		policyTime, _ := time.Parse(time.RFC3339, (*times)[0])
		block["time"] = policyTime.Format("15:04")
	}

	if days := schedule.ScheduleRunDays; days != nil {
		weekdays := make([]any, 0)
		for _, d := range *days {
			weekdays = append(weekdays, string(d))
		}
		block["weekdays"] = pluginsdk.NewSet(pluginsdk.HashString, weekdays)
	}

	return []any{block}
}

func flattenBackupProtectionPolicyVMScheduleV2(schedule protectionpolicies.SimpleSchedulePolicyV2) []any {
	block := map[string]any{}

	frequency := pointer.From(schedule.ScheduleRunFrequency)
	block["frequency"] = string(frequency)

	switch frequency {
	case protectionpolicies.ScheduleRunTypeHourly:
		schedule := schedule.HourlySchedule
		if schedule.Interval != nil {
			block["hour_interval"] = *schedule.Interval
		}

		if schedule.ScheduleWindowDuration != nil {
			block["hour_duration"] = *schedule.ScheduleWindowDuration
		}

		if schedule.ScheduleWindowStartTime != nil {
			policyTime, _ := time.Parse(time.RFC3339, pointer.From(schedule.ScheduleWindowStartTime))
			block["time"] = policyTime.Format("15:04")
		}
	case protectionpolicies.ScheduleRunTypeDaily:
		schedule := schedule.DailySchedule
		if times := schedule.ScheduleRunTimes; times != nil && len(*times) > 0 {
			policyTime, _ := time.Parse(time.RFC3339, (*times)[0])
			block["time"] = policyTime.Format("15:04")
		}
	case protectionpolicies.ScheduleRunTypeWeekly:
		schedule := schedule.WeeklySchedule
		if days := schedule.ScheduleRunDays; days != nil {
			weekdays := make([]any, 0)
			for _, d := range *days {
				weekdays = append(weekdays, string(d))
			}
			block["weekdays"] = pluginsdk.NewSet(pluginsdk.HashString, weekdays)
		}

		if times := schedule.ScheduleRunTimes; times != nil && len(*times) > 0 {
			policyTime, _ := time.Parse(time.RFC3339, (*times)[0])
			block["time"] = policyTime.Format("15:04")
		}
	default:
	}

	return []any{block}
}

func flattenBackupProtectionPolicyVMRetentionDaily(daily *protectionpolicies.DailyRetentionSchedule) []any {
	block := map[string]any{}

	if duration := daily.RetentionDuration; duration != nil {
		if v := duration.Count; v != nil {
			block["count"] = *v
		}
	}

	return []any{block}
}

func flattenBackupProtectionPolicyVMRetentionWeekly(weekly *protectionpolicies.WeeklyRetentionSchedule) []any {
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

func flattenBackupProtectionPolicyVMRetentionMonthly(monthly *protectionpolicies.MonthlyRetentionSchedule) []any {
	block := map[string]any{}

	if duration := monthly.RetentionDuration; duration != nil {
		if v := duration.Count; v != nil {
			block["count"] = *v
		}
	}

	if weekly := monthly.RetentionScheduleWeekly; weekly != nil {
		block["weekdays"], block["weeks"] = flattenBackupProtectionPolicyVMRetentionWeeklyFormat(weekly)
	}

	if daily := monthly.RetentionScheduleDaily; daily != nil {
		block["days"], block["include_last_days"] = flattenBackupProtectionPolicyVMRetentionDailyFormat(daily)
	}

	return []any{block}
}

func flattenBackupProtectionPolicyVMRetentionYearly(yearly *protectionpolicies.YearlyRetentionSchedule) []any {
	block := map[string]any{}

	if duration := yearly.RetentionDuration; duration != nil {
		if v := duration.Count; v != nil {
			block["count"] = *v
		}
	}

	if weekly := yearly.RetentionScheduleWeekly; weekly != nil {
		block["weekdays"], block["weeks"] = flattenBackupProtectionPolicyVMRetentionWeeklyFormat(weekly)
	}

	if months := yearly.MonthsOfYear; months != nil {
		slice := make([]any, 0)
		for _, d := range *months {
			slice = append(slice, string(d))
		}
		block["months"] = pluginsdk.NewSet(pluginsdk.HashString, slice)
	}

	if daily := yearly.RetentionScheduleDaily; daily != nil {
		block["days"], block["include_last_days"] = flattenBackupProtectionPolicyVMRetentionDailyFormat(daily)
	}

	return []any{block}
}

func flattenBackupProtectionPolicyVMRetentionWeeklyFormat(retention *protectionpolicies.WeeklyRetentionFormat) (weekdays, weeks *pluginsdk.Set) {
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

func flattenBackupProtectionPolicyVMRetentionDailyFormat(retention *protectionpolicies.DailyRetentionFormat) (days []any, includeLastDay bool) {
	if dotm := retention.DaysOfTheMonth; dotm != nil {
		for _, d := range *dotm {
			if d.Date != nil && *d.Date != 0 {
				days = append(days, *d.Date)
			}
			if d.IsLast != nil && *d.IsLast {
				includeLastDay = true
			}
		}
	}

	return days, includeLastDay
}

func flattenBackupProtectionPolicyVMTieringPolicy(input *map[string]protectionpolicies.TieringPolicy) []any {
	results := make([]any, 0)
	if input == nil {
		return results
	}

	for k, v := range *input {
		if k == "ArchivedRP" {
			if pointer.From(v.TieringMode) == protectionpolicies.TieringModeDoNotTier && pointer.From(v.DurationType) == protectionpolicies.RetentionDurationTypeInvalid && pointer.From(v.Duration) == 0 {
				return results
			}

			results = append(results, map[string]any{
				"archived_restore_point": flattenBackupProtectionPolicyVMArchivedRP(v),
			})
		}
	}

	return results
}

func flattenBackupProtectionPolicyVMArchivedRP(input protectionpolicies.TieringPolicy) []any {
	results := make([]any, 0)

	result := map[string]any{
		"mode":     pointer.FromEnum(input.TieringMode),
		"duration": int(pointer.From(input.Duration)),
	}

	durationType := ""
	if v := input.DurationType; v != nil && pointer.From(v) != protectionpolicies.RetentionDurationTypeInvalid {
		durationType = pointer.FromEnum(v)
	}
	result["duration_type"] = durationType

	results = append(results, result)

	return results
}
