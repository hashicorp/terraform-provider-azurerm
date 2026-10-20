// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package data_protection_backup_policy_blob_storage

import (
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-07-01/basebackuppolicyresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func flattenBackupPolicyBlobStorageDefaultRetentionRuleDuration(input []basebackuppolicyresources.BasePolicyRule, dsType basebackuppolicyresources.DataStoreTypes) any {
	if input == nil {
		return nil
	}

	for _, item := range input {
		if retentionRule, ok := item.(basebackuppolicyresources.AzureRetentionRule); ok && retentionRule.IsDefault != nil && *retentionRule.IsDefault {
			if len(retentionRule.Lifecycles) > 0 {
				if deleteOption, ok := retentionRule.Lifecycles[0].DeleteAfter.(basebackuppolicyresources.AbsoluteDeleteOption); ok {
					if retentionRule.Lifecycles[0].SourceDataStore.DataStoreType == dsType {
						return deleteOption.Duration
					}
				}
			}
		}
	}
	return nil
}

func flattenBackupPolicyBlobStorageVaultBackupRuleArray(input *[]basebackuppolicyresources.BasePolicyRule) []any {
	if input == nil {
		return make([]any, 0)
	}
	for _, item := range *input {
		if backupRule, ok := item.(basebackuppolicyresources.AzureBackupRule); ok {
			if backupRule.Trigger != nil {
				if scheduleBasedTrigger, ok := backupRule.Trigger.(basebackuppolicyresources.ScheduleBasedTriggerContext); ok {
					return pluginsdk.FlattenSlice(&scheduleBasedTrigger.Schedule.RepeatingTimeIntervals)
				}
			}
		}
	}
	return make([]any, 0)
}

func flattenBackupPolicyBlobStorageVaultBackupTimeZone(input *[]basebackuppolicyresources.BasePolicyRule) string {
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

func flattenBackupPolicyBlobStorageRetentionRuleArray(input *[]basebackuppolicyresources.BasePolicyRule) []any {
	results := make([]any, 0)
	if input == nil {
		return results
	}

	var taggingCriteriaList []basebackuppolicyresources.TaggingCriteria
	for _, item := range *input {
		if backupRule, ok := item.(basebackuppolicyresources.AzureBackupRule); ok {
			if trigger, ok := backupRule.Trigger.(basebackuppolicyresources.ScheduleBasedTriggerContext); ok {
				if trigger.TaggingCriteria != nil {
					taggingCriteriaList = trigger.TaggingCriteria
				}
			}
		}
	}

	for _, item := range *input {
		if retentionRule, ok := item.(basebackuppolicyresources.AzureRetentionRule); ok && (retentionRule.IsDefault == nil || !*retentionRule.IsDefault) {
			name := retentionRule.Name
			var taggingPriority int64
			var taggingCriteria []any
			for _, criteria := range taggingCriteriaList {
				if strings.EqualFold(criteria.TagInfo.TagName, name) {
					taggingPriority = criteria.TaggingPriority
					taggingCriteria = flattenBackupPolicyBlobStorageBackupCriteriaArray(criteria.Criteria)
					break
				}
			}

			var lifeCycle []any
			if v := retentionRule.Lifecycles; len(v) > 0 {
				lifeCycle = flattenBackupPolicyBlobStorageBackupLifeCycleArray(v, basebackuppolicyresources.DataStoreTypesVaultStore)
			}
			results = append(results, map[string]any{
				"name":       name,
				"priority":   taggingPriority,
				"criteria":   taggingCriteria,
				"life_cycle": lifeCycle,
			})
		}
	}
	return results
}

func flattenBackupPolicyBlobStorageBackupLifeCycleArray(input []basebackuppolicyresources.SourceLifeCycle, dsType basebackuppolicyresources.DataStoreTypes) []any {
	results := make([]any, 0)
	if input == nil {
		return results
	}

	for _, item := range input {
		var duration string
		dataStoreType := item.SourceDataStore.DataStoreType
		if deleteOption, ok := item.DeleteAfter.(basebackuppolicyresources.AbsoluteDeleteOption); ok {
			if dataStoreType == dsType {
				duration = deleteOption.Duration
			} else {
				continue
			}
		}

		results = append(results, map[string]any{
			"duration":        duration,
			"data_store_type": string(dataStoreType),
		})
	}
	return results
}

func flattenBackupPolicyBlobStorageBackupCriteriaArray(input *[]basebackuppolicyresources.BackupCriteria) []any {
	results := make([]any, 0)
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
			var daysOfMonth []int
			if criteria.DaysOfMonth != nil {
				daysOfMonth = make([]int, 0)
				for _, item := range *criteria.DaysOfMonth {
					daysOfMonth = append(daysOfMonth, int(pointer.From(item.Date)))
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

			results = append(results, map[string]any{
				"absolute_criteria":      absoluteCriteria,
				"days_of_week":           daysOfWeek,
				"days_of_month":          daysOfMonth,
				"months_of_year":         monthsOfYear,
				"weeks_of_month":         weeksOfMonth,
				"scheduled_backup_times": scheduleTimes,
			})
		}
	}
	return results
}
