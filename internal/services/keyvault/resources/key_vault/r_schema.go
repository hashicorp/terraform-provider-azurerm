// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package key_vault

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/keyvault/2026-02-01/vaults"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/set"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func keyVaultSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.VaultName,
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"rbac_authorization_enabled": {
			Type:     pluginsdk.TypeBool,
			Required: true,
		},

		"sku_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(vaults.PossibleValuesForSkuName(), false),
		},

		"tenant_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.IsUUID,
		},

		"access_policy": {
			Type:       pluginsdk.TypeList,
			ConfigMode: pluginsdk.SchemaConfigModeAttr,
			Optional:   true,
			Computed:   true, // azignore:AZS007 - pre-existing violation
			MaxItems:   1024,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"tenant_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.IsUUID,
					},
					"object_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.IsUUID,
					},
					"application_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validate.IsUUIDOrEmpty,
					},
					"certificate_permissions": helpers.SchemaCertificatePermissions(),
					"key_permissions":         helpers.SchemaKeyPermissions(),
					"secret_permissions":      helpers.SchemaSecretPermissions(),
					"storage_permissions":     helpers.SchemaStoragePermissions(),
				},
			},
		},

		"enabled_for_deployment": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"enabled_for_disk_encryption": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"enabled_for_template_deployment": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"network_acls": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"default_action": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(vaults.PossibleValuesForNetworkRuleAction(), false),
					},
					"bypass": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(vaults.PossibleValuesForNetworkRuleBypassOptions(), false),
					},
					"ip_rules": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
							ValidateFunc: validation.Any(
								validation.IsIPv4Address,
								validation.IsCIDRIPv4,
							),
						},
						Set: set.HashIPv4AddressOrCIDR,
					},
					"virtual_network_subnet_ids": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
						Set:      set.HashStringIgnoreCase,
					},
				},
			},
		},

		"public_network_access_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"purge_protection_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"soft_delete_retention_days": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      90,
			ValidateFunc: validation.IntBetween(7, 90),
		},

		"tags": commonschema.Tags(),

		// Computed
		"vault_uri": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}
