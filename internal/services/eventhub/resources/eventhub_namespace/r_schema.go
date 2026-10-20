// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package eventhub_namespace

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2024-01-01/eventhubsclusters"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2024-01-01/namespaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2024-01-01/networkrulesets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/eventhub/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func eventhubNamespaceSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ValidateEventHubNamespaceName(),
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"sku": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(namespaces.PossibleValuesForSkuName(), false),
		},

		"capacity": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			Default:  1,
		},

		"auto_inflate_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"dedicated_cluster_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: eventhubsclusters.ValidateClusterID,
		},

		"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

		"maximum_throughput_units": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntBetween(0, 40),
		},

		"network_rulesets": {
			Type:       pluginsdk.TypeList,
			Optional:   true,
			MaxItems:   1,
			Computed:   true, // azignore:AZS007 - pre-existing violation
			ConfigMode: pluginsdk.SchemaConfigModeAttr,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"default_action": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(networkrulesets.PossibleValuesForDefaultAction(), false),
					},

					"public_network_access_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  true,
					},

					"trusted_service_access_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
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
									Type:             pluginsdk.TypeString,
									Required:         true,
									ValidateFunc:     commonids.ValidateSubnetID,
									DiffSuppressFunc: suppress.CaseDifference,
								},

								"ignore_missing_virtual_network_service_endpoint": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
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
									Type:     pluginsdk.TypeString,
									Required: true,
								},

								"action": {
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
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"minimum_tls_version": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  string(namespaces.TlsVersionOnePointTwo),
			ValidateFunc: validation.StringInSlice([]string{
				string(namespaces.TlsVersionOnePointTwo),
			}, false),
		},

		"public_network_access_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"default_primary_connection_string_alias": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},

		"default_secondary_connection_string_alias": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},

		"default_primary_connection_string": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},

		"default_primary_key": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},

		"default_secondary_connection_string": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},

		"default_secondary_key": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},

		"tags": commonschema.Tags(),
	}
}
