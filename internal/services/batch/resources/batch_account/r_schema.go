// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package batch_account

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-sdk/resource-manager/batch/2024-07-01/batchaccount"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/batch/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func batchAccountSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the Batch Account. Only lowercase Alphanumeric characters allowed.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.AccountName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"storage_account_id": {
			Description:  "The Storage Account to use for the Batch Account. If not specified, Azure Batch will manage the storage.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: commonids.ValidateStorageAccountID,
			RequiredWith: []string{"storage_account_authentication_mode"},
		},

		"storage_account_authentication_mode": {
			Description:  "The Storage Account authentication mode.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice(batchaccount.PossibleValuesForAutoStorageAuthenticationMode(), false),
			RequiredWith: []string{"storage_account_id"},
		},

		"storage_account_node_identity": {
			Description:  "The User Assigned Identity for the Storage Account.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: commonids.ValidateUserAssignedIdentityID,
			RequiredWith: []string{"storage_account_id"},
		},

		"allowed_authentication_modes": {
			Description: "The allowed authentication mode for the Batch Account.",
			Type:        pluginsdk.TypeSet,
			Optional:    true,
			// NOTE: O+C This can remain since we need to send an empty slice to properly remove this, which means this can't be set back to
			// its default which is to return all three values
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringInSlice(batchaccount.PossibleValuesForAuthenticationMode(), false),
			},
		},

		"pool_allocation_mode": {
			Description:  "The mode to use for pool allocation.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(batchaccount.PoolAllocationModeBatchService),
			ValidateFunc: validation.StringInSlice(batchaccount.PossibleValuesForPoolAllocationMode(), false),
		},

		"public_network_access_enabled": {
			Description: "Whether public network access is allowed for this server.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"network_profile": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"account_access": resourceBatchAccountEndpointAccessProfileSchema(),

					"node_management_access": resourceBatchAccountEndpointAccessProfileSchema(),
				},
			},
		},

		"key_vault_reference": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"id": {
						Description:  "The Azure identifier of the Azure KeyVault to use.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: azure.ValidateResourceID,
					},
					"url": {
						Description:  "The HTTPS URL of the Azure KeyVault to use.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.IsURLWithHTTPS,
					},
				},
			},
		},

		"identity": commonschema.SystemOrUserAssignedIdentityOptional(),

		"primary_access_key": {
			Description: "The Batch Account primary access key.",
			Type:        pluginsdk.TypeString,
			Sensitive:   true,
			Computed:    true,
		},

		"secondary_access_key": {
			Description: "The Batch Account secondary access key.",
			Type:        pluginsdk.TypeString,
			Sensitive:   true,
			Computed:    true,
		},

		"account_endpoint": {
			Description: "The account endpoint used to interact with the Batch service.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},
		"encryption": {
			Description: "Whether customer managed key encryption should be used to encrypt Batch Account data.",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			MaxItems:    1,
			ConfigMode:  pluginsdk.SchemaConfigModeAttr,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"key_vault_key_id": {
						Description:  "The full URL path to the Azure Key Vault Key ID that should be used to encrypt data, as documented [here](https://docs.microsoft.com/azure/batch/batch-customer-managed-key). Both versioned and versionless keys are supported.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: keyvault.ValidateNestedItemID(keyvault.VersionTypeAny, keyvault.NestedItemTypeKey),
					},
				},
			},
		},

		"tags": commonschema.Tags(),
	}
}

func resourceBatchAccountEndpointAccessProfileSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:         pluginsdk.TypeList,
		Optional:     true,
		MaxItems:     1,
		AtLeastOneOf: []string{"network_profile.0.account_access", "network_profile.0.node_management_access"},
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"default_action": {
					Description:  "The default action.",
					Type:         pluginsdk.TypeString,
					Optional:     true,
					Default:      batchaccount.EndpointAccessDefaultActionDeny,
					ValidateFunc: validation.StringInSlice(batchaccount.PossibleValuesForEndpointAccessDefaultAction(), false),
				},

				"ip_rule": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"ip_range": {
								Description:  "The CIDR block from which requests will match the rule.",
								Type:         pluginsdk.TypeString,
								Required:     true,
								ValidateFunc: validate.BatchAccountIpRange,
							},

							"action": {
								Description:  "The action of the IP rule.",
								Type:         pluginsdk.TypeString,
								Optional:     true,
								Default:      string(batchaccount.IPRuleActionAllow),
								ValidateFunc: validation.StringInSlice(batchaccount.PossibleValuesForIPRuleAction(), false),
							},
						},
					},
				},
			},
		},
	}
}
