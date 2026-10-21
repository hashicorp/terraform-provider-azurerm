// Copyright IBM Corp. 2014, 2025
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
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.RouteTableName,
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"route": {
			Type:       pluginsdk.TypeSet,
			ConfigMode: pluginsdk.SchemaConfigModeAttr,
			Optional:   true,
			Computed:   true, // azignore:AZS007 - pre-existing violation
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validate.RouteName,
					},

					"address_prefix": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"next_hop_type": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(routetables.PossibleValuesForRouteNextHopType(), false),
					},

					"next_hop_in_ip_address": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"bgp_route_propagation_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"subnets": {
			Type:     pluginsdk.TypeSet,
			Computed: true,
			Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
			Set:      pluginsdk.HashString,
		},

		"tags": commonschema.Tags(),
	}
}
