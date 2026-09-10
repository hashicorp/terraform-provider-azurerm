// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package dataprotection

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	basebackuppolicyresources20250701 "github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-07-01/basebackuppolicyresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2026-06-01/basebackuppolicyresources"
	azValidate "github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

//go:generate go run ../../tools/generator-tests resourceidentity -resource-name data_protection_backup_policy_cosmosdb_database_account -service-package-name dataprotection -properties "name" -compare-values "subscription_id:vault_id,resource_group_name:vault_id,backup_vault_name:vault_id"

type BackupPolicyCosmosDBDatabaseAccountModel struct {
	Name                         string                                                `tfschema:"name"`
	VaultId                      string                                                `tfschema:"vault_id"`
	BackupRepeatingTimeIntervals []string                                              `tfschema:"backup_repeating_time_intervals"`
	DefaultRetentionRule         []BackupPolicyMySQLFlexibleServerDefaultRetentionRule `tfschema:"default_retention_rule"`
	RetentionRules               []BackupPolicyMySQLFlexibleServerRetentionRule        `tfschema:"retention_rule"`
	TimeZone                     string                                                `tfschema:"time_zone"`
}

type DataProtectionBackupPolicyCosmosDBDatabaseAccountResource struct{}

var (
	_ sdk.Resource             = DataProtectionBackupPolicyCosmosDBDatabaseAccountResource{}
	_ sdk.ResourceWithIdentity = DataProtectionBackupPolicyCosmosDBDatabaseAccountResource{}
)

func (r DataProtectionBackupPolicyCosmosDBDatabaseAccountResource) Identity() resourceids.ResourceId {
	return &basebackuppolicyresources.BackupPolicyId{}
}

func (r DataProtectionBackupPolicyCosmosDBDatabaseAccountResource) ResourceType() string {
	return "azurerm_data_protection_backup_policy_cosmosdb_database_account"
}

func (r DataProtectionBackupPolicyCosmosDBDatabaseAccountResource) ModelObject() interface{} {
	return &BackupPolicyCosmosDBDatabaseAccountModel{}
}

func (r DataProtectionBackupPolicyCosmosDBDatabaseAccountResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return basebackuppolicyresources.ValidateBackupPolicyID
}

func (r DataProtectionBackupPolicyCosmosDBDatabaseAccountResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile("^[a-zA-Z][-a-zA-Z0-9]{2,149}$"),
				"`name` must be 3 - 150 characters long, contain only letters, numbers and hyphens, and cannot start with a number or hyphen",
			),
		},

		"vault_id": commonschema.ResourceIDReferenceRequiredForceNew(pointer.To(basebackuppolicyresources.BackupVaultId{})),

		"backup_repeating_time_intervals": {
			Type:     pluginsdk.TypeList,
			Required: true,
			ForceNew: true,
			MinItems: 1,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: azValidate.ISO8601RepeatingTime,
			},
		},

		"default_retention_rule": {
			Type:     pluginsdk.TypeList,
			Required: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"life_cycle": cosmosDBDatabaseAccountLifeCycleSchema(),
				},
			},
		},

		"retention_rule": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"criteria": {
						Type:     pluginsdk.TypeList,
						Required: true,
						ForceNew: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"absolute_criteria": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ForceNew:     true,
									ValidateFunc: validation.StringInSlice(basebackuppolicyresources.PossibleValuesForAbsoluteMarker(), false),
								},

								"days_of_week": cosmosDBDatabaseAccountStringSet(basebackuppolicyresources.PossibleValuesForDayOfWeek()),

								"months_of_year": cosmosDBDatabaseAccountStringSet(basebackuppolicyresources.PossibleValuesForMonth()),

								"scheduled_backup_times": {
									Type:     pluginsdk.TypeSet,
									Optional: true,
									ForceNew: true,
									MinItems: 1,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.IsRFC3339Time,
									},
								},

								"weeks_of_month": cosmosDBDatabaseAccountStringSet(basebackuppolicyresources.PossibleValuesForWeekNumber()),
							},
						},
					},

					"life_cycle": cosmosDBDatabaseAccountLifeCycleSchema(),

					"priority": {
						Type:     pluginsdk.TypeInt,
						Required: true,
						ForceNew: true,
					},
				},
			},
		},

		"time_zone": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func (r DataProtectionBackupPolicyCosmosDBDatabaseAccountResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r DataProtectionBackupPolicyCosmosDBDatabaseAccountResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model BackupPolicyCosmosDBDatabaseAccountModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.DataProtection.BackupPolicyClient20260601
			vaultId, err := basebackuppolicyresources.ParseBackupVaultID(model.VaultId)
			if err != nil {
				return err
			}
			id := basebackuppolicyresources.NewBackupPolicyID(vaultId.SubscriptionId, vaultId.ResourceGroupName, vaultId.BackupVaultName, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.BackupPoliciesGet(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for existing %s: %+v", id, err)
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			policyRules := make([]basebackuppolicyresources20250701.BasePolicyRule, 0)
			policyRules = append(policyRules, expandBackupPolicyMySQLFlexibleServerAzureBackupRules(model.BackupRepeatingTimeIntervals, model.TimeZone, expandBackupPolicyMySQLFlexibleServerTaggingCriteria(model.RetentionRules))...)
			policyRules = append(policyRules, expandBackupPolicyMySQLFlexibleServerDefaultAzureRetentionRule(model.DefaultRetentionRule))
			policyRules = append(policyRules, expandBackupPolicyMySQLFlexibleServerAzureRetentionRules(model.RetentionRules)...)

			legacyParameters := basebackuppolicyresources20250701.BaseBackupPolicyResource{
				Properties: &basebackuppolicyresources20250701.BackupPolicy{
					ObjectType:      "BackupPolicy",
					PolicyRules:     policyRules,
					DatasourceTypes: []string{"Microsoft.DocumentDB/databaseAccounts"},
				},
			}
			parameters, err := convertDataProtectionModel[basebackuppolicyresources.BaseBackupPolicyResource](legacyParameters)
			if err != nil {
				return err
			}
			if _, err := client.BackupPoliciesCreateOrUpdate(ctx, id, *parameters); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return pluginsdk.SetResourceIdentityData(metadata.ResourceData, &id)
		},
	}
}

func (r DataProtectionBackupPolicyCosmosDBDatabaseAccountResource) Read() sdk.ResourceFunc {
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

			state := BackupPolicyCosmosDBDatabaseAccountModel{
				Name:    id.BackupPolicyName,
				VaultId: basebackuppolicyresources.NewBackupVaultID(id.SubscriptionId, id.ResourceGroupName, id.BackupVaultName).ID(),
			}
			if model := resp.Model; model != nil {
				legacyModel, err := convertDataProtectionModel[basebackuppolicyresources20250701.BaseBackupPolicyResource](*model)
				if err != nil {
					return err
				}
				if properties, ok := legacyModel.Properties.(basebackuppolicyresources20250701.BackupPolicy); ok {
					state.DefaultRetentionRule = flattenBackupPolicyMySQLFlexibleServerDefaultRetentionRule(properties.PolicyRules)
					state.RetentionRules = flattenBackupPolicyMySQLFlexibleServerRetentionRules(properties.PolicyRules)
					state.BackupRepeatingTimeIntervals = flattenBackupPolicyMySQLFlexibleServerBackupRules(properties.PolicyRules)
					state.TimeZone = flattenBackupPolicyMySQLFlexibleServerBackupTimeZone(properties.PolicyRules)
				}
			}

			if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
				return err
			}
			return metadata.Encode(&state)
		},
	}
}

func (r DataProtectionBackupPolicyCosmosDBDatabaseAccountResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DataProtection.BackupPolicyClient20260601
			id, err := basebackuppolicyresources.ParseBackupPolicyID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}
			if _, err := client.BackupPoliciesDelete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}
			return nil
		},
	}
}

func cosmosDBDatabaseAccountLifeCycleSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Required: true,
		ForceNew: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"data_store_type": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ForceNew:     true,
					ValidateFunc: validation.StringInSlice([]string{string(basebackuppolicyresources.DataStoreTypesVaultStore)}, false),
				},

				"duration": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ForceNew:     true,
					ValidateFunc: azValidate.ISO8601Duration,
				},
			},
		},
	}
}

func cosmosDBDatabaseAccountStringSet(values []string) *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeSet,
		Optional: true,
		ForceNew: true,
		MinItems: 1,
		Elem: &pluginsdk.Schema{
			Type:         pluginsdk.TypeString,
			ValidateFunc: validation.StringInSlice(values, false),
		},
	}
}
