// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package eventhub_namespace

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2026-01-01/clusters"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2026-01-01/ehnamespaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2026-01-01/networkrulesets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/eventhub/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func eventhubNamespaceSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the EventHub Namespace resource.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ValidateEventHubNamespaceName(),
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"sku": {
			Description:  "Which tier to use. Setting this field to `Premium` will force the creation of a new resource.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(ehnamespaces.PossibleValuesForSkuName(), false),
		},

		"capacity": {
			Description: "The Capacity / Throughput Units for a `Standard` SKU namespace.",
			Type:        pluginsdk.TypeInt,
			Optional:    true,
			Default:     1,
		},

		"auto_inflate_enabled": {
			Description: "Whether Auto Inflate is enabled for the EventHub Namespace.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"dedicated_cluster_id": {
			Description:  "The ID of the EventHub Dedicated Cluster where this Namespace should be created.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: clusters.ValidateClusterID,
		},

		"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

		"maximum_throughput_units": {
			Description:  "The maximum number of throughput units when Auto Inflate is Enabled.",
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntBetween(0, 40),
		},

		"network_rulesets": {
			Type:       pluginsdk.TypeList,
			Optional:   true,
			Computed:   true, // azignore:AZS007 - pre-existing violation
			MaxItems:   1,
			ConfigMode: pluginsdk.SchemaConfigModeAttr,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"default_action": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(networkrulesets.PossibleValuesForDefaultAction(), false),
					},

					"public_network_access_enabled": {
						Description: "Whether public network access is enabled for the EventHub Namespace.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     true,
					},

					"trusted_service_access_enabled": {
						Description: "Whether Trusted Microsoft Services are allowed to bypass firewall.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
					},

					// Returned value of the `virtual_network_rule` array does not honour the input order,
					// possibly a service design, thus changed to TypeSet
					"virtual_network_rule": {
						Type:       pluginsdk.TypeSet,
						Optional:   true,
						ConfigMode: pluginsdk.SchemaConfigModeAttr,
						Set:        resourceVnetRuleHash,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								// the API returns the subnet ID's resource group name in lowercase
								// https://github.com/Azure/azure-sdk-for-go/issues/5855
								"subnet_id": {
									Description:      "The ID of the subnet to match on.",
									Type:             pluginsdk.TypeString,
									Required:         true,
									ValidateFunc:     commonids.ValidateSubnetID,
									DiffSuppressFunc: suppress.CaseDifference,
								},

								"ignore_missing_virtual_network_service_endpoint": {
									Description: "Whether missing Virtual Network service endpoints are ignored.",
									Type:        pluginsdk.TypeBool,
									Optional:    true,
								},
							},
						},
					},

					"ip_rule": {
						Type:       pluginsdk.TypeList,
						Optional:   true,
						ConfigMode: pluginsdk.SchemaConfigModeAttr,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"ip_mask": {
									Description: "The IP mask to match on.",
									Type:        pluginsdk.TypeString,
									Required:    true,
								},

								"action": {
									Description:  "The action to take when the rule is matched.",
									Type:         pluginsdk.TypeString,
									Optional:     true,
									Default:      string(networkrulesets.NetworkRuleIPActionAllow),
									ValidateFunc: validation.StringInSlice(networkrulesets.PossibleValuesForNetworkRuleIPAction(), false),
								},
							},
						},
					},
				},
			},
		},

		"local_authentication_enabled": {
			Description: "Whether SAS authentication is enabled for the EventHub Namespace.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"minimum_tls_version": {
			Description: "The minimum supported TLS version for this EventHub Namespace.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			Default:     string(ehnamespaces.TlsVersionOnePointTwo),
			ValidateFunc: validation.StringInSlice([]string{
				string(ehnamespaces.TlsVersionOnePointTwo),
			}, false),
		},

		"public_network_access_enabled": {
			Description: "Whether public network access is enabled for the EventHub Namespace.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"default_primary_connection_string_alias": {
			Description: "The alias of the primary connection string for the authorization rule `RootManageSharedAccessKey`, which is generated when disaster recovery is enabled.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Sensitive:   true,
		},

		"default_secondary_connection_string_alias": {
			Description: "The alias of the secondary connection string for the authorization rule `RootManageSharedAccessKey`, which is generated when disaster recovery is enabled.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Sensitive:   true,
		},

		"default_primary_connection_string": {
			Description: "The primary connection string for the authorization rule `RootManageSharedAccessKey`.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Sensitive:   true,
		},

		"default_primary_key": {
			Description: "The primary access key for the authorization rule `RootManageSharedAccessKey`.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Sensitive:   true,
		},

		"default_secondary_connection_string": {
			Description: "The secondary connection string for the authorization rule `RootManageSharedAccessKey`.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Sensitive:   true,
		},

		"default_secondary_key": {
			Description: "The secondary access key for the authorization rule `RootManageSharedAccessKey`.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Sensitive:   true,
		},

		"tags": commonschema.Tags(),
	}
}
