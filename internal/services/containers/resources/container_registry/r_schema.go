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
			Description:  "The name of the Container Registry. Only Alphanumeric characters allowed.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ContainerRegistryName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"sku": {
			Description: "The SKU name of the container registry.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ValidateFunc: validation.StringInSlice([]string{
				string(registries.SkuNameBasic),
				string(registries.SkuNameStandard),
				string(registries.SkuNamePremium),
			}, false),
		},

		"admin_enabled": {
			Description: "Whether the admin user is enabled.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
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
						Description: "Whether this geo-replicated location participates in global endpoint routing for the Container Registry's geo-replicated login server.",
						Type:        pluginsdk.TypeBool,
						Required:    true,
					},

					"zone_redundancy_enabled": {
						Description: "Whether zone redundancy is enabled for this replication location.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     false,
					},

					"tags": commonschema.Tags(),
				},
			},
		},

		"public_network_access_enabled": {
			Description: "Whether public network access is allowed for the container registry.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"login_server": {
			Description: "The URL that can be used to log into the container registry.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"admin_username": {
			Description: "The Username associated with the Container Registry Admin account - if the admin account is enabled.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"admin_password": {
			Description: "The Password associated with the Container Registry Admin account - if the admin account is enabled.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Sensitive:   true,
		},

		"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

		"encryption": {
			Type:       pluginsdk.TypeList,
			Optional:   true,
			MaxItems:   1,
			ConfigMode: pluginsdk.SchemaConfigModeAttr,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"identity_client_id": {
						Description:  "The client ID of the Managed Identity associated with the encryption key.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.IsUUID,
					},
					"key_vault_key_id": {
						Description:  "The ID of the Key Vault Key.",
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
			MaxItems: 1,
			// ConfigModeAttr ensures we can set this to an empty array for Premium -> Basic
			ConfigMode: pluginsdk.SchemaConfigModeAttr,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"default_action": {
						Description:  "The behaviour for requests matching no rules.",
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
									Description:  "The behaviour for requests matching this rule.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringInSlice(registries.PossibleValuesForAction(), false),
								},
								"ip_range": {
									Description:  "The CIDR block from which requests will match the rule.",
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
			Description: "Boolean value that indicates whether quarantine policy is enabled.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
		},

		"retention_policy_in_days": {
			Description:  "The number of days to retain an untagged manifest after which it gets purged.",
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntBetween(0, 365),
		},

		"export_policy_enabled": {
			Description: "Boolean value that indicates whether export policy is enabled. To set it to `false`, make sure the `public_network_access_enabled` is also set to `false`.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"azuread_authentication_as_arm_policy_enabled": {
			Description: "Whether to use Azure Resource Manager audience token for this Container Registry.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"zone_redundancy_enabled": {
			Description: "Whether zone redundancy is enabled for this Container Registry.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			ForceNew:    true,
			Default:     false,
		},

		"anonymous_pull_enabled": {
			Description: "Whether to allow anonymous (unauthenticated) pull access to this Container Registry. This is only supported on resources with the `Standard` or `Premium` SKU.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
		},

		"data_endpoint_enabled": {
			Description: "Whether to enable dedicated data endpoints for this Container Registry. This is only supported on resources with the `Premium` SKU.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
		},

		"data_endpoint_host_names": {
			Description: "A set of data endpoint hostnames associated with the container registry if data endpoints are enabled.",
			Type:        pluginsdk.TypeSet,
			Computed:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"network_rule_bypass_option": {
			Description:  "Whether to allow trusted Azure services to access a network-restricted Container Registry.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(registries.NetworkRuleBypassOptionsAzureServices),
			ValidateFunc: validation.StringInSlice(registries.PossibleValuesForNetworkRuleBypassOptions(), false),
		},

		"network_rule_bypass_for_tasks_enabled": {
			Description: "Whether to allow Container Registry Tasks to access a network-restricted Container Registry.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"role_assignment_mode": {
			Description:  "The Role Assignment mode of this Container Registry.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      registries.RoleAssignmentModeLegacyRegistryPermissions,
			ValidateFunc: validation.StringInSlice(registries.PossibleValuesForRoleAssignmentMode(), false),
		},

		"tags": commonschema.Tags(),
	}
}
