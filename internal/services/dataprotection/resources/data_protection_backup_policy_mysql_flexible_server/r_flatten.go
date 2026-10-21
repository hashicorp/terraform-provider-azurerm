// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package data_protection_backup_policy_mysql_flexible_server

import (
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-07-01/basebackuppolicyresources"
)

func flattenBackupPolicyMySQLFlexibleServerBackupRules(input []basebackuppolicyresources.BasePolicyRule) []string {
	backupRules := make([]string, 0)

	for _, item := range input {
		if v, ok := item.(basebackuppolicyresources.AzureBackupRule); ok {
			if v.Trigger != nil {
				if scheduleBasedTrigger, ok := v.Trigger.(basebackuppolicyresources.ScheduleBasedTriggerContext); ok {
					backupRules = scheduleBasedTrigger.Schedule.RepeatingTimeIntervals
					return backupRules
				}
			}
		}
	}

	return backupRules
}

func flattenBackupPolicyMySQLFlexibleServerBackupTimeZone(input []basebackuppolicyresources.BasePolicyRule) string {
	var timeZone string

	for _, item := range input {
		if backupRule, ok := item.(basebackuppolicyresources.AzureBackupRule); ok {
			if backupRule.Trigger != nil {
				if scheduleBasedTrigger, ok := backupRule.Trigger.(basebackuppolicyresources.ScheduleBasedTriggerContext); ok {
					timeZone = pointer.From(scheduleBasedTrigger.Schedule.TimeZone)
					return timeZone
				}
			}
		}
	}

	return timeZone
}

func flattenBackupPolicyMySQLFlexibleServerDefaultRetentionRule(input []basebackuppolicyresources.BasePolicyRule) []BackupPolicyMySQLFlexibleServerDefaultRetentionRule {
	results := make([]BackupPolicyMySQLFlexibleServerDefaultRetentionRule, 0)

	for _, item := range input {
		if retentionRule, ok := item.(basebackuppolicyresources.AzureRetentionRule); ok {
			if pointer.From(retentionRule.IsDefault) {
				var lifeCycle []BackupPolicyMySQLFlexibleServerLifeCycle
				if v := retentionRule.Lifecycles; len(v) > 0 {
					lifeCycle = flattenBackupPolicyMySQLFlexibleServerLifeCycles(v)
				}

				results = append(results, BackupPolicyMySQLFlexibleServerDefaultRetentionRule{
					LifeCycle: lifeCycle,
				})
			}
		}
	}

	return results
}

func flattenBackupPolicyMySQLFlexibleServerRetentionRules(input []basebackuppolicyresources.BasePolicyRule) []BackupPolicyMySQLFlexibleServerRetentionRule {
	results := make([]BackupPolicyMySQLFlexibleServerRetentionRule, 0)
	var taggingCriteriaList []basebackuppolicyresources.TaggingCriteria

	for _, item := range input {
		if backupRule, ok := item.(basebackuppolicyresources.AzureBackupRule); ok {
			if trigger, ok := backupRule.Trigger.(basebackuppolicyresources.ScheduleBasedTriggerContext); ok {
				if trigger.TaggingCriteria != nil {
					taggingCriteriaList = trigger.TaggingCriteria
				}
			}
		}
	}

	for _, item := range input {
		if retentionRule, ok := item.(basebackuppolicyresources.AzureRetentionRule); ok {
			var name string
			var taggingPriority int64
			var taggingCriteria []BackupPolicyMySQLFlexibleServerCriteria

			if !pointer.From(retentionRule.IsDefault) {
				name = retentionRule.Name

				for _, criteria := range taggingCriteriaList {
					if strings.EqualFold(criteria.TagInfo.TagName, name) {
						taggingPriority = criteria.TaggingPriority
						taggingCriteria = flattenBackupPolicyMySQLFlexibleServerBackupCriteria(criteria.Criteria)
						break
					}
				}

				var lifeCycle []BackupPolicyMySQLFlexibleServerLifeCycle
				if v := retentionRule.Lifecycles; len(v) > 0 {
					lifeCycle = flattenBackupPolicyMySQLFlexibleServerLifeCycles(v)
				}

				results = append(results, BackupPolicyMySQLFlexibleServerRetentionRule{
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

func flattenBackupPolicyMySQLFlexibleServerLifeCycles(input []basebackuppolicyresources.SourceLifeCycle) []BackupPolicyMySQLFlexibleServerLifeCycle {
	results := make([]BackupPolicyMySQLFlexibleServerLifeCycle, 0)

	for _, item := range input {
		var duration string
		var dataStoreType string

		if deleteOption, ok := item.DeleteAfter.(basebackuppolicyresources.AbsoluteDeleteOption); ok {
			duration = deleteOption.Duration
		}

		dataStoreType = string(item.SourceDataStore.DataStoreType)

		results = append(results, BackupPolicyMySQLFlexibleServerLifeCycle{
			Duration:      duration,
			DataStoreType: dataStoreType,
		})
	}

	return results
}

func flattenBackupPolicyMySQLFlexibleServerBackupCriteria(input *[]basebackuppolicyresources.BackupCriteria) []BackupPolicyMySQLFlexibleServerCriteria {
	results := make([]BackupPolicyMySQLFlexibleServerCriteria, 0)
	if input == nil {
		return results
	}

	for _, item := range pointer.From(input) {
		if criteria, ok := item.(basebackuppolicyresources.ScheduleBasedBackupCriteria); ok {
			var absoluteCriteria string
			if criteria.AbsoluteCriteria != nil && len(pointer.From(criteria.AbsoluteCriteria)) > 0 {
				absoluteCriteria = string(pointer.From(criteria.AbsoluteCriteria)[0])
			}

			daysOfWeek := make([]string, 0)
			if criteria.DaysOfTheWeek != nil {
				for _, item := range pointer.From(criteria.DaysOfTheWeek) {
					daysOfWeek = append(daysOfWeek, string(item))
				}
			}

			monthsOfYear := make([]string, 0)
			if criteria.MonthsOfYear != nil {
				for _, item := range pointer.From(criteria.MonthsOfYear) {
					monthsOfYear = append(monthsOfYear, string(item))
				}
			}

			weeksOfMonth := make([]string, 0)
			if criteria.WeeksOfTheMonth != nil {
				for _, item := range pointer.From(criteria.WeeksOfTheMonth) {
					weeksOfMonth = append(weeksOfMonth, string(item))
				}
			}

			scheduleTimes := make([]string, 0)
			if criteria.ScheduleTimes != nil {
				scheduleTimes = append(scheduleTimes, pointer.From(criteria.ScheduleTimes)...)
			}

			results = append(results, BackupPolicyMySQLFlexibleServerCriteria{
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
