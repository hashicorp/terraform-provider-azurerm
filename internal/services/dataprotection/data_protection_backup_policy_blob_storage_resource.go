// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package dataprotection

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2024-04-01/backuppolicies"
	helperValidate "github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/dataprotection/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/dataprotection/validate"

	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type DataProtectionBackupPolicyBlobStorageResource struct{}

type DataProtectionBackupPolicyBlobStorageResourceModel struct {
	Name                                string                  `tfschema:"name"`
	VaultId                             string                  `tfschema:"vault_id"`
	BackupRepeatingTimeIntervals        []string                `tfschema:"backup_repeating_time_intervals"`
	OperationalDefaultRetentionDuration string                  `tfschema:"operational_default_retention_duration"`
	TimeZone                            string                  `tfschema:"time_zone"`
	VaultDefaultRetentionDuration       string                  `tfschema:"vault_default_retention_duration"`
	RetentionRule                       []helpers.RetentionRule `tfschema:"retention_rule"`
}

var _ sdk.Resource = DataProtectionBackupPolicyBlobStorageResource{}

func (r DataProtectionBackupPolicyBlobStorageResource) ResourceType() string {
	return "azurerm_data_protection_backup_policy_blob_storage"
}

func (r DataProtectionBackupPolicyBlobStorageResource) ModelObject() interface{} {
	return &DataProtectionBackupPolicyBlobStorageResourceModel{}
}

func (r DataProtectionBackupPolicyBlobStorageResource) Arguments() map[string]*pluginsdk.Schema {
	arguments := map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.BackupPolicyBlobStorageName,
		},

		"vault_id": commonschema.ResourceIDReferenceRequiredForceNew(pointer.To(backuppolicies.BackupVaultId{})),

		"backup_repeating_time_intervals": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			MinItems: 1,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: helperValidate.ISO8601RepeatingTime,
			},
		},

		"operational_default_retention_duration": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			AtLeastOneOf: []string{"operational_default_retention_duration", "vault_default_retention_duration"},
			ValidateFunc: helperValidate.ISO8601Duration,
		},

		"time_zone": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validate.BackupPolicyTimeZone(),
		},

		"vault_default_retention_duration": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			AtLeastOneOf: []string{"operational_default_retention_duration", "vault_default_retention_duration"},
			RequiredWith: []string{"backup_repeating_time_intervals"},
			ValidateFunc: helperValidate.ISO8601Duration,
		},

		"retention_rule": {
			Type:         pluginsdk.TypeList,
			Optional:     true,
			ForceNew:     true,
			RequiredWith: []string{"vault_default_retention_duration"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ForceNew: true,
					},

					"criteria": {
						Type:     pluginsdk.TypeList,
						Required: true,
						ForceNew: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"absolute_criteria": {
									Type:     pluginsdk.TypeString,
									Optional: true,
									ForceNew: true,
									ValidateFunc: validation.StringInSlice(
										backuppolicies.PossibleValuesForAbsoluteMarker(), false),
								},

								"days_of_month": {
									Type:     pluginsdk.TypeSet,
									Optional: true,
									ForceNew: true,
									MinItems: 1,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeInt,
										ValidateFunc: validation.Any(
											validation.IntBetween(0, 28),
										),
									},
								},

								"days_of_week": {
									Type:     pluginsdk.TypeSet,
									Optional: true,
									ForceNew: true,
									MinItems: 1,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.IsDayOfTheWeek(false),
									},
								},

								"months_of_year": {
									Type:     pluginsdk.TypeSet,
									Optional: true,
									ForceNew: true,
									MinItems: 1,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.IsMonth(false),
									},
								},

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

								"weeks_of_month": {
									Type:     pluginsdk.TypeSet,
									Optional: true,
									ForceNew: true,
									MinItems: 1,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.StringInSlice(backuppolicies.PossibleValuesForWeekNumber(), false),
									},
								},
							},
						},
					},

					"life_cycle": {
						Type:     pluginsdk.TypeList,
						Required: true,
						ForceNew: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"data_store_type": {
									Type:     pluginsdk.TypeString,
									Required: true,
									ForceNew: true,
									ValidateFunc: validation.StringInSlice([]string{
										// confirmed with the service team that currently only `VaultStore` is supported.
										// However, since `ArchiveStore` may be supported in the future, it is open to user specification.
										string(backuppolicies.DataStoreTypesVaultStore),
									}, false),
								},

								"duration": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ForceNew:     true,
									ValidateFunc: helperValidate.ISO8601Duration,
								},
							},
						},
					},

					"priority": {
						Type:     pluginsdk.TypeInt,
						Required: true,
						ForceNew: true,
					},
				},
			},
		},
	}
	return arguments
}

func (r DataProtectionBackupPolicyBlobStorageResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r DataProtectionBackupPolicyBlobStorageResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model DataProtectionBackupPolicyBlobStorageResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.DataProtection.BackupPolicyClient
			subscriptionId := metadata.Client.Account.SubscriptionId
			vaultId, _ := backuppolicies.ParseBackupVaultID(model.VaultId)
			id := backuppolicies.NewBackupPolicyID(subscriptionId, vaultId.ResourceGroupName, vaultId.BackupVaultName, model.Name)
			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for existing %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			backupType := "Discrete"
			policyRules := make([]backuppolicies.BasePolicyRule, 0)

			// expand the default operational retention rule when the operational default duration is specified
			operationalDefaultDuration := model.OperationalDefaultRetentionDuration
			if operationalDefaultDuration != "" {
				policyRules = append(policyRules, helpers.ExpandBackupPolicyDefaultRetentionRuleArray(operationalDefaultDuration, backuppolicies.DataStoreTypesOperationalStore))
			}

			vaultDefaultRetentionDuration := model.VaultDefaultRetentionDuration
			if vaultDefaultRetentionDuration != "" {
				taggingCriteria, err := helpers.ExpandBackupPolicyTaggingCriteriaArray(model.RetentionRule)
				if err != nil {
					return err
				}
				policyRules = append(policyRules, helpers.ExpandBackupPolicyAzureBackupRuleArray(model.BackupRepeatingTimeIntervals, model.TimeZone, backupType, backuppolicies.DataStoreTypesVaultStore, taggingCriteria)...)
				policyRules = append(policyRules, helpers.ExpandBackupPolicyDefaultRetentionRuleArray(vaultDefaultRetentionDuration, backuppolicies.DataStoreTypesVaultStore))
			}

			// expand the vault retention rule when the vault retention rules are specified, the operational backup cannot specify retention rules.
			retentionRule := model.RetentionRule
			if len(retentionRule) > 0 {
				log.Printf("[DEBUG] retention rules %v", retentionRule)
				policyRules = append(policyRules, helpers.ExpandBackupPolicyAzureRetentionRules(model.RetentionRule)...)
			}

			parameters := backuppolicies.BaseBackupPolicyResource{
				Properties: &backuppolicies.BackupPolicy{
					PolicyRules:     policyRules,
					DatasourceTypes: []string{"Microsoft.Storage/storageAccounts/blobServices"},
				},
			}

			if _, err := client.CreateOrUpdate(ctx, id, parameters); err != nil {
				return fmt.Errorf("creating/updating DataProtection BackupPolicy (%q): %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (r DataProtectionBackupPolicyBlobStorageResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DataProtection.BackupPolicyClient

			id, err := backuppolicies.ParseBackupPolicyID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			vaultId := backuppolicies.NewBackupVaultID(id.SubscriptionId, id.ResourceGroupName, id.BackupVaultName)
			state := DataProtectionBackupPolicyBlobStorageResourceModel{
				Name:    id.BackupPolicyName,
				VaultId: vaultId.ID(),
			}

			if model := resp.Model; model != nil {
				if properties, ok := model.Properties.(backuppolicies.BackupPolicy); ok {
					state.RetentionRule = helpers.FlattenBackupPolicyRetentionRules(&properties.PolicyRules)
					state.BackupRepeatingTimeIntervals = helpers.FlattenBackupPolicyBackupRuleArray(&properties.PolicyRules)
					state.TimeZone = helpers.FlattenBackupPolicyBackupTimeZone(&properties.PolicyRules)
					state.OperationalDefaultRetentionDuration = helpers.FlattenBackupPolicyDefaultRetentionRuleDuration(&properties.PolicyRules, backuppolicies.DataStoreTypesOperationalStore)
					state.VaultDefaultRetentionDuration = helpers.FlattenBackupPolicyDefaultRetentionRuleDuration(&properties.PolicyRules, backuppolicies.DataStoreTypesVaultStore)
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r DataProtectionBackupPolicyBlobStorageResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DataProtection.BackupPolicyClient

			id, err := backuppolicies.ParseBackupPolicyID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.Delete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r DataProtectionBackupPolicyBlobStorageResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return backuppolicies.ValidateBackupPolicyID
}
