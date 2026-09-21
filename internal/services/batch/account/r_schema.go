// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package account

import (
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-sdk/resource-manager/batch/2024-07-01/batchaccount"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/batch/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func resourceBatchAccountSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.AccountName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"storage_account_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: commonids.ValidateStorageAccountID,
			RequiredWith: []string{"storage_account_authentication_mode"},
		},

		"storage_account_authentication_mode": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice(batchaccount.PossibleValuesForAutoStorageAuthenticationMode(), false),
			RequiredWith: []string{"storage_account_id"},
		},

		"storage_account_node_identity": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: commonids.ValidateUserAssignedIdentityID,
			RequiredWith: []string{"storage_account_id"},
		},

		"allowed_authentication_modes": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			// NOTE: O+C This can remain since we need to send an empty slice to properly remove this, which means this can't be set back to
			// its default which is to return all three values
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringInSlice(batchaccount.PossibleValuesForAuthenticationMode(), false),
			},
		},

		"pool_allocation_mode": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(batchaccount.PoolAllocationModeBatchService),
			ValidateFunc: validation.StringInSlice(batchaccount.PossibleValuesForPoolAllocationMode(), false),
		},

		"public_network_access_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
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
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: azure.ValidateResourceID,
					},
					"url": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.IsURLWithHTTPS,
					},
				},
			},
		},

		"identity": commonschema.SystemOrUserAssignedIdentityOptional(),

		"primary_access_key": {
			Type:      pluginsdk.TypeString,
			Sensitive: true,
			Computed:  true,
		},

		"secondary_access_key": {
			Type:      pluginsdk.TypeString,
			Sensitive: true,
			Computed:  true,
		},

		"account_endpoint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"encryption": {
			Type:       pluginsdk.TypeList,
			Optional:   true,
			MaxItems:   1,
			ConfigMode: pluginsdk.SchemaConfigModeAttr,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"key_vault_key_id": {
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
								Type:         pluginsdk.TypeString,
								Required:     true,
								ValidateFunc: validate.BatchAccountIpRange,
							},

							"action": {
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

func isShardKeyAllowed(input []interface{}) bool {
	if len(input) == 0 {
		return false
	}
	for _, authMod := range input {
		if strings.EqualFold(authMod.(string), string(batchaccount.AuthenticationModeSharedKey)) {
			return true
		}
	}
	return false
}
