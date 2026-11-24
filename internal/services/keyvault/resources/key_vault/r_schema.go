// Copyright IBM Corp. 2014, 2026
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
			Description:  "The name of the Key Vault. The name must be globally unique. If the vault is in a recoverable state then the vault will need to be purged before reusing the name.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.VaultName,
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"rbac_authorization_enabled": {
			Description: "Boolean flag to specify whether Azure Key Vault uses Role Based Access Control (RBAC) for authorization of data actions.",
			Type:        pluginsdk.TypeBool,
			Required:    true,
		},

		"sku_name": {
			Description:  "The Name of the SKU used for this Key Vault.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(vaults.PossibleValuesForSkuName(), false),
		},

		"tenant_id": {
			Description:  "The Azure Active Directory tenant ID that should be used for authenticating requests to the Key Vault.",
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
						Description:  "The Azure Active Directory tenant ID that should be used for authenticating requests to the Key Vault. Must match the `tenant_id` used above.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.IsUUID,
					},
					"object_id": {
						Description:  "The object ID of a user, Service Principal or security group in the Azure Active Directory tenant for the vault. The object ID must be unique for the list of access policies.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.IsUUID,
					},
					"application_id": {
						Description:  "The object ID of an Application in Azure Active Directory.",
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
			Description: "Boolean flag to specify whether Azure Virtual Machines are permitted to retrieve certificates stored as secrets from the Key Vault.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
		},

		"enabled_for_disk_encryption": {
			Description: "Boolean flag to specify whether Azure Disk Encryption is permitted to retrieve secrets from the vault and unwrap keys.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
		},

		"enabled_for_template_deployment": {
			Description: "Boolean flag to specify whether Azure Resource Manager is permitted to retrieve secrets from the Key Vault.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
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
						Description:  "Which traffic can bypass the network rules.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(vaults.PossibleValuesForNetworkRuleBypassOptions(), false),
					},
					"ip_rules": {
						Description: "One or more IP Addresses, or CIDR Blocks which should be able to access the Key Vault.",
						Type:        pluginsdk.TypeSet,
						Optional:    true,
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
						Description: "One or more Subnet IDs which should be able to access this Key Vault.",
						Type:        pluginsdk.TypeSet,
						Optional:    true,
						Elem:        &pluginsdk.Schema{Type: pluginsdk.TypeString},
						Set:         set.HashStringIgnoreCase,
					},
				},
			},
		},

		"public_network_access_enabled": {
			Description: "Whether public network access is allowed for this Key Vault.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"purge_protection_enabled": {
			Description: "Is Purge Protection enabled for this Key Vault?",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
		},

		"soft_delete_retention_days": {
			Description:  "The number of days that items should be retained for once soft-deleted.",
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      90,
			ValidateFunc: validation.IntBetween(7, 90),
		},

		"tags": commonschema.Tags(),

		// Computed
		"vault_uri": {
			Description: "The URI of the Key Vault, used for performing operations on keys and secrets.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},
	}
}
