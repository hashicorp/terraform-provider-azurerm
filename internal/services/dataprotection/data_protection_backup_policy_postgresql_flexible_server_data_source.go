// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dataprotection

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-09-01/basebackuppolicyresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type BackupPolicyPostgreSQLFlexibleServerDataSourceModel struct {
	Name                         string                                                     `tfschema:"name"`
	VaultId                      string                                                     `tfschema:"vault_id"`
	BackupRepeatingTimeIntervals []string                                                   `tfschema:"backup_repeating_time_intervals"`
	DefaultRetentionRule         []BackupPolicyPostgreSQLFlexibleServerDefaultRetentionRule `tfschema:"default_retention_rule"`
	RetentionRule                []BackupPolicyPostgreSQLFlexibleServerRetentionRule        `tfschema:"retention_rule"`
	TimeZone                     string                                                     `tfschema:"time_zone"`
}

type DataProtectionBackupPolicyPostgreSQLFlexibleServerDataSource struct{}

var _ sdk.DataSource = DataProtectionBackupPolicyPostgreSQLFlexibleServerDataSource{}

func (r DataProtectionBackupPolicyPostgreSQLFlexibleServerDataSource) ResourceType() string {
	return "azurerm_data_protection_backup_policy_postgresql_flexible_server"
}

func (r DataProtectionBackupPolicyPostgreSQLFlexibleServerDataSource) ModelObject() interface{} {
	return &BackupPolicyPostgreSQLFlexibleServerDataSourceModel{}
}

func (r DataProtectionBackupPolicyPostgreSQLFlexibleServerDataSource) Arguments() map[string]*pluginsdk.Schema {
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

func (r DataProtectionBackupPolicyPostgreSQLFlexibleServerDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"backup_repeating_time_intervals": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"default_retention_rule": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"life_cycle": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"data_store_type": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},

								"duration": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
							},
						},
					},
				},
			},
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

					"criteria": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"absolute_criteria": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},

								"days_of_week": {
									Type:     pluginsdk.TypeSet,
									Computed: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},

								"months_of_year": {
									Type:     pluginsdk.TypeSet,
									Computed: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},

								"scheduled_backup_times": {
									Type:     pluginsdk.TypeSet,
									Computed: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},

								"weeks_of_month": {
									Type:     pluginsdk.TypeSet,
									Computed: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
							},
						},
					},

					"life_cycle": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"data_store_type": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},

								"duration": {
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

func (r DataProtectionBackupPolicyPostgreSQLFlexibleServerDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DataProtection.BackupPolicyClient

			var model BackupPolicyPostgreSQLFlexibleServerDataSourceModel
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

			state := BackupPolicyPostgreSQLFlexibleServerDataSourceModel{
				Name:    id.BackupPolicyName,
				VaultId: vaultId.ID(),
			}

			if respModel := resp.Model; respModel != nil {
				if properties, ok := respModel.Properties.(basebackuppolicyresources.BackupPolicy); ok {
					state.BackupRepeatingTimeIntervals = flattenBackupPolicyPostgreSQLFlexibleServerBackupRules(properties.PolicyRules)
					state.TimeZone = flattenBackupPolicyPostgreSQLFlexibleServerBackupTimeZone(properties.PolicyRules)
					state.DefaultRetentionRule = flattenBackupPolicyPostgreSQLFlexibleServerDefaultRetentionRule(properties.PolicyRules)
					state.RetentionRule = flattenBackupPolicyPostgreSQLFlexibleServerRetentionRules(properties.PolicyRules)
				}
			}

			metadata.SetID(id)
			return metadata.Encode(&state)
		},
	}
}
