// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package data_protection_backup_policy_disk

import (
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-07-01/basebackuppolicyresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func flattenBackupPolicyDiskBackupRuleArray(input *[]basebackuppolicyresources.BasePolicyRule) []any {
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

func flattenBackupPolicyDiskBackupTimeZone(input *[]basebackuppolicyresources.BasePolicyRule) string {
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

func flattenBackupPolicyDiskDefaultRetentionRuleDuration(input *[]basebackuppolicyresources.BasePolicyRule) any {
	if input == nil {
		return nil
	}

	for _, item := range *input {
		if retentionRule, ok := item.(basebackuppolicyresources.AzureRetentionRule); ok && retentionRule.IsDefault != nil && *retentionRule.IsDefault {
			if len(retentionRule.Lifecycles) > 0 {
				if deleteOption, ok := retentionRule.Lifecycles[0].DeleteAfter.(basebackuppolicyresources.AbsoluteDeleteOption); ok {
					return deleteOption.Duration
				}
			}
		}
	}
	return nil
}

func flattenBackupPolicyDiskRetentionRuleArray(input *[]basebackuppolicyresources.BasePolicyRule) []any {
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
					taggingCriteria = flattenBackupPolicyDiskBackupCriteriaArray(criteria.Criteria)
				}
			}
			var duration string
			if len(retentionRule.Lifecycles) > 0 {
				if deleteOption, ok := retentionRule.Lifecycles[0].DeleteAfter.(basebackuppolicyresources.AbsoluteDeleteOption); ok {
					duration = deleteOption.Duration
				}
			}
			results = append(results, map[string]any{
				"name":     name,
				"priority": taggingPriority,
				"criteria": taggingCriteria,
				"duration": duration,
			})
		}
	}
	return results
}

func flattenBackupPolicyDiskBackupCriteriaArray(input *[]basebackuppolicyresources.BackupCriteria) []any {
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

			results = append(results, map[string]any{
				"absolute_criteria": absoluteCriteria,
			})
		}
	}
	return results
}
