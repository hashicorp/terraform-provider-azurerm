// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package databricks_workspace

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-sdk/resource-manager/databricks/2026-01-01/workspaces"
	mlworkspace "github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2025-06-01/workspaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/loadbalancers"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/databricks/validate"
	storageValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func databricksWorkspaceSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.WorkspaceName,
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"sku": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ValidateFunc: validation.StringInSlice([]string{
				"standard",
				"premium",
				"trial",
			}, false),
		},

		"managed_resource_group_name": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
			// NOTE: O+C We set a value for this if omitted so this should remain Computed
			Computed:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"customer_managed_key_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"managed_disk_identity": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"principal_id": {
						Type:      pluginsdk.TypeString,
						Sensitive: true,
						Computed:  true,
					},

					"tenant_id": {
						Type:      pluginsdk.TypeString,
						Sensitive: true,
						Computed:  true,
					},

					"type": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"infrastructure_encryption_enabled": {
			Type:     pluginsdk.TypeBool,
			ForceNew: true,
			Optional: true,
			Default:  false,
		},

		"public_network_access_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"default_storage_firewall_enabled": {
			Type:         pluginsdk.TypeBool,
			Optional:     true,
			RequiredWith: []string{"access_connector_id"},
		},

		"access_connector_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			RequiredWith: []string{"default_storage_firewall_enabled"},
		},

		"network_security_group_rules_required": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice(workspaces.PossibleValuesForRequiredNsgRules(), false),
		},

		"load_balancer_backend_address_pool_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: loadbalancers.ValidateLoadBalancerBackendAddressPoolID,
		},

		"custom_parameters": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			// NOTE: O+C The API populates these and since many are ForceNew there doesn't appear to be a need to remove this once set to use the defaults
			Computed: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"machine_learning_workspace_id": {
						Type:         pluginsdk.TypeString,
						ForceNew:     true,
						Optional:     true,
						ValidateFunc: mlworkspace.ValidateWorkspaceID,
						AtLeastOneOf: workspaceCustomParametersString(),
					},

					"nat_gateway_name": {
						Type:         pluginsdk.TypeString,
						ForceNew:     true,
						Optional:     true,
						Computed:     true, // azignore:AZS007 - pre-existing violation
						AtLeastOneOf: workspaceCustomParametersString(),
					},

					"no_public_ip": {
						Type:         pluginsdk.TypeBool,
						Optional:     true,
						Default:      true,
						AtLeastOneOf: workspaceCustomParametersString(),
					},

					"public_ip_name": {
						Type:         pluginsdk.TypeString,
						ForceNew:     true,
						Optional:     true,
						Computed:     true, // azignore:AZS007 - pre-existing violation
						AtLeastOneOf: workspaceCustomParametersString(),
					},

					"public_subnet_name": {
						Type:         pluginsdk.TypeString,
						ForceNew:     true,
						Optional:     true,
						AtLeastOneOf: workspaceCustomParametersString(),
					},

					"public_subnet_network_security_group_association_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: azure.ValidateResourceID,
						AtLeastOneOf: workspaceCustomParametersString(),
					},

					"private_subnet_name": {
						Type:         pluginsdk.TypeString,
						ForceNew:     true,
						Optional:     true,
						AtLeastOneOf: workspaceCustomParametersString(),
					},

					"private_subnet_network_security_group_association_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: azure.ValidateResourceID,
						AtLeastOneOf: workspaceCustomParametersString(),
					},

					"virtual_network_id": {
						Type:         pluginsdk.TypeString,
						ForceNew:     true,
						Optional:     true,
						ValidateFunc: commonids.ValidateVirtualNetworkID,
						AtLeastOneOf: workspaceCustomParametersString(),
					},

					"storage_account_name": {
						Type:         pluginsdk.TypeString,
						ForceNew:     true,
						Optional:     true,
						Computed:     true, // azignore:AZS007 - pre-existing violation
						ValidateFunc: storageValidate.StorageAccountName,
						AtLeastOneOf: workspaceCustomParametersString(),
					},

					"storage_account_sku_name": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Computed:     true, // azignore:AZS007 - pre-existing violation
						AtLeastOneOf: workspaceCustomParametersString(),
					},

					"vnet_address_prefix": {
						Type:         pluginsdk.TypeString,
						ForceNew:     true,
						Optional:     true,
						Computed:     true, // azignore:AZS007 - pre-existing violation
						AtLeastOneOf: workspaceCustomParametersString(),
					},
				},
			},
		},

		"managed_resource_group_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"workspace_url": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"workspace_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"managed_services_cmk_key_vault_key_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: keyvault.ValidateNestedItemID(keyvault.VersionTypeVersioned, keyvault.NestedItemTypeKey),
		},

		"managed_disk_cmk_key_vault_key_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: keyvault.ValidateNestedItemID(keyvault.VersionTypeVersioned, keyvault.NestedItemTypeKey),
		},

		"managed_disk_cmk_rotation_to_latest_version_enabled": {
			Type:         pluginsdk.TypeBool,
			Optional:     true,
			RequiredWith: []string{"managed_disk_cmk_key_vault_key_id"},
		},

		"disk_encryption_set_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"storage_account_identity": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"principal_id": {
						Type:      pluginsdk.TypeString,
						Sensitive: true,
						Computed:  true,
					},

					"tenant_id": {
						Type:      pluginsdk.TypeString,
						Sensitive: true,
						Computed:  true,
					},

					"type": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"enhanced_security_compliance": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"automatic_cluster_update_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},
					"compliance_security_profile_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},
					"compliance_security_profile_standards": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringInSlice(validate.PossibleValuesForComplianceStandard(), false),
						},
					},
					"enhanced_security_monitoring_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},
				},
			},
		},

		"tags": commonschema.Tags(),
	}
}
