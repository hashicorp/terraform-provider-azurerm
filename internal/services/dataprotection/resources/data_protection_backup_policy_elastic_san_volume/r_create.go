// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_protection_backup_policy_elastic_san_volume

import (
	"context"
	"fmt"
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

			var model BackupPolicyElasticSanVolumeModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			vaultId, err := basebackuppolicyresources.ParseBackupVaultIDInsensitively(model.DataProtectionBackupVaultId)
			if err != nil {
				return err
			}
			id := basebackuppolicyresources.NewBackupPolicyID(vaultId.SubscriptionId, vaultId.ResourceGroupName, vaultId.BackupVaultName, model.Name)

			existing, err := client.BackupPoliciesGet(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
				if existing.Model == nil {
					return fmt.Errorf("retrieving %s: model was nil", id)
				}
				if !backupPolicySupportsElasticSanVolumes(existing.Model.Properties) {
					return fmt.Errorf("%s is not a backup policy for `%s`", id, elasticSanVolumeDatasourceType)
				}
			}

			parameters := basebackuppolicyresources.BaseBackupPolicyResource{
				Properties: basebackuppolicyresources.BackupPolicy{
					DatasourceTypes: []string{elasticSanVolumeDatasourceType},
					PolicyRules:     expandBackupPolicyElasticSanVolumePolicyRules(model),
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

func expandBackupPolicyElasticSanVolumePolicyRules(input BackupPolicyElasticSanVolumeModel) []basebackuppolicyresources.BasePolicyRule {
	taggingCriteria := []basebackuppolicyresources.TaggingCriteria{{
		IsDefault:       true,
		TaggingPriority: 99,
		TagInfo: basebackuppolicyresources.RetentionTag{
			Id:      pointer.To("Default_"),
			TagName: "Default",
		},
	}}
	rules := []basebackuppolicyresources.BasePolicyRule{
		expandBackupPolicyElasticSanVolumeRetentionRule("Default", input.DefaultRetentionDuration, true),
	}

	if input.DailyRetentionDuration != "" {
		taggingCriteria = append(taggingCriteria, expandBackupPolicyElasticSanVolumeTaggingCriteria("Daily", 3, basebackuppolicyresources.ScheduleBasedBackupCriteria{
			AbsoluteCriteria: pointer.To([]basebackuppolicyresources.AbsoluteMarker{basebackuppolicyresources.AbsoluteMarkerFirstOfDay}),
		}))
		rules = append(rules, expandBackupPolicyElasticSanVolumeRetentionRule("Daily", input.DailyRetentionDuration, false))
	}
	if input.WeeklyRetentionDuration != "" {
		taggingCriteria = append(taggingCriteria, expandBackupPolicyElasticSanVolumeTaggingCriteria("Weekly", 2, basebackuppolicyresources.ScheduleBasedBackupCriteria{
			AbsoluteCriteria: pointer.To([]basebackuppolicyresources.AbsoluteMarker{basebackuppolicyresources.AbsoluteMarkerFirstOfWeek}),
		}))
		rules = append(rules, expandBackupPolicyElasticSanVolumeRetentionRule("Weekly", input.WeeklyRetentionDuration, false))
	}
	if input.MonthlyRetentionDuration != "" {
		taggingCriteria = append(taggingCriteria, expandBackupPolicyElasticSanVolumeTaggingCriteria("Monthly", 1, basebackuppolicyresources.ScheduleBasedBackupCriteria{
			AbsoluteCriteria: pointer.To([]basebackuppolicyresources.AbsoluteMarker{basebackuppolicyresources.AbsoluteMarkerFirstOfMonth}),
		}))
		rules = append(rules, expandBackupPolicyElasticSanVolumeRetentionRule("Monthly", input.MonthlyRetentionDuration, false))
	}

	backupSchedule := input.DailyBackupSchedule
	backupRuleName := "BackupDaily"
	if backupSchedule == "" {
		backupSchedule = input.WeeklyBackupSchedule
		backupRuleName = "BackupWeekly"
	}
	schedule := basebackuppolicyresources.BackupSchedule{
		RepeatingTimeIntervals: []string{backupSchedule},
	}
	rules = append(rules, basebackuppolicyresources.AzureBackupRule{
		Name: backupRuleName,
		BackupParameters: basebackuppolicyresources.AzureBackupParams{
			BackupType: "Incremental",
		},
		DataStore: basebackuppolicyresources.DataStoreInfoBase{
			DataStoreType: basebackuppolicyresources.DataStoreTypesOperationalStore,
			ObjectType:    "DataStoreInfoBase",
		},
		Trigger: basebackuppolicyresources.ScheduleBasedTriggerContext{
			Schedule:        schedule,
			TaggingCriteria: taggingCriteria,
		},
	})
	return rules
}

func expandBackupPolicyElasticSanVolumeRetentionRule(name, duration string, isDefault bool) basebackuppolicyresources.AzureRetentionRule {
	return basebackuppolicyresources.AzureRetentionRule{
		Name:      name,
		IsDefault: pointer.To(isDefault),
		Lifecycles: []basebackuppolicyresources.SourceLifeCycle{{
			DeleteAfter: basebackuppolicyresources.AbsoluteDeleteOption{
				Duration: duration,
			},
			SourceDataStore: basebackuppolicyresources.DataStoreInfoBase{
				DataStoreType: basebackuppolicyresources.DataStoreTypesOperationalStore,
				ObjectType:    "DataStoreInfoBase",
			},
		}},
	}
}

func expandBackupPolicyElasticSanVolumeTaggingCriteria(name string, priority int64, criteria basebackuppolicyresources.ScheduleBasedBackupCriteria) basebackuppolicyresources.TaggingCriteria {
	return basebackuppolicyresources.TaggingCriteria{
		Criteria:        pointer.To([]basebackuppolicyresources.BackupCriteria{criteria}),
		TaggingPriority: priority,
		TagInfo: basebackuppolicyresources.RetentionTag{
			Id:      pointer.To(name + "_"),
			TagName: name,
		},
	}
}
