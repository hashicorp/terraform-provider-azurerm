// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_protection_backup_policy_cosmosdb_account

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2026-06-01/basebackuppolicyresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DataProtection.BackupPolicyClient20260601
			id, err := basebackuppolicyresources.ParseBackupPolicyID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.BackupPoliciesGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if resp.Model == nil {
				return fmt.Errorf("retrieving %s: model was nil", *id)
			}
			properties, ok := resp.Model.Properties.(basebackuppolicyresources.BackupPolicy)
			if !ok || !backupPolicySupportsCosmosdbAccounts(properties) {
				return fmt.Errorf("%s is not a backup policy for `Microsoft.DocumentDB/databaseAccounts`", *id)
			}

			return r.flatten(metadata, id, resp.Model)
		},
	}
}

func (r Resource) flatten(metadata sdk.ResourceMetaData, id *basebackuppolicyresources.BackupPolicyId, model *basebackuppolicyresources.BaseBackupPolicyResource) error {
	vaultId := basebackuppolicyresources.NewBackupVaultID(id.SubscriptionId, id.ResourceGroupName, id.BackupVaultName)
	state := BackupPolicyCosmosdbAccountModel{
		Name:                        id.BackupPolicyName,
		DataProtectionBackupVaultId: vaultId.ID(),
	}

	if model != nil {
		if properties, ok := model.Properties.(basebackuppolicyresources.BackupPolicy); ok {
			var retentionRules []BackupPolicyCosmosdbAccountRetentionRule
			state.DefaultRetentionDuration, retentionRules, state.BackupSchedule, state.IncrementalBackupSchedules, state.BackupFrequency = flattenBackupPolicyCosmosdbAccountPolicyRules(properties.PolicyRules)

			if err := flattenBackupPolicyCosmosdbAccountRetentionRules(retentionRules, &state); err != nil {
				return fmt.Errorf("flattening retention rules for %s: %+v", *id, err)
			}
		}
	}

	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
		return err
	}
	return metadata.Encode(&state)
}

func flattenBackupPolicyCosmosdbAccountRetentionRules(input []BackupPolicyCosmosdbAccountRetentionRule, state *BackupPolicyCosmosdbAccountModel) error {
	for _, rule := range input {
		switch {
		case rule.BackupOccurrence == string(basebackuppolicyresources.AbsoluteMarkerFirstOfYear) || len(rule.MonthsOfYear) > 0:
			if len(state.YearlyRetentionRule) > 0 {
				return errors.New("multiple yearly retention rules cannot be represented by `yearly_retention_rule`")
			}
			state.YearlyRetentionRule = []BackupPolicyCosmosdbAccountYearlyRetentionRule{{
				Duration:     rule.Duration,
				DaysOfWeek:   rule.DaysOfWeek,
				WeeksOfMonth: rule.WeeksOfMonth,
				MonthsOfYear: rule.MonthsOfYear,
			}}
		case rule.BackupOccurrence == string(basebackuppolicyresources.AbsoluteMarkerFirstOfMonth) || len(rule.WeeksOfMonth) > 0:
			if len(state.MonthlyRetentionRule) > 0 {
				return errors.New("multiple monthly retention rules cannot be represented by `monthly_retention_rule`")
			}
			state.MonthlyRetentionRule = []BackupPolicyCosmosdbAccountMonthlyRetentionRule{{
				Duration:     rule.Duration,
				DaysOfWeek:   rule.DaysOfWeek,
				WeeksOfMonth: rule.WeeksOfMonth,
			}}
		default:
			if len(state.WeeklyRetentionRule) > 0 {
				return errors.New("multiple weekly retention rules cannot be represented by `weekly_retention_rule`")
			}
			state.WeeklyRetentionRule = []BackupPolicyCosmosdbAccountWeeklyRetentionRule{{
				Duration:   rule.Duration,
				DaysOfWeek: rule.DaysOfWeek,
			}}
		}
	}
	return nil
}

func flattenBackupPolicyCosmosdbAccountPolicyRules(input []basebackuppolicyresources.BasePolicyRule) (string, []BackupPolicyCosmosdbAccountRetentionRule, string, []string, string) {
	var taggingCriteria []basebackuppolicyresources.TaggingCriteria
	var nonDefaultRetentionRules []basebackuppolicyresources.AzureRetentionRule
	backupFrequency := backupPolicyCosmosdbAccountFrequencyWeekly
	var fullBackupSchedule string
	var defaultRetentionDuration string

	incrementalBackupSchedule := make([]string, 0)
	retentionRules := make([]BackupPolicyCosmosdbAccountRetentionRule, 0)

	for _, item := range input {
		switch rule := item.(type) {
		case basebackuppolicyresources.AzureBackupRule:
			if trigger, ok := rule.Trigger.(basebackuppolicyresources.ScheduleBasedTriggerContext); ok {
				if parameters, ok := rule.BackupParameters.(basebackuppolicyresources.AzureBackupParams); ok {
					switch {
					case strings.EqualFold(parameters.BackupType, "Full"):
						if len(trigger.Schedule.RepeatingTimeIntervals) > 0 {
							fullBackupSchedule = trigger.Schedule.RepeatingTimeIntervals[0]
						}
						taggingCriteria = trigger.TaggingCriteria
					case strings.EqualFold(parameters.BackupType, "Incremental"):
						backupFrequency = backupPolicyCosmosdbAccountFrequencyDaily
						incrementalBackupSchedule = append(incrementalBackupSchedule, trigger.Schedule.RepeatingTimeIntervals...)
					}
				}
			}
		case basebackuppolicyresources.AzureRetentionRule:
			if pointer.From(rule.IsDefault) {
				if v := rule.Lifecycles; len(v) > 0 {
					if deleteOption, ok := v[0].DeleteAfter.(basebackuppolicyresources.AbsoluteDeleteOption); ok {
						defaultRetentionDuration = deleteOption.Duration
					}
				}
			} else {
				nonDefaultRetentionRules = append(nonDefaultRetentionRules, rule)
			}
		}
	}

	for _, rule := range nonDefaultRetentionRules {
		result := BackupPolicyCosmosdbAccountRetentionRule{
			Name: rule.Name,
		}

		for _, criteria := range taggingCriteria {
			if strings.EqualFold(criteria.TagInfo.TagName, rule.Name) {
				flattenBackupPolicyCosmosdbAccountCriteriaIntoRule(criteria.Criteria, &result)
				break
			}
		}

		if v := rule.Lifecycles; len(v) > 0 {
			if deleteOption, ok := v[0].DeleteAfter.(basebackuppolicyresources.AbsoluteDeleteOption); ok {
				result.Duration = deleteOption.Duration
			}
		}

		retentionRules = append(retentionRules, result)
	}

	return defaultRetentionDuration, retentionRules, fullBackupSchedule, incrementalBackupSchedule, backupFrequency
}

func flattenBackupPolicyCosmosdbAccountCriteriaIntoRule(input *[]basebackuppolicyresources.BackupCriteria, rule *BackupPolicyCosmosdbAccountRetentionRule) {
	if input == nil {
		return
	}

	for _, item := range pointer.From(input) {
		if criteria, ok := item.(basebackuppolicyresources.ScheduleBasedBackupCriteria); ok {
			if criteria.AbsoluteCriteria != nil && len(pointer.From(criteria.AbsoluteCriteria)) > 0 {
				rule.BackupOccurrence = string(pointer.From(criteria.AbsoluteCriteria)[0])
			}

			if criteria.DaysOfTheWeek != nil {
				daysOfWeek := make([]string, 0)
				for _, item := range pointer.From(criteria.DaysOfTheWeek) {
					daysOfWeek = append(daysOfWeek, string(item))
				}
				rule.DaysOfWeek = daysOfWeek
			}

			if criteria.MonthsOfYear != nil {
				monthsOfYear := make([]string, 0)
				for _, item := range pointer.From(criteria.MonthsOfYear) {
					monthsOfYear = append(monthsOfYear, string(item))
				}
				rule.MonthsOfYear = monthsOfYear
			}

			if criteria.WeeksOfTheMonth != nil {
				weeksOfMonth := make([]string, 0)
				for _, item := range pointer.From(criteria.WeeksOfTheMonth) {
					weeksOfMonth = append(weeksOfMonth, string(item))
				}
				rule.WeeksOfMonth = weeksOfMonth
			}
		}
	}
}
