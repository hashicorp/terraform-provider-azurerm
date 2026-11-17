// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_protection_backup_policy_cosmosdb_account

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2026-06-01/basebackuppolicyresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DataProtection.BackupPolicyClient20260601
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model BackupPolicyCosmosdbAccountModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			vaultId, err := basebackuppolicyresources.ParseBackupVaultID(model.DataProtectionBackupVaultId)
			if err != nil {
				return err
			}
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

			retentionRules := expandBackupPolicyCosmosdbAccountRetentionRuleModels(model)
			policyRules := make([]basebackuppolicyresources.BasePolicyRule, 0)
			policyRules = append(policyRules, expandBackupPolicyCosmosdbAccountRetentionRules(retentionRules)...)
			policyRules = append(policyRules, expandBackupPolicyCosmosdbAccountDefaultRetentionRule(model.DefaultRetentionDuration))
			backupRules, err := expandBackupPolicyCosmosdbAccountBackupRules(model.BackupSchedule, model.BackupFrequency == backupPolicyCosmosdbAccountFrequencyDaily, expandBackupPolicyCosmosdbAccountTaggingCriteria(retentionRules))
			if err != nil {
				return fmt.Errorf("parsing `backup schedule`: %+v", err)
			}
			policyRules = append(policyRules, backupRules...)

			parameters := basebackuppolicyresources.BaseBackupPolicyResource{
				Properties: &basebackuppolicyresources.BackupPolicy{
					ObjectType:      "BackupPolicy",
					PolicyRules:     policyRules,
					DatasourceTypes: []string{"Microsoft.DocumentDB/databaseAccounts"},
				},
			}
			if _, err := client.BackupPoliciesCreateOrUpdate(ctx, id, parameters); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return pluginsdk.SetResourceIdentityData(metadata.ResourceData, &id)
		},
	}
}

func expandBackupPolicyCosmosdbAccountRetentionRuleModels(input BackupPolicyCosmosdbAccountModel) []BackupPolicyCosmosdbAccountRetentionRule {
	results := make([]BackupPolicyCosmosdbAccountRetentionRule, 0)
	for _, rule := range input.YearlyRetentionRule {
		result := BackupPolicyCosmosdbAccountRetentionRule{
			Name:         "Yearly",
			Duration:     rule.Duration,
			DaysOfWeek:   rule.DaysOfWeek,
			WeeksOfMonth: rule.WeeksOfMonth,
			MonthsOfYear: rule.MonthsOfYear,
		}
		if len(rule.DaysOfWeek) == 0 {
			result.BackupOccurrence = string(basebackuppolicyresources.AbsoluteMarkerFirstOfYear)
		}
		results = append(results, result)
	}
	for _, rule := range input.MonthlyRetentionRule {
		result := BackupPolicyCosmosdbAccountRetentionRule{
			Name:         "Monthly",
			Duration:     rule.Duration,
			DaysOfWeek:   rule.DaysOfWeek,
			WeeksOfMonth: rule.WeeksOfMonth,
		}
		if len(rule.DaysOfWeek) == 0 {
			result.BackupOccurrence = string(basebackuppolicyresources.AbsoluteMarkerFirstOfMonth)
		}
		results = append(results, result)
	}
	for _, rule := range input.WeeklyRetentionRule {
		result := BackupPolicyCosmosdbAccountRetentionRule{
			Name:       "Weekly",
			Duration:   rule.Duration,
			DaysOfWeek: rule.DaysOfWeek,
		}
		if len(rule.DaysOfWeek) == 0 {
			result.BackupOccurrence = string(basebackuppolicyresources.AbsoluteMarkerFirstOfWeek)
		}
		results = append(results, result)
	}
	return results
}

func expandBackupPolicyCosmosdbAccountRetentionRules(input []BackupPolicyCosmosdbAccountRetentionRule) []basebackuppolicyresources.BasePolicyRule {
	results := make([]basebackuppolicyresources.BasePolicyRule, 0)

	for _, item := range input {
		results = append(results, basebackuppolicyresources.AzureRetentionRule{
			Name:       item.Name,
			IsDefault:  pointer.To(false),
			Lifecycles: expandBackupPolicyCosmosdbAccountLifeCycle(item.Duration),
		})
	}

	return results
}

func expandBackupPolicyCosmosdbAccountDefaultRetentionRule(duration string) basebackuppolicyresources.BasePolicyRule {
	return basebackuppolicyresources.AzureRetentionRule{
		Name:       "Default",
		IsDefault:  pointer.To(true),
		Lifecycles: expandBackupPolicyCosmosdbAccountLifeCycle(duration),
	}
}

func expandBackupPolicyCosmosdbAccountBackupRules(fullBackupSchedule string, incrementalBackupEnabled bool, taggingCriteria []basebackuppolicyresources.TaggingCriteria) ([]basebackuppolicyresources.BasePolicyRule, error) {
	results := []basebackuppolicyresources.BasePolicyRule{
		basebackuppolicyresources.AzureBackupRule{
			Name: "BackupWeekly",
			DataStore: basebackuppolicyresources.DataStoreInfoBase{
				DataStoreType: basebackuppolicyresources.DataStoreTypesVaultStore,
				ObjectType:    "DataStoreInfoBase",
			},
			BackupParameters: basebackuppolicyresources.AzureBackupParams{
				BackupType: "Full",
			},
			Trigger: basebackuppolicyresources.ScheduleBasedTriggerContext{
				Schedule: basebackuppolicyresources.BackupSchedule{
					RepeatingTimeIntervals: []string{fullBackupSchedule},
				},
				TaggingCriteria: taggingCriteria,
			},
		},
	}

	if !incrementalBackupEnabled {
		return results, nil
	}

	incrementalBackupSchedule, err := generateBackupPolicyCosmosdbAccountIncrementalSchedules(fullBackupSchedule)
	if err != nil {
		return nil, err
	}

	results = append(results, basebackuppolicyresources.AzureBackupRule{
		Name: "BackupWeeklyIncremental",
		DataStore: basebackuppolicyresources.DataStoreInfoBase{
			DataStoreType: basebackuppolicyresources.DataStoreTypesVaultStore,
			ObjectType:    "DataStoreInfoBase",
		},
		BackupParameters: basebackuppolicyresources.AzureBackupParams{
			BackupType: "Incremental",
		},
		Trigger: basebackuppolicyresources.ScheduleBasedTriggerContext{
			Schedule: basebackuppolicyresources.BackupSchedule{
				RepeatingTimeIntervals: incrementalBackupSchedule,
			},
			TaggingCriteria: []basebackuppolicyresources.TaggingCriteria{
				expandBackupPolicyCosmosdbAccountDefaultTaggingCriteria(),
			},
		},
	})

	return results, nil
}

func generateBackupPolicyCosmosdbAccountIncrementalSchedules(fullBackupSchedule string) ([]string, error) {
	const (
		recurrencePrefix = "R/"
		recurrenceSuffix = "/P1W"
		timeFormat       = "2006-01-02T15:04:05-07:00"
	)

	timestamp := strings.TrimSuffix(strings.TrimPrefix(fullBackupSchedule, recurrencePrefix), recurrenceSuffix)
	fullBackupTime, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		return nil, fmt.Errorf("parsing timestamp in `backup_schedule` value `%s`: %+v", fullBackupSchedule, err)
	}

	results := make([]string, 0, 6)
	for day := 1; day <= 6; day++ {
		incrementalBackupTime := fullBackupTime.Add(time.Duration(day) * 24 * time.Hour)
		results = append(results, fmt.Sprintf("%s%s%s", recurrencePrefix, incrementalBackupTime.Format(timeFormat), recurrenceSuffix))
	}

	return results, nil
}

func expandBackupPolicyCosmosdbAccountLifeCycle(duration string) []basebackuppolicyresources.SourceLifeCycle {
	// NOTE: currently only `VaultStore` is supported by the service team. When other options are supported
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

func expandBackupPolicyCosmosdbAccountTaggingCriteria(input []BackupPolicyCosmosdbAccountRetentionRule) []basebackuppolicyresources.TaggingCriteria {
	results := []basebackuppolicyresources.TaggingCriteria{
		expandBackupPolicyCosmosdbAccountDefaultTaggingCriteria(),
	}

	for i, item := range input {
		result := basebackuppolicyresources.TaggingCriteria{
			Criteria:        expandBackupPolicyCosmosdbAccountRetentionRuleCriteria(item),
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

func expandBackupPolicyCosmosdbAccountDefaultTaggingCriteria() basebackuppolicyresources.TaggingCriteria {
	return basebackuppolicyresources.TaggingCriteria{
		IsDefault:       true,
		TaggingPriority: 99,
		TagInfo: basebackuppolicyresources.RetentionTag{
			Id:      pointer.To("Default_"),
			TagName: "Default",
		},
	}
}

func expandBackupPolicyCosmosdbAccountRetentionRuleCriteria(input BackupPolicyCosmosdbAccountRetentionRule) *[]basebackuppolicyresources.BackupCriteria {
	var absoluteCriteria []basebackuppolicyresources.AbsoluteMarker
	if len(input.BackupOccurrence) > 0 {
		absoluteCriteria = []basebackuppolicyresources.AbsoluteMarker{basebackuppolicyresources.AbsoluteMarker(input.BackupOccurrence)}
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

	if len(absoluteCriteria) == 0 && len(daysOfWeek) == 0 && len(monthsOfYear) == 0 && len(weeksOfMonth) == 0 {
		return nil
	}

	return &[]basebackuppolicyresources.BackupCriteria{
		basebackuppolicyresources.ScheduleBasedBackupCriteria{
			AbsoluteCriteria: pointer.To(absoluteCriteria),
			DaysOfTheWeek:    pointer.To(daysOfWeek),
			MonthsOfYear:     pointer.To(monthsOfYear),
			WeeksOfTheMonth:  pointer.To(weeksOfMonth),
		},
	}
}
