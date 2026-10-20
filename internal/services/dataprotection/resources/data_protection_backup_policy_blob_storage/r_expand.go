// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package data_protection_backup_policy_blob_storage

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-07-01/basebackuppolicyresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func expandBackupPolicyBlobStorageTaggingCriteriaArray(input []any) (*[]basebackuppolicyresources.TaggingCriteria, error) {
	results := []basebackuppolicyresources.TaggingCriteria{
		{
			IsDefault:       true,
			TaggingPriority: 99,
			TagInfo: basebackuppolicyresources.RetentionTag{
				Id:      pointer.To("Default_"),
				TagName: blobStorageVaultDefaultRetentionRuleName,
			},
		},
	}

	for _, item := range input {
		v := item.(map[string]any)
		result := basebackuppolicyresources.TaggingCriteria{
			TaggingPriority: int64(v["priority"].(int)),
			TagInfo: basebackuppolicyresources.RetentionTag{
				Id:      pointer.To(v["name"].(string) + "_"),
				TagName: v["name"].(string),
			},
		}

		criteria, err := expandBackupPolicyBlobStorageCriteriaArray(v["criteria"].([]any))
		if err != nil {
			return nil, err
		}
		result.Criteria = criteria
		results = append(results, result)
	}
	return &results, nil
}

func expandBackupPolicyBlobStorageCriteriaArray(input []any) (*[]basebackuppolicyresources.BackupCriteria, error) {
	if len(input) == 0 || input[0] == nil {
		return nil, fmt.Errorf("criteria is a required field, cannot leave blank")
	}
	results := make([]basebackuppolicyresources.BackupCriteria, 0)

	for _, item := range input {
		v := item.(map[string]any)
		var absoluteCriteria []basebackuppolicyresources.AbsoluteMarker
		if absoluteCriteriaRaw := v["absolute_criteria"].(string); len(absoluteCriteriaRaw) > 0 {
			absoluteCriteria = []basebackuppolicyresources.AbsoluteMarker{basebackuppolicyresources.AbsoluteMarker(absoluteCriteriaRaw)}
		}

		var daysOfMonth []basebackuppolicyresources.Day
		if v["days_of_month"].(*pluginsdk.Set).Len() > 0 {
			daysOfMonth = make([]basebackuppolicyresources.Day, 0)
			for _, value := range v["days_of_month"].(*pluginsdk.Set).List() {
				isLast := false
				if value == 0 {
					isLast = true
				}
				daysOfMonth = append(daysOfMonth, basebackuppolicyresources.Day{
					Date: pointer.To(int64(value.(int))), IsLast: pointer.To(isLast),
				})
			}
		}

		var daysOfWeek []basebackuppolicyresources.DayOfWeek
		if v["days_of_week"].(*pluginsdk.Set).Len() > 0 {
			daysOfWeek = make([]basebackuppolicyresources.DayOfWeek, 0)
			for _, value := range v["days_of_week"].(*pluginsdk.Set).List() {
				daysOfWeek = append(daysOfWeek, basebackuppolicyresources.DayOfWeek(value.(string)))
			}
		}

		var monthsOfYear []basebackuppolicyresources.Month
		if v["months_of_year"].(*pluginsdk.Set).Len() > 0 {
			monthsOfYear = make([]basebackuppolicyresources.Month, 0)
			for _, value := range v["months_of_year"].(*pluginsdk.Set).List() {
				monthsOfYear = append(monthsOfYear, basebackuppolicyresources.Month(value.(string)))
			}
		}

		var scheduleTimes *[]string
		if v["scheduled_backup_times"].(*pluginsdk.Set).Len() > 0 {
			scheduleTimes = pluginsdk.ExpandStringSlice(v["scheduled_backup_times"].(*pluginsdk.Set).List())
		}

		var weeksOfMonth []basebackuppolicyresources.WeekNumber
		if v["weeks_of_month"].(*pluginsdk.Set).Len() > 0 {
			weeksOfMonth = make([]basebackuppolicyresources.WeekNumber, 0)
			for _, value := range v["weeks_of_month"].(*pluginsdk.Set).List() {
				weeksOfMonth = append(weeksOfMonth, basebackuppolicyresources.WeekNumber(value.(string)))
			}
		}

		results = append(results, basebackuppolicyresources.ScheduleBasedBackupCriteria{
			AbsoluteCriteria: &absoluteCriteria,
			DaysOfMonth:      &daysOfMonth,
			DaysOfTheWeek:    &daysOfWeek,
			MonthsOfYear:     &monthsOfYear,
			ScheduleTimes:    scheduleTimes,
			WeeksOfTheMonth:  &weeksOfMonth,
		})
	}
	return &results, nil
}

func expandBackupPolicyBlobStorageAzureBackupRuleArray(input []any, timeZone string, taggingCriteria *[]basebackuppolicyresources.TaggingCriteria) []basebackuppolicyresources.BasePolicyRule {
	results := make([]basebackuppolicyresources.BasePolicyRule, 0)
	results = append(results, basebackuppolicyresources.AzureBackupRule{
		Name: "BackupIntervals",
		DataStore: basebackuppolicyresources.DataStoreInfoBase{
			DataStoreType: basebackuppolicyresources.DataStoreTypesVaultStore,
			ObjectType:    "DataStoreInfoBase",
		},
		BackupParameters: &basebackuppolicyresources.AzureBackupParams{
			BackupType: "Discrete",
		},
		Trigger: basebackuppolicyresources.ScheduleBasedTriggerContext{
			Schedule: basebackuppolicyresources.BackupSchedule{
				RepeatingTimeIntervals: *pluginsdk.ExpandStringSlice(input),
				TimeZone:               pointer.To(timeZone),
			},
			TaggingCriteria: *taggingCriteria,
		},
	})

	return results
}

func expandBackupPolicyBlobStorageDefaultRetentionRuleArray(input any, dataStoreType basebackuppolicyresources.DataStoreTypes) basebackuppolicyresources.BasePolicyRule {
	name := blobStorageVaultDefaultRetentionRuleName
	if dataStoreType == basebackuppolicyresources.DataStoreTypesOperationalStore {
		name = blobStorageOperationalDefaultRetentionRuleName
	}

	return basebackuppolicyresources.AzureRetentionRule{
		Name:      name,
		IsDefault: pointer.To(true),
		Lifecycles: []basebackuppolicyresources.SourceLifeCycle{
			{
				DeleteAfter: basebackuppolicyresources.AbsoluteDeleteOption{
					Duration: input.(string),
				},
				SourceDataStore: basebackuppolicyresources.DataStoreInfoBase{
					DataStoreType: dataStoreType,
					ObjectType:    "DataStoreInfoBase",
				},
				TargetDataStoreCopySettings: &[]basebackuppolicyresources.TargetCopySetting{},
			},
		},
	}
}

func expandBackupPolicyBlobStorageAzureRetentionRuleArray(input []any) []basebackuppolicyresources.BasePolicyRule {
	results := make([]basebackuppolicyresources.BasePolicyRule, 0)
	for _, item := range input {
		v := item.(map[string]any)
		results = append(results, basebackuppolicyresources.AzureRetentionRule{
			Name:       v["name"].(string),
			IsDefault:  pointer.To(false),
			Lifecycles: expandBackupPolicyBlobStorageLifeCycle(v["life_cycle"].([]any)),
		})
	}
	return results
}

func expandBackupPolicyBlobStorageLifeCycle(input []any) []basebackuppolicyresources.SourceLifeCycle {
	results := make([]basebackuppolicyresources.SourceLifeCycle, 0)
	for _, item := range input {
		v := item.(map[string]any)
		sourceLifeCycle := basebackuppolicyresources.SourceLifeCycle{
			DeleteAfter: basebackuppolicyresources.AbsoluteDeleteOption{
				Duration: v["duration"].(string),
			},
			SourceDataStore: basebackuppolicyresources.DataStoreInfoBase{
				DataStoreType: basebackuppolicyresources.DataStoreTypes(v["data_store_type"].(string)),
				ObjectType:    "DataStoreInfoBase",
			},
			TargetDataStoreCopySettings: &[]basebackuppolicyresources.TargetCopySetting{},
		}
		results = append(results, sourceLifeCycle)
	}

	return results
}
