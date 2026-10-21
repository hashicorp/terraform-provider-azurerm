// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_protection_backup_policy_mysql_flexible_server

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
			var model BackupPolicyMySQLFlexibleServerModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.DataProtection.BackupPolicyClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			vaultId, _ := basebackuppolicyresources.ParseBackupVaultID(model.VaultId)
			id := basebackuppolicyresources.NewBackupPolicyID(subscriptionId, vaultId.ResourceGroupName, vaultId.BackupVaultName, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.BackupPoliciesGet(ctx, id)
				if err != nil {
					if !response.WasNotFound(existing.HttpResponse) {
						return fmt.Errorf("checking for existing %s: %+v", id, err)
					}
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			policyRules := make([]basebackuppolicyresources.BasePolicyRule, 0)
			policyRules = append(policyRules, expandBackupPolicyMySQLFlexibleServerAzureBackupRules(model.BackupRepeatingTimeIntervals, model.TimeZone, expandBackupPolicyMySQLFlexibleServerTaggingCriteria(model.RetentionRules))...)
			policyRules = append(policyRules, expandBackupPolicyMySQLFlexibleServerDefaultAzureRetentionRule(model.DefaultRetentionRule))
			policyRules = append(policyRules, expandBackupPolicyMySQLFlexibleServerAzureRetentionRules(model.RetentionRules)...)

			parameters := basebackuppolicyresources.BaseBackupPolicyResource{
				Properties: basebackuppolicyresources.BackupPolicy{
					PolicyRules:     policyRules,
					DatasourceTypes: []string{"Microsoft.DBforMySQL/flexibleServers"},
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

func expandBackupPolicyMySQLFlexibleServerAzureBackupRules(input []string, timeZone string, taggingCriteria []basebackuppolicyresources.TaggingCriteria) []basebackuppolicyresources.BasePolicyRule {
	results := make([]basebackuppolicyresources.BasePolicyRule, 0)

	results = append(results, basebackuppolicyresources.AzureBackupRule{
		Name: "BackupIntervals",
		DataStore: basebackuppolicyresources.DataStoreInfoBase{
			DataStoreType: basebackuppolicyresources.DataStoreTypesVaultStore,
			ObjectType:    "DataStoreInfoBase",
		},
		BackupParameters: basebackuppolicyresources.AzureBackupParams{
			BackupType: "Full",
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

func expandBackupPolicyMySQLFlexibleServerAzureRetentionRules(input []BackupPolicyMySQLFlexibleServerRetentionRule) []basebackuppolicyresources.BasePolicyRule {
	results := make([]basebackuppolicyresources.BasePolicyRule, 0)

	for _, item := range input {
		results = append(results, basebackuppolicyresources.AzureRetentionRule{
			Name:       item.Name,
			IsDefault:  pointer.To(false),
			Lifecycles: expandBackupPolicyMySQLFlexibleServerLifeCycle(item.LifeCycle),
		})
	}

	return results
}

func expandBackupPolicyMySQLFlexibleServerDefaultAzureRetentionRule(input []BackupPolicyMySQLFlexibleServerDefaultRetentionRule) basebackuppolicyresources.BasePolicyRule {
	result := basebackuppolicyresources.AzureRetentionRule{
		Name:      "Default",
		IsDefault: pointer.To(true),
	}

	if len(input) > 0 {
		result.Lifecycles = expandBackupPolicyMySQLFlexibleServerLifeCycle(input[0].LifeCycle)
	}

	return result
}

func expandBackupPolicyMySQLFlexibleServerLifeCycle(input []BackupPolicyMySQLFlexibleServerLifeCycle) []basebackuppolicyresources.SourceLifeCycle {
	results := make([]basebackuppolicyresources.SourceLifeCycle, 0)

	for _, item := range input {
		sourceLifeCycle := basebackuppolicyresources.SourceLifeCycle{
			DeleteAfter: basebackuppolicyresources.AbsoluteDeleteOption{
				Duration: item.Duration,
			},
			SourceDataStore: basebackuppolicyresources.DataStoreInfoBase{
				DataStoreType: basebackuppolicyresources.DataStoreTypes(item.DataStoreType),
				ObjectType:    "DataStoreInfoBase",
			},
			TargetDataStoreCopySettings: &[]basebackuppolicyresources.TargetCopySetting{},
		}

		results = append(results, sourceLifeCycle)
	}

	return results
}

func expandBackupPolicyMySQLFlexibleServerTaggingCriteria(input []BackupPolicyMySQLFlexibleServerRetentionRule) []basebackuppolicyresources.TaggingCriteria {
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
		result := basebackuppolicyresources.TaggingCriteria{
			Criteria:        expandBackupPolicyMySQLFlexibleServerCriteria(item.Criteria),
			TaggingPriority: item.Priority,
			TagInfo: basebackuppolicyresources.RetentionTag{
				Id:      pointer.To(item.Name + "_"),
				TagName: item.Name,
			},
		}

		results = append(results, result)
	}

	return results
}

func expandBackupPolicyMySQLFlexibleServerCriteria(input []BackupPolicyMySQLFlexibleServerCriteria) *[]basebackuppolicyresources.BackupCriteria {
	if len(input) == 0 {
		return nil
	}

	results := make([]basebackuppolicyresources.BackupCriteria, 0)

	for _, item := range input {
		var absoluteCriteria []basebackuppolicyresources.AbsoluteMarker
		if absoluteCriteriaRaw := item.AbsoluteCriteria; len(absoluteCriteriaRaw) > 0 {
			absoluteCriteria = []basebackuppolicyresources.AbsoluteMarker{basebackuppolicyresources.AbsoluteMarker(absoluteCriteriaRaw)}
		}

		var daysOfWeek []basebackuppolicyresources.DayOfWeek
		if len(item.DaysOfWeek) > 0 {
			daysOfWeek = make([]basebackuppolicyresources.DayOfWeek, 0)
			for _, value := range item.DaysOfWeek {
				daysOfWeek = append(daysOfWeek, basebackuppolicyresources.DayOfWeek(value))
			}
		}

		var monthsOfYear []basebackuppolicyresources.Month
		if len(item.MonthsOfYear) > 0 {
			monthsOfYear = make([]basebackuppolicyresources.Month, 0)
			for _, value := range item.MonthsOfYear {
				monthsOfYear = append(monthsOfYear, basebackuppolicyresources.Month(value))
			}
		}

		var weeksOfMonth []basebackuppolicyresources.WeekNumber
		if len(item.WeeksOfMonth) > 0 {
			weeksOfMonth = make([]basebackuppolicyresources.WeekNumber, 0)
			for _, value := range item.WeeksOfMonth {
				weeksOfMonth = append(weeksOfMonth, basebackuppolicyresources.WeekNumber(value))
			}
		}

		var scheduleTimes []string
		if len(item.ScheduledBackupTimes) > 0 {
			scheduleTimes = item.ScheduledBackupTimes
		}

		results = append(results, basebackuppolicyresources.ScheduleBasedBackupCriteria{
			AbsoluteCriteria: pointer.To(absoluteCriteria),
			DaysOfTheWeek:    pointer.To(daysOfWeek),
			MonthsOfYear:     pointer.To(monthsOfYear),
			ScheduleTimes:    pointer.To(scheduleTimes),
			WeeksOfTheMonth:  pointer.To(weeksOfMonth),
		})
	}

	return &results
}
