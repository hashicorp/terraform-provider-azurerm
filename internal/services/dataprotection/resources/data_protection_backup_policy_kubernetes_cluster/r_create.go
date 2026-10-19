// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package data_protection_backup_policy_kubernetes_cluster

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
			var model BackupPolicyKubernatesClusterModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.DataProtection.BackupPolicyClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			id := basebackuppolicyresources.NewBackupPolicyID(subscriptionId, model.ResourceGroupName, model.VaultName, model.Name)

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

			taggingCriteria, err := expandBackupPolicyKubernetesClusterTaggingCriteriaArray(model.RetentionRule)
			if err != nil {
				return err
			}

			policyRules := make([]basebackuppolicyresources.BasePolicyRule, 0)
			policyRules = append(policyRules, expandBackupPolicyKubernetesClusterAzureBackupRuleArray(model.BackupRepeatingTimeIntervals, model.TimeZone, taggingCriteria)...)
			if v := expandBackupPolicyKubernetesClusterDefaultRetentionRule(model.DefaultRetentionRule); v != nil {
				policyRules = append(policyRules, pointer.From(v))
			}
			policyRules = append(policyRules, expandBackupPolicyKubernetesClusterAzureRetentionRules(model.RetentionRule)...)

			parameters := basebackuppolicyresources.BaseBackupPolicyResource{
				Properties: &basebackuppolicyresources.BackupPolicy{
					PolicyRules:     policyRules,
					DatasourceTypes: []string{"Microsoft.ContainerService/managedClusters"},
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

func expandBackupPolicyKubernetesClusterAzureBackupRuleArray(input []string, timeZone string, taggingCriteria *[]basebackuppolicyresources.TaggingCriteria) []basebackuppolicyresources.BasePolicyRule {
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
				RepeatingTimeIntervals: input,
				TimeZone:               pointer.To(timeZone),
			},
			TaggingCriteria: *taggingCriteria,
		},
	})

	return results
}

func expandBackupPolicyKubernetesClusterDefaultRetentionRule(input []DefaultRetentionRule) *basebackuppolicyresources.AzureRetentionRule {
	if len(input) == 0 {
		return nil
	}
	return &basebackuppolicyresources.AzureRetentionRule{
		Name:       "Default",
		IsDefault:  pointer.To(true),
		Lifecycles: expandBackupPolicyKubernetesClusterLifeCycle(input[0].LifeCycle),
	}
}

func expandBackupPolicyKubernetesClusterAzureRetentionRules(input []RetentionRule) []basebackuppolicyresources.BasePolicyRule {
	results := make([]basebackuppolicyresources.BasePolicyRule, 0)
	for _, item := range input {
		lifeCycle := expandBackupPolicyKubernetesClusterLifeCycle(item.LifeCycle)

		results = append(results, basebackuppolicyresources.AzureRetentionRule{
			Name:       item.Name,
			IsDefault:  pointer.To(false),
			Lifecycles: lifeCycle,
		})
	}
	return results
}

func expandBackupPolicyKubernetesClusterLifeCycle(input []LifeCycle) []basebackuppolicyresources.SourceLifeCycle {
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

func expandBackupPolicyKubernetesClusterTaggingCriteriaArray(input []RetentionRule) (*[]basebackuppolicyresources.TaggingCriteria, error) {
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
			TaggingPriority: item.Priority,
			TagInfo: basebackuppolicyresources.RetentionTag{
				Id:      pointer.To(item.Name + "_"),
				TagName: item.Name,
			},
		}

		criteria, err := expandBackupPolicyKubernetesClusterCriteriaArray(item.Criteria)
		if err != nil {
			return nil, err
		}
		result.Criteria = criteria
		results = append(results, result)
	}
	return &results, nil
}

func expandBackupPolicyKubernetesClusterCriteriaArray(input []Criteria) (*[]basebackuppolicyresources.BackupCriteria, error) {
	if len(input) == 0 {
		return nil, fmt.Errorf("criteria is a required field, cannot leave blank")
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

		results = append(results, basebackuppolicyresources.ScheduleBasedBackupCriteria{
			AbsoluteCriteria: &absoluteCriteria,
			DaysOfTheWeek:    &daysOfWeek,
			MonthsOfYear:     &monthsOfYear,
			ScheduleTimes:    pointer.To(item.ScheduledBackupTimes),
			WeeksOfTheMonth:  &weeksOfMonth,
		})
	}
	return &results, nil
}
