// Copyright IBM Corp. 2014, 2026
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
			Description:  "The name of the Databricks Workspace resource.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.WorkspaceName,
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"sku": {
			Description: "The `sku` to use for the Databricks Workspace.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ValidateFunc: validation.StringInSlice([]string{
				"standard",
				"premium",
				"trial",
			}, false),
		},

		"managed_resource_group_name": {
			Description: "The name of the Resource Group where Azure should place the managed Databricks resources.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			ForceNew:    true,
			// NOTE: O+C We set a value for this if omitted so this should remain Computed
			Computed:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"customer_managed_key_enabled": {
			Description: "Is the workspace enabled for customer managed key encryption? If `true` this enables the Managed Identity for the managed Storage Account. This field is only valid if the Databricks Workspace `sku` is set to `premium`.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"managed_disk_identity": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"principal_id": {
						Description: "The principal UUID for the internal databricks disks identity needed to provide access to the workspace for enabling Customer Managed Keys.",
						Type:        pluginsdk.TypeString,
						Sensitive:   true,
						Computed:    true,
					},

					"tenant_id": {
						Description: "The UUID of the tenant where the internal databricks disks identity was created.",
						Type:        pluginsdk.TypeString,
						Sensitive:   true,
						Computed:    true,
					},

					"type": {
						Description: "The type of the internal databricks disks identity.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},
				},
			},
		},

		"infrastructure_encryption_enabled": {
			Description: "Is the Databricks File System root file system enabled with a secondary layer of encryption with platform managed keys? This field is only valid if the Databricks Workspace `sku` is set to `premium`.",
			Type:        pluginsdk.TypeBool,
			ForceNew:    true,
			Optional:    true,
			Default:     false,
		},

		"public_network_access_enabled": {
			Description: "Allow public access for accessing workspace. Set value to `false` to access workspace only via Private Link endpoint.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"default_storage_firewall_enabled": {
			Description:  "Disallow public access to default Storage Account.",
			Type:         pluginsdk.TypeBool,
			Optional:     true,
			RequiredWith: []string{"access_connector_id"},
		},

		"access_connector_id": {
			Description:  "Access Connector ID to use when default Storage Account firewall is enabled.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			RequiredWith: []string{"default_storage_firewall_enabled"},
		},

		"network_security_group_rules_required": {
			Description:  "Does the data plane (clusters) to control plane communication happen over Private Link endpoint only or publicly?",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice(workspaces.PossibleValuesForRequiredNsgRules(), false),
		},

		"load_balancer_backend_address_pool_id": {
			Description:  "Resource ID of the Outbound Load balancer Backend Address Pool for Secure Cluster Connectivity (No Public IP) workspace with managed Virtual Network.",
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
						Description:  "The ID of an Azure Machine Learning Workspace to link with Databricks Workspace.",
						Type:         pluginsdk.TypeString,
						ForceNew:     true,
						Optional:     true,
						ValidateFunc: mlworkspace.ValidateWorkspaceID,
						AtLeastOneOf: workspaceCustomParametersString(),
					},

					"nat_gateway_name": {
						Description:  "Name of the NAT gateway for Secure Cluster Connectivity (No Public IP) workspace subnets (only for workspace with managed Virtual Network).",
						Type:         pluginsdk.TypeString,
						ForceNew:     true,
						Optional:     true,
						Computed:     true, // azignore:AZS007 - pre-existing violation
						AtLeastOneOf: workspaceCustomParametersString(),
					},

					"no_public_ip": {
						Description:  "Are public IP Addresses not allowed?",
						Type:         pluginsdk.TypeBool,
						Optional:     true,
						Default:      true,
						AtLeastOneOf: workspaceCustomParametersString(),
					},

					"public_ip_name": {
						Description:  "Name of the Public IP for No Public IP workspace with managed Virtual Network.",
						Type:         pluginsdk.TypeString,
						ForceNew:     true,
						Optional:     true,
						Computed:     true, // azignore:AZS007 - pre-existing violation
						AtLeastOneOf: workspaceCustomParametersString(),
					},

					"public_subnet_name": {
						Description:  "The name of the Public Subnet within the Virtual Network.",
						Type:         pluginsdk.TypeString,
						ForceNew:     true,
						Optional:     true,
						AtLeastOneOf: workspaceCustomParametersString(),
					},

					"public_subnet_network_security_group_association_id": {
						Description:  "The resource ID of the `azurerm_subnet_network_security_group_association` resource which is referred to by the `public_subnet_name` field. This is the same as the ID of the subnet referred to by the `public_subnet_name` field.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: azure.ValidateResourceID,
						AtLeastOneOf: workspaceCustomParametersString(),
					},

					"private_subnet_name": {
						Description:  "The name of the Private Subnet within the Virtual Network.",
						Type:         pluginsdk.TypeString,
						ForceNew:     true,
						Optional:     true,
						AtLeastOneOf: workspaceCustomParametersString(),
					},

					"private_subnet_network_security_group_association_id": {
						Description:  "The resource ID of the `azurerm_subnet_network_security_group_association` resource which is referred to by the `private_subnet_name` field. This is the same as the ID of the subnet referred to by the `private_subnet_name` field.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: azure.ValidateResourceID,
						AtLeastOneOf: workspaceCustomParametersString(),
					},

					"virtual_network_id": {
						Description:  "The ID of a Virtual Network where this Databricks Cluster should be created.",
						Type:         pluginsdk.TypeString,
						ForceNew:     true,
						Optional:     true,
						ValidateFunc: commonids.ValidateVirtualNetworkID,
						AtLeastOneOf: workspaceCustomParametersString(),
					},

					"storage_account_name": {
						Description:  "Default Databricks File Storage Account name.",
						Type:         pluginsdk.TypeString,
						ForceNew:     true,
						Optional:     true,
						Computed:     true, // azignore:AZS007 - pre-existing violation
						ValidateFunc: storageValidate.StorageAccountName,
						AtLeastOneOf: workspaceCustomParametersString(),
					},

					"storage_account_sku_name": {
						Description:  "Storage Account SKU name.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Computed:     true, // azignore:AZS007 - pre-existing violation
						AtLeastOneOf: workspaceCustomParametersString(),
					},

					"vnet_address_prefix": {
						Description:  "Address prefix for Managed Virtual Network.",
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
			Description: "The ID of the Managed Resource Group created by the Databricks Workspace.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"workspace_url": {
			Description: "The workspace URL which is of the format `adb-{workspaceId}.{random}.azuredatabricks.net`.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"workspace_id": {
			Description: "The unique identifier of the Databricks Workspace in Databricks control plane.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"managed_services_cmk_key_vault_key_id": {
			Description:  "Customer managed encryption properties for the Databricks Workspace managed resources (e.g. Notebooks and Artifacts).",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: keyvault.ValidateNestedItemID(keyvault.VersionTypeVersioned, keyvault.NestedItemTypeKey),
		},

		"managed_disk_cmk_key_vault_key_id": {
			Description:  "Customer managed encryption properties for the Databricks Workspace Managed Disks.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: keyvault.ValidateNestedItemID(keyvault.VersionTypeVersioned, keyvault.NestedItemTypeKey),
		},

		"managed_disk_cmk_rotation_to_latest_version_enabled": {
			Description:  "Whether customer managed keys for disk encryption will automatically be rotated to the latest version.",
			Type:         pluginsdk.TypeBool,
			Optional:     true,
			RequiredWith: []string{"managed_disk_cmk_key_vault_key_id"},
		},

		"disk_encryption_set_id": {
			Description: "The ID of Managed Disk Encryption Set created by the Databricks Workspace.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"storage_account_identity": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"principal_id": {
						Description: "The principal UUID for the internal databricks Storage Account needed to provide access to the workspace for enabling Customer Managed Keys.",
						Type:        pluginsdk.TypeString,
						Sensitive:   true,
						Computed:    true,
					},

					"tenant_id": {
						Description: "The UUID of the tenant where the internal databricks Storage Account was created.",
						Type:        pluginsdk.TypeString,
						Sensitive:   true,
						Computed:    true,
					},

					"type": {
						Description: "The type of the internal databricks Storage Account.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},
				},
			},
		},

		"enhanced_security_compliance": {
			Description: "This feature is only valid if `sku` is set to `premium`.",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			MaxItems:    1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"automatic_cluster_update_enabled": {
						Description: "Enables automatic cluster updates for this workspace.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     false,
					},
					"compliance_security_profile_enabled": {
						Description: "Enables compliance security profile for this workspace.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     false,
					},
					"compliance_security_profile_standards": {
						Description: "A list of standards to enforce on this workspace.",
						Type:        pluginsdk.TypeSet,
						Optional:    true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringInSlice(validate.PossibleValuesForComplianceStandard(), false),
						},
					},
					"enhanced_security_monitoring_enabled": {
						Description: "Enables enhanced security monitoring for this workspace.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     false,
					},
				},
			},
		},

		"tags": commonschema.Tags(),
	}
}
