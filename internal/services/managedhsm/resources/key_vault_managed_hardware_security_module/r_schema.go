// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package key_vault_managed_hardware_security_module

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-sdk/resource-manager/keyvault/2026-02-01/managedhsms"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/managedhsm/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func keyVaultManagedHardwareSecurityModuleSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ManagedHardwareSecurityModuleName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"sku_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(managedhsms.ManagedHsmSkuNameStandardBOne),
			}, false),
		},

		"admin_object_ids": {
			Type:     pluginsdk.TypeSet,
			Required: true,
			ForceNew: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.IsUUID,
			},
		},

		"tenant_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsUUID,
		},

		"purge_protection_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			ForceNew: true,
		},

		"soft_delete_retention_days": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ForceNew:     true,
			Default:      90,
			ValidateFunc: validation.IntBetween(7, 90),
		},

		"hsm_uri": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"public_network_access_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
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
						ValidateFunc: validation.StringInSlice(managedhsms.PossibleValuesForNetworkRuleAction(), false),
					},
					"bypass": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(managedhsms.PossibleValuesForNetworkRuleBypassOptions(), false),
					},
				},
			},
		},

		"security_domain_key_vault_certificate_ids": {
			Type:         pluginsdk.TypeList,
			MinItems:     3,
			MaxItems:     10,
			Optional:     true,
			RequiredWith: []string{"security_domain_quorum"},
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: keyvault.ValidateNestedItemID(keyvault.VersionTypeVersioned, keyvault.NestedItemTypeCertificate),
			},
		},

		"security_domain_quorum": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			RequiredWith: []string{"security_domain_key_vault_certificate_ids"},
			ValidateFunc: validation.IntBetween(2, 10),
		},

		"security_domain_encrypted_data": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},

		// https://github.com/Azure/azure-rest-api-specs/issues/13365
		"tags": commonschema.Tags(),
	}
}
