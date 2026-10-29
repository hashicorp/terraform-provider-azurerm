// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_protection_backup_policy_postgresql_flexible_server

import (
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-07-01/basebackuppolicyresources"
)

func flattenBackupPolicyPostgreSQLFlexibleServerBackupRules(input []basebackuppolicyresources.BasePolicyRule) []string {
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

func flattenBackupPolicyPostgreSQLFlexibleServerBackupTimeZone(input []basebackuppolicyresources.BasePolicyRule) string {
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

func flattenBackupPolicyPostgreSQLFlexibleServerDefaultRetentionRule(input []basebackuppolicyresources.BasePolicyRule) []BackupPolicyPostgreSQLFlexibleServerDefaultRetentionRule {
	results := make([]BackupPolicyPostgreSQLFlexibleServerDefaultRetentionRule, 0)

	for _, item := range input {
		if retentionRule, ok := item.(basebackuppolicyresources.AzureRetentionRule); ok {
			if pointer.From(retentionRule.IsDefault) {
				var lifeCycle []BackupPolicyPostgreSQLFlexibleServerLifeCycle
				if v := retentionRule.Lifecycles; len(v) > 0 {
					lifeCycle = flattenBackupPolicyPostgreSQLFlexibleServerLifeCycles(v)
				}

				results = append(results, BackupPolicyPostgreSQLFlexibleServerDefaultRetentionRule{
					LifeCycle: lifeCycle,
				})
			}
		}
	}

	return results
}

func flattenBackupPolicyPostgreSQLFlexibleServerRetentionRules(input []basebackuppolicyresources.BasePolicyRule) []BackupPolicyPostgreSQLFlexibleServerRetentionRule {
	results := make([]BackupPolicyPostgreSQLFlexibleServerRetentionRule, 0)
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
			var taggingCriteria []BackupPolicyPostgreSQLFlexibleServerCriteria

			if !pointer.From(retentionRule.IsDefault) {
				name = retentionRule.Name

				for _, criteria := range taggingCriteriaList {
					if strings.EqualFold(criteria.TagInfo.TagName, name) {
						taggingPriority = criteria.TaggingPriority
						taggingCriteria = flattenBackupPolicyPostgreSQLFlexibleServerBackupCriteria(criteria.Criteria)
						break
					}
				}

				var lifeCycle []BackupPolicyPostgreSQLFlexibleServerLifeCycle
				if v := retentionRule.Lifecycles; len(v) > 0 {
					lifeCycle = flattenBackupPolicyPostgreSQLFlexibleServerLifeCycles(v)
				}

				results = append(results, BackupPolicyPostgreSQLFlexibleServerRetentionRule{
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

func flattenBackupPolicyPostgreSQLFlexibleServerLifeCycles(input []basebackuppolicyresources.SourceLifeCycle) []BackupPolicyPostgreSQLFlexibleServerLifeCycle {
	results := make([]BackupPolicyPostgreSQLFlexibleServerLifeCycle, 0)

	for _, item := range input {
		var duration string
		var dataStoreType string

		if deleteOption, ok := item.DeleteAfter.(basebackuppolicyresources.AbsoluteDeleteOption); ok {
			duration = deleteOption.Duration
		}

		dataStoreType = string(item.SourceDataStore.DataStoreType)

		results = append(results, BackupPolicyPostgreSQLFlexibleServerLifeCycle{
			Duration:      duration,
			DataStoreType: dataStoreType,
		})
	}

	return results
}

func flattenBackupPolicyPostgreSQLFlexibleServerBackupCriteria(input *[]basebackuppolicyresources.BackupCriteria) []BackupPolicyPostgreSQLFlexibleServerCriteria {
	results := make([]BackupPolicyPostgreSQLFlexibleServerCriteria, 0)
	if input == nil {
		return results
	}

	for _, item := range pointer.From(input) {
		if criteria, ok := item.(basebackuppolicyresources.ScheduleBasedBackupCriteria); ok {
			var absoluteCriteria string
			if criteria.AbsoluteCriteria != nil && len(pointer.From(criteria.AbsoluteCriteria)) > 0 {
				absoluteCriteria = string(pointer.From(criteria.AbsoluteCriteria)[0])
			}

			var daysOfWeek []string
			if criteria.DaysOfTheWeek != nil {
				daysOfWeek = make([]string, 0)

				for _, item := range pointer.From(criteria.DaysOfTheWeek) {
					daysOfWeek = append(daysOfWeek, string(item))
				}
			}

			var monthsOfYear []string
			if criteria.MonthsOfYear != nil {
				monthsOfYear = make([]string, 0)

				for _, item := range pointer.From(criteria.MonthsOfYear) {
					monthsOfYear = append(monthsOfYear, string(item))
				}
			}

			var weeksOfMonth []string
			if criteria.WeeksOfTheMonth != nil {
				weeksOfMonth = make([]string, 0)

				for _, item := range pointer.From(criteria.WeeksOfTheMonth) {
					weeksOfMonth = append(weeksOfMonth, string(item))
				}
			}

			var scheduleTimes []string
			if criteria.ScheduleTimes != nil {
				scheduleTimes = make([]string, 0)
				scheduleTimes = append(scheduleTimes, pointer.From(criteria.ScheduleTimes)...)
			}

			results = append(results, BackupPolicyPostgreSQLFlexibleServerCriteria{
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
