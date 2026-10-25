// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_protection_backup_policy_disk

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-07-01/basebackuppolicyresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func expandBackupPolicyDiskAzureBackupRuleArray(input []any, timeZone string, taggingCriteria *[]basebackuppolicyresources.TaggingCriteria) []basebackuppolicyresources.BasePolicyRule {
	results := make([]basebackuppolicyresources.BasePolicyRule, 0)

	results = append(results, basebackuppolicyresources.AzureBackupRule{
		Name: "BackupIntervals",
		DataStore: basebackuppolicyresources.DataStoreInfoBase{
			DataStoreType: basebackuppolicyresources.DataStoreTypesOperationalStore,
			ObjectType:    "DataStoreInfoBase",
		},
		BackupParameters: &basebackuppolicyresources.AzureBackupParams{
			BackupType: "Incremental",
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

func expandBackupPolicyDiskAzureRetentionRuleArray(input []any) []basebackuppolicyresources.BasePolicyRule {
	results := make([]basebackuppolicyresources.BasePolicyRule, 0)
	for _, item := range input {
		v := item.(map[string]any)
		results = append(results, basebackuppolicyresources.AzureRetentionRule{
			Name:      v["name"].(string),
			IsDefault: pointer.To(false),
			Lifecycles: []basebackuppolicyresources.SourceLifeCycle{
				{
					DeleteAfter: basebackuppolicyresources.AbsoluteDeleteOption{
						Duration: v["duration"].(string),
					},
					SourceDataStore: basebackuppolicyresources.DataStoreInfoBase{
						DataStoreType: "OperationalStore",
						ObjectType:    "DataStoreInfoBase",
					},
					TargetDataStoreCopySettings: &[]basebackuppolicyresources.TargetCopySetting{},
				},
			},
		})
	}
	return results
}

func expandBackupPolicyDiskDefaultAzureRetentionRule(input any) basebackuppolicyresources.BasePolicyRule {
	return basebackuppolicyresources.AzureRetentionRule{
		Name:      "Default",
		IsDefault: pointer.To(true),
		Lifecycles: []basebackuppolicyresources.SourceLifeCycle{
			{
				DeleteAfter: basebackuppolicyresources.AbsoluteDeleteOption{
					Duration: input.(string),
				},
				SourceDataStore: basebackuppolicyresources.DataStoreInfoBase{
					DataStoreType: "OperationalStore",
					ObjectType:    "DataStoreInfoBase",
				},
				TargetDataStoreCopySettings: &[]basebackuppolicyresources.TargetCopySetting{},
			},
		},
	}
}

func expandBackupPolicyDiskTaggingCriteriaArray(input []any) *[]basebackuppolicyresources.TaggingCriteria {
	results := []basebackuppolicyresources.TaggingCriteria{
		{
			IsDefault:       true,
			TaggingPriority: 99,
			TagInfo: basebackuppolicyresources.RetentionTag{
				Id:      pointer.To("Default_"),
				TagName: "Default",
			},
		},
	}
	for _, item := range input {
		v := item.(map[string]any)
		results = append(results, basebackuppolicyresources.TaggingCriteria{
			Criteria:        expandBackupPolicyDiskCriteriaArray(v["criteria"].([]any)),
			TaggingPriority: int64(v["priority"].(int)),
			TagInfo: basebackuppolicyresources.RetentionTag{
				Id:      pointer.To(v["name"].(string) + "_"),
				TagName: v["name"].(string),
			},
		})
	}
	return &results
}

func expandBackupPolicyDiskCriteriaArray(input []any) *[]basebackuppolicyresources.BackupCriteria {
	results := make([]basebackuppolicyresources.BackupCriteria, 0)
	for _, item := range input {
		v := item.(map[string]any)
		var absoluteCriteria []basebackuppolicyresources.AbsoluteMarker
		if absoluteCriteriaRaw := v["absolute_criteria"].(string); len(absoluteCriteriaRaw) > 0 {
			absoluteCriteria = []basebackuppolicyresources.AbsoluteMarker{basebackuppolicyresources.AbsoluteMarker(absoluteCriteriaRaw)}
		}
		results = append(results, basebackuppolicyresources.ScheduleBasedBackupCriteria{
			AbsoluteCriteria: &absoluteCriteria,
		})
	}
	return &results
}
