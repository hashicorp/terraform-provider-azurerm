// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_protection_backup_policy_data_lake_storage

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-07-01/basebackuppolicyresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DataProtection.BackupPolicyClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model BackupPolicyDataLakeStorageModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			for _, rule := range model.RetentionRules {
				if rule.AbsoluteCriteria == "" && len(rule.DaysOfWeek) == 0 {
					return fmt.Errorf("`retention_rule` %q requires at least one of `absolute_criteria` and `days_of_week` to be specified", rule.Name)
				}
			}

			vaultId, _ := basebackuppolicyresources.ParseBackupVaultID(model.DataProtectionBackupVaultId)
			id := basebackuppolicyresources.NewBackupPolicyID(subscriptionId, vaultId.ResourceGroupName, vaultId.BackupVaultName, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.BackupPoliciesGet(ctx, id)
				if err != nil {
					if !response.WasNotFound(existing.HttpResponse) {
						return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
					}
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			policyRules := make([]basebackuppolicyresources.BasePolicyRule, 0)
			policyRules = append(policyRules, expandBackupPolicyDataLakeStorageAzureRetentionRules(model.RetentionRules)...)
			policyRules = append(policyRules, expandBackupPolicyDataLakeStorageDefaultAzureRetentionRule(model.DefaultRetentionDuration))
			policyRules = append(policyRules, expandBackupPolicyDataLakeStorageAzureBackupRules(model.BackupSchedule, model.TimeZone, expandBackupPolicyDataLakeStorageTaggingCriteria(model.RetentionRules))...)

			parameters := basebackuppolicyresources.BaseBackupPolicyResource{
				Properties: &basebackuppolicyresources.BackupPolicy{
					ObjectType:      "BackupPolicy",
					PolicyRules:     policyRules,
					DatasourceTypes: []string{"Microsoft.Storage/storageAccounts/adlsBlobServices"},
				},
			}

			if _, err := client.BackupPoliciesCreateOrUpdate(ctx, id, parameters); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, &id); err != nil {
				return err
			}

			return nil
		},
	}
}

func expandBackupPolicyDataLakeStorageAzureBackupRules(input []string, timeZone string, taggingCriteria []basebackuppolicyresources.TaggingCriteria) []basebackuppolicyresources.BasePolicyRule {
	results := make([]basebackuppolicyresources.BasePolicyRule, 0)

	results = append(results, basebackuppolicyresources.AzureBackupRule{
		Name: "BackupRule",
		DataStore: basebackuppolicyresources.DataStoreInfoBase{
			DataStoreType: basebackuppolicyresources.DataStoreTypesVaultStore,
			ObjectType:    "DataStoreInfoBase",
		},
		BackupParameters: basebackuppolicyresources.AzureBackupParams{
			BackupType: "Discrete",
		},
		Trigger: basebackuppolicyresources.ScheduleBasedTriggerContext{
			Schedule: basebackuppolicyresources.BackupSchedule{
				RepeatingTimeIntervals: input,
				TimeZone:               pointer.To(timeZone),
			},
			TaggingCriteria: taggingCriteria,
		},
	})

	return results
}

func expandBackupPolicyDataLakeStorageAzureRetentionRules(input []BackupPolicyDataLakeStorageRetentionRule) []basebackuppolicyresources.BasePolicyRule {
	results := make([]basebackuppolicyresources.BasePolicyRule, 0)

	for _, item := range input {
		results = append(results, basebackuppolicyresources.AzureRetentionRule{
			Name:       item.Name,
			IsDefault:  pointer.To(false),
			Lifecycles: expandBackupPolicyDataLakeStorageLifeCycle(item.Duration),
		})
	}

	return results
}

func expandBackupPolicyDataLakeStorageDefaultAzureRetentionRule(duration string) basebackuppolicyresources.BasePolicyRule {
	return basebackuppolicyresources.AzureRetentionRule{
		Name:       "Default",
		IsDefault:  pointer.To(true),
		Lifecycles: expandBackupPolicyDataLakeStorageLifeCycle(duration),
	}
}

func expandBackupPolicyDataLakeStorageLifeCycle(duration string) []basebackuppolicyresources.SourceLifeCycle {
	// NOTE: currently only `VaultStore` is supported by the service team. When `ArchiveStore` is supported
	// in the future, export `data_store_type` as a schema field and use `VaultStore` as the default value.
	return []basebackuppolicyresources.SourceLifeCycle{
		{
			DeleteAfter: basebackuppolicyresources.AbsoluteDeleteOption{
				Duration: duration,
			},
			SourceDataStore: basebackuppolicyresources.DataStoreInfoBase{
				DataStoreType: basebackuppolicyresources.DataStoreTypesVaultStore,
				ObjectType:    "DataStoreInfoBase",
			},
			TargetDataStoreCopySettings: &[]basebackuppolicyresources.TargetCopySetting{},
		},
	}
}

func expandBackupPolicyDataLakeStorageTaggingCriteria(input []BackupPolicyDataLakeStorageRetentionRule) []basebackuppolicyresources.TaggingCriteria {
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

	for i, item := range input {
		result := basebackuppolicyresources.TaggingCriteria{
			Criteria:        expandBackupPolicyDataLakeStorageRetentionRuleCriteria(item),
			TaggingPriority: int64(i + 1),
			TagInfo: basebackuppolicyresources.RetentionTag{
				Id:      pointer.To(item.Name + "_"),
				TagName: item.Name,
			},
		}

		results = append(results, result)
	}

	return results
}

func expandBackupPolicyDataLakeStorageRetentionRuleCriteria(input BackupPolicyDataLakeStorageRetentionRule) *[]basebackuppolicyresources.BackupCriteria {
	var absoluteCriteria []basebackuppolicyresources.AbsoluteMarker
	if len(input.AbsoluteCriteria) > 0 {
		absoluteCriteria = []basebackuppolicyresources.AbsoluteMarker{basebackuppolicyresources.AbsoluteMarker(input.AbsoluteCriteria)}
	}

	var daysOfWeek []basebackuppolicyresources.DayOfWeek
	if len(input.DaysOfWeek) > 0 {
		daysOfWeek = make([]basebackuppolicyresources.DayOfWeek, 0)
		for _, value := range input.DaysOfWeek {
			daysOfWeek = append(daysOfWeek, basebackuppolicyresources.DayOfWeek(value))
		}
	}

	var monthsOfYear []basebackuppolicyresources.Month
	if len(input.MonthsOfYear) > 0 {
		monthsOfYear = make([]basebackuppolicyresources.Month, 0)
		for _, value := range input.MonthsOfYear {
			monthsOfYear = append(monthsOfYear, basebackuppolicyresources.Month(value))
		}
	}

	var weeksOfMonth []basebackuppolicyresources.WeekNumber
	if len(input.WeeksOfMonth) > 0 {
		weeksOfMonth = make([]basebackuppolicyresources.WeekNumber, 0)
		for _, value := range input.WeeksOfMonth {
			weeksOfMonth = append(weeksOfMonth, basebackuppolicyresources.WeekNumber(value))
		}
	}

	var scheduleTimes []string
	if len(input.ScheduledBackupTimes) > 0 {
		scheduleTimes = input.ScheduledBackupTimes
	}

	if len(absoluteCriteria) == 0 && len(daysOfWeek) == 0 && len(monthsOfYear) == 0 && len(weeksOfMonth) == 0 && len(scheduleTimes) == 0 {
		return nil
	}

	return &[]basebackuppolicyresources.BackupCriteria{
		basebackuppolicyresources.ScheduleBasedBackupCriteria{
			AbsoluteCriteria: pointer.To(absoluteCriteria),
			DaysOfTheWeek:    pointer.To(daysOfWeek),
			MonthsOfYear:     pointer.To(monthsOfYear),
			ScheduleTimes:    pointer.To(scheduleTimes),
			WeeksOfTheMonth:  pointer.To(weeksOfMonth),
		},
	}
}
