// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package backup_policy_vm_workload

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicesbackup/2024-10-01/protectionpolicies"
)

func flattenBackupProtectionPolicyVMWorkloadSettings(input *protectionpolicies.Settings) []Settings {
	if input == nil {
		return make([]Settings, 0)
	}

	result := make([]Settings, 0)

	result = append(result, Settings{
		CompressionEnabled: *input.IsCompression,
		TimeZone:           *input.TimeZone,
	})

	return result
}

func flattenBackupProtectionPolicyVMWorkloadProtectionPolicies(input *[]protectionpolicies.SubProtectionPolicy) []ProtectionPolicy {
	results := make([]ProtectionPolicy, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		result := ProtectionPolicy{
			PolicyType: pointer.FromEnum(item.PolicyType),
			Backup:     flattenBackupProtectionPolicyVMWorkloadSchedulePolicy(item.SchedulePolicy, pointer.From(item.PolicyType)),
		}

		if retentionPolicy := item.RetentionPolicy; retentionPolicy != nil {
			if longTermRetentionPolicy, ok := retentionPolicy.(protectionpolicies.LongTermRetentionPolicy); ok {
				result.RetentionDaily = flattenBackupProtectionPolicyVMWorkloadRetentionDaily(longTermRetentionPolicy.DailySchedule)
				result.RetentionWeekly = flattenBackupProtectionPolicyVMWorkloadRetentionWeekly(longTermRetentionPolicy.WeeklySchedule)
				result.RetentionMonthly = flattenBackupProtectionPolicyVMWorkloadRetentionMonthly(longTermRetentionPolicy.MonthlySchedule)
				result.RetentionYearly = flattenBackupProtectionPolicyVMWorkloadRetentionYearly(longTermRetentionPolicy.YearlySchedule)
			} else {
				simpleRetentionPolicy, _ := retentionPolicy.(protectionpolicies.SimpleRetentionPolicy)
				result.SimpleRetention = flattenBackupProtectionPolicyVMWorkloadSimpleRetention(simpleRetentionPolicy.RetentionDuration)
			}
		}

		results = append(results, result)
	}

	return results
}

func flattenBackupProtectionPolicyVMWorkloadSchedulePolicy(input protectionpolicies.SchedulePolicy, policyType protectionpolicies.PolicyType) []Backup {
	if input == nil {
		return []Backup{}
	}

	backupBlock := Backup{}

	if policyType == protectionpolicies.PolicyTypeLog {
		logSchedulePolicy, _ := input.(protectionpolicies.LogSchedulePolicy)

		if v := logSchedulePolicy.ScheduleFrequencyInMins; v != nil {
			backupBlock.FrequencyInMinutes = *v
		}
	} else {
		simpleSchedulePolicy, _ := input.(protectionpolicies.SimpleSchedulePolicy)

		backupBlock.Frequency = pointer.FromEnum(simpleSchedulePolicy.ScheduleRunFrequency)

		if times := simpleSchedulePolicy.ScheduleRunTimes; times != nil && len(*times) > 0 {
			policyTime, _ := time.Parse(time.RFC3339, (*times)[0])
			backupBlock.Time = policyTime.Format("15:04")
		}

		if days := simpleSchedulePolicy.ScheduleRunDays; days != nil {
			weekdays := make([]string, 0)
			for _, d := range *days {
				weekdays = append(weekdays, string(d))
			}
			backupBlock.Weekdays = weekdays
		}
	}

	return []Backup{backupBlock}
}

func flattenBackupProtectionPolicyVMWorkloadRetentionDaily(input *protectionpolicies.DailyRetentionSchedule) []RetentionDaily {
	if input == nil {
		return []RetentionDaily{}
	}

	retentionDailyBlock := RetentionDaily{}

	if duration := input.RetentionDuration; duration != nil {
		if v := duration.Count; v != nil {
			retentionDailyBlock.Count = *v
		}
	}

	return []RetentionDaily{retentionDailyBlock}
}

func flattenBackupProtectionPolicyVMWorkloadRetentionWeekly(input *protectionpolicies.WeeklyRetentionSchedule) []RetentionWeekly {
	if input == nil {
		return []RetentionWeekly{}
	}

	retentionWeeklyBlock := RetentionWeekly{}

	if duration := input.RetentionDuration; duration != nil {
		if v := duration.Count; v != nil {
			retentionWeeklyBlock.Count = *v
		}
	}

	if days := input.DaysOfTheWeek; days != nil {
		weekdays := make([]string, 0)
		for _, d := range *days {
			weekdays = append(weekdays, string(d))
		}
		retentionWeeklyBlock.Weekdays = weekdays
	}

	return []RetentionWeekly{retentionWeeklyBlock}
}

func flattenBackupProtectionPolicyVMWorkloadRetentionMonthly(input *protectionpolicies.MonthlyRetentionSchedule) []RetentionMonthly {
	if input == nil {
		return []RetentionMonthly{}
	}

	retentionMonthlyBlock := RetentionMonthly{}

	if duration := input.RetentionDuration; duration != nil {
		if v := duration.Count; v != nil {
			retentionMonthlyBlock.Count = *v
		}
	}

	if formatType := pointer.From(input.RetentionScheduleFormatType); formatType != "" {
		retentionMonthlyBlock.FormatType = string(formatType)
	}

	if weekly := input.RetentionScheduleWeekly; weekly != nil {
		retentionMonthlyBlock.Weekdays, retentionMonthlyBlock.Weeks = flattenBackupProtectionPolicyVMWorkloadRetentionWeeklyFormat(weekly)
	}

	if daily := input.RetentionScheduleDaily; daily != nil {
		retentionMonthlyBlock.Monthdays = flattenBackupProtectionPolicyVMWorkloadRetentionDailyFormat(daily)
	}

	return []RetentionMonthly{retentionMonthlyBlock}
}

func flattenBackupProtectionPolicyVMWorkloadRetentionYearly(input *protectionpolicies.YearlyRetentionSchedule) []RetentionYearly {
	if input == nil {
		return []RetentionYearly{}
	}

	retentionYearlyBlock := RetentionYearly{}

	if duration := input.RetentionDuration; duration != nil {
		if v := duration.Count; v != nil {
			retentionYearlyBlock.Count = *v
		}
	}

	if formatType := pointer.From(input.RetentionScheduleFormatType); formatType != "" {
		retentionYearlyBlock.FormatType = string(formatType)
	}

	if weekly := input.RetentionScheduleWeekly; weekly != nil {
		retentionYearlyBlock.Weekdays, retentionYearlyBlock.Weeks = flattenBackupProtectionPolicyVMWorkloadRetentionWeeklyFormat(weekly)
	}

	if v := input.MonthsOfYear; v != nil {
		months := make([]string, 0)
		for _, d := range *v {
			months = append(months, string(d))
		}
		retentionYearlyBlock.Months = months
	}

	if daily := input.RetentionScheduleDaily; daily != nil {
		retentionYearlyBlock.Monthdays = flattenBackupProtectionPolicyVMWorkloadRetentionDailyFormat(daily)
	}

	return []RetentionYearly{retentionYearlyBlock}
}

func flattenBackupProtectionPolicyVMWorkloadSimpleRetention(input *protectionpolicies.RetentionDuration) []SimpleRetention {
	if input == nil {
		return []SimpleRetention{}
	}

	simpleRetentionBlock := SimpleRetention{}

	if v := input.Count; v != nil {
		simpleRetentionBlock.Count = *v
	}

	return []SimpleRetention{simpleRetentionBlock}
}

func flattenBackupProtectionPolicyVMWorkloadRetentionWeeklyFormat(input *protectionpolicies.WeeklyRetentionFormat) (weekdays, weeks []string) {
	if v := input.DaysOfTheWeek; v != nil {
		days := make([]string, 0)
		for _, d := range *v {
			days = append(days, string(d))
		}
		weekdays = days
	}

	if v := input.WeeksOfTheMonth; v != nil {
		days := make([]string, 0)
		for _, d := range *v {
			days = append(days, string(d))
		}
		weeks = days
	}

	return weekdays, weeks
}

func flattenBackupProtectionPolicyVMWorkloadRetentionDailyFormat(input *protectionpolicies.DailyRetentionFormat) []int64 {
	result := make([]int64, 0)

	if days := input.DaysOfTheMonth; days != nil {
		for _, d := range *days {
			result = append(result, pointer.From(d.Date))
		}
	}

	return result
}
