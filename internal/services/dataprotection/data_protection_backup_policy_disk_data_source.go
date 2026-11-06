// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dataprotection

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-07-01/basebackuppolicyresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type BackupPolicyDiskDataSourceModel struct {
	Name                         string                                    `tfschema:"name"`
	VaultId                      string                                    `tfschema:"vault_id"`
	BackupRepeatingTimeIntervals []string                                  `tfschema:"backup_repeating_time_intervals"`
	DefaultRetentionDuration     string                                    `tfschema:"default_retention_duration"`
	RetentionRule                []BackupPolicyDiskDataSourceRetentionRule `tfschema:"retention_rule"`
	TimeZone                     string                                    `tfschema:"time_zone"`
}

type BackupPolicyDiskDataSourceRetentionRule struct {
	Name     string                               `tfschema:"name"`
	Duration string                               `tfschema:"duration"`
	Criteria []BackupPolicyDiskDataSourceCriteria `tfschema:"criteria"`
	Priority int64                                `tfschema:"priority"`
}

type BackupPolicyDiskDataSourceCriteria struct {
	AbsoluteCriteria string `tfschema:"absolute_criteria"`
}

type DataProtectionBackupPolicyDiskDataSource struct{}

var _ sdk.DataSource = DataProtectionBackupPolicyDiskDataSource{}

func (r DataProtectionBackupPolicyDiskDataSource) ResourceType() string {
	return "azurerm_data_protection_backup_policy_disk"
}

func (r DataProtectionBackupPolicyDiskDataSource) ModelObject() interface{} {
	return &BackupPolicyDiskDataSourceModel{}
}

func (r DataProtectionBackupPolicyDiskDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"vault_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: basebackuppolicyresources.ValidateBackupVaultID,
		},
	}
}

func (r DataProtectionBackupPolicyDiskDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"backup_repeating_time_intervals": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"default_retention_duration": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"retention_rule": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"duration": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"criteria": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"absolute_criteria": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
							},
						},
					},

					"priority": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},
				},
			},
		},

		"time_zone": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r DataProtectionBackupPolicyDiskDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DataProtection.BackupPolicyClient

			var model BackupPolicyDiskDataSourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			vaultId, err := basebackuppolicyresources.ParseBackupVaultID(model.VaultId)
			if err != nil {
				return err
			}

			id := basebackuppolicyresources.NewBackupPolicyID(vaultId.SubscriptionId, vaultId.ResourceGroupName, vaultId.BackupVaultName, model.Name)

			resp, err := client.BackupPoliciesGet(ctx, id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			state := BackupPolicyDiskDataSourceModel{
				Name:    id.BackupPolicyName,
				VaultId: vaultId.ID(),
			}

			if respModel := resp.Model; respModel != nil {
				if properties, ok := respModel.Properties.(basebackuppolicyresources.BackupPolicy); ok {
					state.BackupRepeatingTimeIntervals = flattenBackupPolicyBackupRepeatingTimeIntervals(properties.PolicyRules)
					state.TimeZone = flattenBackupPolicyBackupTimeZone(properties.PolicyRules)
					state.DefaultRetentionDuration = flattenBackupPolicyDiskDataSourceDefaultRetentionDuration(properties.PolicyRules)
					state.RetentionRule = flattenBackupPolicyDiskDataSourceRetentionRules(properties.PolicyRules)
				}
			}

			metadata.SetID(id)
			return metadata.Encode(&state)
		},
	}
}

func flattenBackupPolicyDiskDataSourceDefaultRetentionDuration(input []basebackuppolicyresources.BasePolicyRule) string {
	for _, item := range input {
		if retentionRule, ok := item.(basebackuppolicyresources.AzureRetentionRule); ok && retentionRule.IsDefault != nil && *retentionRule.IsDefault {
			if len(retentionRule.Lifecycles) > 0 {
				if deleteOption, ok := (retentionRule.Lifecycles)[0].DeleteAfter.(basebackuppolicyresources.AbsoluteDeleteOption); ok {
					return deleteOption.Duration
				}
			}
		}
	}
	return ""
}

func flattenBackupPolicyDiskDataSourceRetentionRules(input []basebackuppolicyresources.BasePolicyRule) []BackupPolicyDiskDataSourceRetentionRule {
	results := make([]BackupPolicyDiskDataSourceRetentionRule, 0)
	if len(input) == 0 {
		return results
	}

	var taggingCriterias []basebackuppolicyresources.TaggingCriteria
	for _, item := range input {
		if backupRule, ok := item.(basebackuppolicyresources.AzureBackupRule); ok {
			if trigger, ok := backupRule.Trigger.(basebackuppolicyresources.ScheduleBasedTriggerContext); ok {
				if trigger.TaggingCriteria != nil {
					taggingCriterias = trigger.TaggingCriteria
				}
			}
		}
	}

	for _, item := range input {
		if retentionRule, ok := item.(basebackuppolicyresources.AzureRetentionRule); ok && (retentionRule.IsDefault == nil || !*retentionRule.IsDefault) {
			name := retentionRule.Name
			var taggingPriority int64
			var taggingCriteria []BackupPolicyDiskDataSourceCriteria
			for _, criteria := range taggingCriterias {
				if strings.EqualFold(criteria.TagInfo.TagName, name) {
					taggingPriority = criteria.TaggingPriority
					taggingCriteria = flattenBackupPolicyDiskDataSourceCriteria(criteria.Criteria)
					break
				}
			}
			var duration string
			if len(retentionRule.Lifecycles) > 0 {
				if deleteOption, ok := (retentionRule.Lifecycles)[0].DeleteAfter.(basebackuppolicyresources.AbsoluteDeleteOption); ok {
					duration = deleteOption.Duration
				}
			}
			results = append(results, BackupPolicyDiskDataSourceRetentionRule{
				Name:     name,
				Priority: taggingPriority,
				Criteria: taggingCriteria,
				Duration: duration,
			})
		}
	}
	return results
}

func flattenBackupPolicyDiskDataSourceCriteria(input *[]basebackuppolicyresources.BackupCriteria) []BackupPolicyDiskDataSourceCriteria {
	results := make([]BackupPolicyDiskDataSourceCriteria, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		if criteria, ok := item.(basebackuppolicyresources.ScheduleBasedBackupCriteria); ok {
			var absoluteCriteria string
			if criteria.AbsoluteCriteria != nil && len(*criteria.AbsoluteCriteria) > 0 {
				absoluteCriteria = string((*criteria.AbsoluteCriteria)[0])
			}

			results = append(results, BackupPolicyDiskDataSourceCriteria{
				AbsoluteCriteria: absoluteCriteria,
			})
		}
	}
	return results
}
