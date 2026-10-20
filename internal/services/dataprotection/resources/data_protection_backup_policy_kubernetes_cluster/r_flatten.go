// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package data_protection_backup_policy_kubernetes_cluster

import (
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-07-01/basebackuppolicyresources"
)

func flattenBackupPolicyKubernetesClusterBackupRuleArray(input *[]basebackuppolicyresources.BasePolicyRule) []string {
	if input == nil {
		return make([]string, 0)
	}
	for _, item := range *input {
		if backupRule, ok := item.(basebackuppolicyresources.AzureBackupRule); ok {
			if backupRule.Trigger != nil {
				if scheduleBasedTrigger, ok := backupRule.Trigger.(basebackuppolicyresources.ScheduleBasedTriggerContext); ok {
					return scheduleBasedTrigger.Schedule.RepeatingTimeIntervals
				}
			}
		}
	}
	return make([]string, 0)
}

func flattenBackupPolicyKubernetesClusterBackupTimeZone(input *[]basebackuppolicyresources.BasePolicyRule) string {
	if input == nil {
		return ""
	}
	for _, item := range *input {
		if backupRule, ok := item.(basebackuppolicyresources.AzureBackupRule); ok {
			if backupRule.Trigger != nil {
				if scheduleBasedTrigger, ok := backupRule.Trigger.(basebackuppolicyresources.ScheduleBasedTriggerContext); ok {
					return pointer.From(scheduleBasedTrigger.Schedule.TimeZone)
				}
			}
		}
	}
	return ""
}

func flattenBackupPolicyKubernetesClusterRetentionRules(input *[]basebackuppolicyresources.BasePolicyRule) []RetentionRule {
	results := make([]RetentionRule, 0)
	if input == nil {
		return results
	}

	var taggingCriteriaList []basebackuppolicyresources.TaggingCriteria
	for _, item := range *input {
		if backupRule, ok := item.(basebackuppolicyresources.AzureBackupRule); ok {
			if trigger, ok := backupRule.Trigger.(basebackuppolicyresources.ScheduleBasedTriggerContext); ok {
				taggingCriteriaList = trigger.TaggingCriteria
			}
		}
	}

	for _, item := range *input {
		if retentionRule, ok := item.(basebackuppolicyresources.AzureRetentionRule); ok {
			var name string
			var taggingPriority int64
			var taggingCriteria []Criteria
			if retentionRule.IsDefault == nil || !*retentionRule.IsDefault {
				name = retentionRule.Name
				for _, criteria := range taggingCriteriaList {
					if strings.EqualFold(criteria.TagInfo.TagName, name) {
						taggingPriority = criteria.TaggingPriority
						taggingCriteria = flattenBackupPolicyKubernetesClusterBackupCriteriaArray(criteria.Criteria)
						break
					}
				}

				var lifeCycle []LifeCycle
				if v := retentionRule.Lifecycles; len(v) > 0 {
					lifeCycle = flattenBackupPolicyKubernetesClusterBackupLifeCycleArray(v)
				}
				results = append(results, RetentionRule{
					Name:      name,
					Priority:  taggingPriority,
					Criteria:  taggingCriteria,
					LifeCycle: lifeCycle,
				})
			}
		}
	}
	return results
}

func flattenBackupPolicyKubernetesClusterDefaultRetentionRule(input *[]basebackuppolicyresources.BasePolicyRule) []DefaultRetentionRule {
	results := make([]DefaultRetentionRule, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		if retentionRule, ok := item.(basebackuppolicyresources.AzureRetentionRule); ok {
			if pointer.From(retentionRule.IsDefault) {
				var lifeCycle []LifeCycle
				if v := retentionRule.Lifecycles; len(v) > 0 {
					lifeCycle = flattenBackupPolicyKubernetesClusterBackupLifeCycleArray(v)
				}

				results = append(results, DefaultRetentionRule{
					LifeCycle: lifeCycle,
				})
			}
		}
	}
	return results
}

func flattenBackupPolicyKubernetesClusterBackupCriteriaArray(input *[]basebackuppolicyresources.BackupCriteria) []Criteria {
	results := make([]Criteria, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		if criteria, ok := item.(basebackuppolicyresources.ScheduleBasedBackupCriteria); ok {
			var absoluteCriteria string
			if criteria.AbsoluteCriteria != nil && len(*criteria.AbsoluteCriteria) > 0 {
				absoluteCriteria = string((*criteria.AbsoluteCriteria)[0])
			}
			var daysOfWeek []string
			if criteria.DaysOfTheWeek != nil {
				daysOfWeek = make([]string, 0)
				for _, item := range *criteria.DaysOfTheWeek {
					daysOfWeek = append(daysOfWeek, string(item))
				}
			}
			var monthsOfYear []string
			if criteria.MonthsOfYear != nil {
				monthsOfYear = make([]string, 0)
				for _, item := range *criteria.MonthsOfYear {
					monthsOfYear = append(monthsOfYear, string(item))
				}
			}
			var weeksOfMonth []string
			if criteria.WeeksOfTheMonth != nil {
				weeksOfMonth = make([]string, 0)
				for _, item := range *criteria.WeeksOfTheMonth {
					weeksOfMonth = append(weeksOfMonth, string(item))
				}
			}
			var scheduleTimes []string
			if criteria.ScheduleTimes != nil {
				scheduleTimes = make([]string, 0)
				scheduleTimes = append(scheduleTimes, *criteria.ScheduleTimes...)
			}

			results = append(results, Criteria{
				AbsoluteCriteria:     absoluteCriteria,
				DaysOfWeek:           daysOfWeek,
				MonthsOfYear:         monthsOfYear,
				WeeksOfMonth:         weeksOfMonth,
				ScheduledBackupTimes: scheduleTimes,
			})
		}
	}
	return results
}

func flattenBackupPolicyKubernetesClusterBackupLifeCycleArray(input []basebackuppolicyresources.SourceLifeCycle) []LifeCycle {
	results := make([]LifeCycle, 0)
	if input == nil {
		return results
	}

	for _, item := range input {
		var duration string
		var dataStoreType string
		if deleteOption, ok := item.DeleteAfter.(basebackuppolicyresources.AbsoluteDeleteOption); ok {
			duration = deleteOption.Duration
		}
		dataStoreType = string(item.SourceDataStore.DataStoreType)

		results = append(results, LifeCycle{
			Duration:      duration,
			DataStoreType: dataStoreType,
		})
	}
	return results
}
