// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package container_registry

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2025-11-01/registries"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func containerRegistrySchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ContainerRegistryName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"sku": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(registries.SkuNameBasic),
				string(registries.SkuNameStandard),
				string(registries.SkuNamePremium),
			}, false),
		},

		"admin_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"georeplications": {
			// Don't make this a TypeSet since TypeSet has bugs when there is a nested property using `StateFunc`.
			// See: https://github.com/hashicorp/terraform-plugin-sdk/issues/160
			Type:       pluginsdk.TypeList,
			Optional:   true,
			ConfigMode: pluginsdk.SchemaConfigModeAuto,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"location": commonschema.LocationWithoutForceNew(),

					"global_endpoint_routing_enabled": {
						Type:     pluginsdk.TypeBool,
						Required: true,
					},

					"zone_redundancy_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},

					"tags": commonschema.Tags(),
				},
			},
		},

		"public_network_access_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"login_server": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"admin_username": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"admin_password": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},

		"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

		"encryption": {
			Type:       pluginsdk.TypeList,
			Optional:   true,
			ConfigMode: pluginsdk.SchemaConfigModeAttr,
			MaxItems:   1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"identity_client_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.IsUUID,
					},
					"key_vault_key_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: keyvault.ValidateNestedItemID(keyvault.VersionTypeAny, keyvault.NestedItemTypeKey),
					},
				},
			},
		},

		"network_rule_set": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
			// ConfigModeAttr ensures we can set this to an empty array for Premium -> Basic
			ConfigMode: pluginsdk.SchemaConfigModeAttr,
			MaxItems:   1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"default_action": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      registries.DefaultActionAllow,
						ValidateFunc: validation.StringInSlice(registries.PossibleValuesForDefaultAction(), false),
					},

					"ip_rule": {
						Type:       pluginsdk.TypeSet,
						Optional:   true,
						Computed:   true, // azignore:AZS007 - pre-existing violation
						ConfigMode: pluginsdk.SchemaConfigModeAttr,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"action": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringInSlice(registries.PossibleValuesForAction(), false),
								},
								"ip_range": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.IsCIDRIPv4,
								},
							},
						},
					},
				},
			},
		},

		"quarantine_policy_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"retention_policy_in_days": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntBetween(0, 365),
		},

		"export_policy_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"azuread_authentication_as_arm_policy_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"zone_redundancy_enabled": {
			Type:     pluginsdk.TypeBool,
			ForceNew: true,
			Optional: true,
			Default:  false,
		},

		"anonymous_pull_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"data_endpoint_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"data_endpoint_host_names": {
			Type:     pluginsdk.TypeSet,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"network_rule_bypass_option": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice(registries.PossibleValuesForNetworkRuleBypassOptions(), false),
			Default:      string(registries.NetworkRuleBypassOptionsAzureServices),
		},

		"network_rule_bypass_for_tasks_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"role_assignment_mode": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice(registries.PossibleValuesForRoleAssignmentMode(), false),
			Default:      registries.RoleAssignmentModeLegacyRegistryPermissions,
		},

		"tags": commonschema.Tags(),
	}
}
