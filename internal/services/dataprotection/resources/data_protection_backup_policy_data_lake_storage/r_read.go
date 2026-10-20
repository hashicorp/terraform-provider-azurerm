// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package data_protection_backup_policy_data_lake_storage

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-07-01/basebackuppolicyresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DataProtection.BackupPolicyClient

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

			vaultId := basebackuppolicyresources.NewBackupVaultID(id.SubscriptionId, id.ResourceGroupName, id.BackupVaultName)
			state := BackupPolicyDataLakeStorageModel{
				Name:                        id.BackupPolicyName,
				DataProtectionBackupVaultId: vaultId.ID(),
			}

			if model := resp.Model; model != nil {
				if properties, ok := model.Properties.(basebackuppolicyresources.BackupPolicy); ok {
					state.DefaultRetentionDuration, state.RetentionRules, state.BackupSchedule, state.TimeZone = flattenBackupPolicyDataLakeStoragePolicyRules(properties.PolicyRules)
				}
			}

			if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
				return err
			}

			return metadata.Encode(&state)
		},
	}
}

func flattenBackupPolicyDataLakeStoragePolicyRules(input []basebackuppolicyresources.BasePolicyRule) (string, []BackupPolicyDataLakeStorageRetentionRule, []string, string) {
	var taggingCriteria []basebackuppolicyresources.TaggingCriteria
	var nonDefaultRetentionRules []basebackuppolicyresources.AzureRetentionRule
	var backupSchedule []string
	var timeZone string
	var defaultRetentionDuration string
	retentionRules := make([]BackupPolicyDataLakeStorageRetentionRule, 0)

	for _, item := range input {
		switch rule := item.(type) {
		case basebackuppolicyresources.AzureBackupRule:
			if trigger, ok := rule.Trigger.(basebackuppolicyresources.ScheduleBasedTriggerContext); ok {
				backupSchedule = trigger.Schedule.RepeatingTimeIntervals
				timeZone = pointer.From(trigger.Schedule.TimeZone)
				taggingCriteria = trigger.TaggingCriteria
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
		result := BackupPolicyDataLakeStorageRetentionRule{
			Name: rule.Name,
		}

		for _, criteria := range taggingCriteria {
			if strings.EqualFold(criteria.TagInfo.TagName, rule.Name) {
				flattenBackupPolicyDataLakeStorageCriteriaIntoRule(criteria.Criteria, &result)
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

	return defaultRetentionDuration, retentionRules, backupSchedule, timeZone
}

func flattenBackupPolicyDataLakeStorageCriteriaIntoRule(input *[]basebackuppolicyresources.BackupCriteria, rule *BackupPolicyDataLakeStorageRetentionRule) {
	if input == nil {
		return
	}

	for _, item := range pointer.From(input) {
		if criteria, ok := item.(basebackuppolicyresources.ScheduleBasedBackupCriteria); ok {
			if criteria.AbsoluteCriteria != nil && len(pointer.From(criteria.AbsoluteCriteria)) > 0 {
				rule.AbsoluteCriteria = string(pointer.From(criteria.AbsoluteCriteria)[0])
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

			if criteria.ScheduleTimes != nil {
				scheduleTimes := make([]string, 0)
				scheduleTimes = append(scheduleTimes, pointer.From(criteria.ScheduleTimes)...)
				rule.ScheduledBackupTimes = scheduleTimes
			}
		}
	}
}

type BackupPolicyDataLakeStorageRetentionRule struct {
	Name                 string   `tfschema:"name"`
	Duration             string   `tfschema:"duration"`
	AbsoluteCriteria     string   `tfschema:"absolute_criteria"`
	DaysOfWeek           []string `tfschema:"days_of_week"`
	MonthsOfYear         []string `tfschema:"months_of_year"`
	ScheduledBackupTimes []string `tfschema:"scheduled_backup_times"`
	WeeksOfMonth         []string `tfschema:"weeks_of_month"`
}
