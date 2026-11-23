// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_protection_backup_policy_elastic_san_volume

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

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := basebackuppolicyresources.ParseBackupPolicyIDInsensitively(metadata.ResourceData.Id())
			if err != nil {
				return err
			}
			resp, err := metadata.Client.DataProtection.BackupPolicyClient20260601.BackupPoliciesGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}
			return r.flatten(metadata, id, resp.Model)
		},
	}
}

func (r Resource) flatten(metadata sdk.ResourceMetaData, id *basebackuppolicyresources.BackupPolicyId, model *basebackuppolicyresources.BaseBackupPolicyResource) error {
	if model == nil {
		return fmt.Errorf("retrieving %s: model was nil", *id)
	}
	if model.Properties == nil {
		return fmt.Errorf("retrieving %s: properties were nil", *id)
	}
	if !backupPolicySupportsElasticSanVolumes(model.Properties) {
		return fmt.Errorf("%s is not a backup policy for `%s`", *id, elasticSanVolumeDatasourceType)
	}

	vaultId := basebackuppolicyresources.NewBackupVaultID(id.SubscriptionId, id.ResourceGroupName, id.BackupVaultName)
	state := BackupPolicyElasticSanVolumeModel{
		Name:                        id.BackupPolicyName,
		DataProtectionBackupVaultId: vaultId.ID(),
	}
	properties := model.Properties.(basebackuppolicyresources.BackupPolicy)
	if err := flattenBackupPolicyElasticSanVolumePolicyRules(properties.PolicyRules, &state); err != nil {
		return fmt.Errorf("flattening policy rules for %s: %+v", *id, err)
	}

	metadata.SetID(id)
	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
		return err
	}
	return metadata.Encode(&state)
}

func flattenBackupPolicyElasticSanVolumePolicyRules(input []basebackuppolicyresources.BasePolicyRule, state *BackupPolicyElasticSanVolumeModel) error {
	var taggingCriteria []basebackuppolicyresources.TaggingCriteria
	for _, item := range input {
		switch rule := item.(type) {
		case basebackuppolicyresources.AzureBackupRule:
			trigger, ok := rule.Trigger.(basebackuppolicyresources.ScheduleBasedTriggerContext)
			if !ok {
				return fmt.Errorf("unsupported trigger for backup rule `%s`", rule.Name)
			}
			if len(trigger.Schedule.RepeatingTimeIntervals) != 1 {
				return fmt.Errorf("backup rule `%s` must have exactly one repeating time interval for `daily_backup_schedule` or `weekly_backup_schedule`", rule.Name)
			}
			backupSchedule := trigger.Schedule.RepeatingTimeIntervals[0]
			switch {
			case strings.HasSuffix(backupSchedule, "/P1D"):
				state.DailyBackupSchedule = backupSchedule
			case strings.HasSuffix(backupSchedule, "/P1W"):
				state.WeeklyBackupSchedule = backupSchedule
			default:
				return fmt.Errorf("unsupported repeating time interval `%s` for backup rule `%s`", backupSchedule, rule.Name)
			}
			taggingCriteria = append(taggingCriteria, trigger.TaggingCriteria...)
		case basebackuppolicyresources.AzureRetentionRule:
			if len(rule.Lifecycles) != 1 {
				return fmt.Errorf("retention rule `%s` must have exactly one lifecycle", rule.Name)
			}
			deleteAfter, ok := rule.Lifecycles[0].DeleteAfter.(basebackuppolicyresources.AbsoluteDeleteOption)
			if !ok {
				return fmt.Errorf("unsupported delete option for retention rule `%s`", rule.Name)
			}
			switch {
			case pointer.From(rule.IsDefault):
				state.DefaultRetentionDuration = deleteAfter.Duration
			case strings.EqualFold(rule.Name, "Daily"):
				state.DailyRetentionDuration = deleteAfter.Duration
			case strings.EqualFold(rule.Name, "Monthly"):
				state.MonthlyRetentionDuration = deleteAfter.Duration
			case strings.EqualFold(rule.Name, "Weekly"):
				state.WeeklyRetentionDuration = deleteAfter.Duration
			default:
				return fmt.Errorf("unsupported retention rule `%s`", rule.Name)
			}
		default:
			return fmt.Errorf("unsupported policy rule type `%T`", item)
		}
	}

	return nil
}
