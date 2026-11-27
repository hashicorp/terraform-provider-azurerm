// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package route_table

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/routetables"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func routeTableSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the Route Table.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.RouteTableName,
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"route": {
			Type:       pluginsdk.TypeSet,
			Optional:   true,
			Computed:   true, // azignore:AZS007 - pre-existing violation
			ConfigMode: pluginsdk.SchemaConfigModeAttr,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Description:  "The name of the route.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validate.RouteName,
					},

					"address_prefix": {
						Description:  "The destination to which the route applies. Can be CIDR (such as `10.1.0.0/16`) or [Azure Service Tag](https://learn.microsoft.com/azure/virtual-network/service-tags-overview) (such as `ApiManagement`, `AzureBackup` or `AzureMonitor`) format.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"next_hop_type": {
						Description:  "The type of Azure hop the packet should be sent to.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(routetables.PossibleValuesForRouteNextHopType(), false),
					},

					"next_hop_in_ip_address": {
						Description:  "Contains the IP address packets should be forwarded to. Next hop values are only allowed in routes where the next hop type is `VirtualAppliance`.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"bgp_route_propagation_enabled": {
			Description: "Boolean flag which controls propagation of routes learned by BGP on that Route Table.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"subnets": {
			Description: "The collection of Subnets associated with this Route Table.",
			Type:        pluginsdk.TypeSet,
			Computed:    true,
			Set:         pluginsdk.HashString,
			Elem:        &pluginsdk.Schema{Type: pluginsdk.TypeString},
		},

		"tags": commonschema.Tags(),
	}
}
