// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package route_table

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/routetables"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func resourceRouteTableFlatten(d *pluginsdk.ResourceData, id *routetables.RouteTableId, table *routetables.RouteTable) error {
	d.Set("name", id.RouteTableName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if table != nil {
		d.Set("location", location.NormalizeNilable(table.Location))

		if props := table.Properties; props != nil {
			d.Set("bgp_route_propagation_enabled", !pointer.From(props.DisableBgpRoutePropagation))
			if err := d.Set("route", flattenRouteTableRoutes(props.Routes)); err != nil {
				return err
			}

			if err := d.Set("subnets", flattenRouteTableSubnets(props.Subnets)); err != nil {
				return err
			}
		}

		if err := tags.FlattenAndSet(d, table.Tags); err != nil {
			return err
		}
	}

	return pluginsdk.SetResourceIdentityData(d, id)
}

func flattenRouteTableRoutes(input *[]routetables.Route) []any {
	results := make([]any, 0)

	if routes := input; routes != nil {
		for _, route := range *routes {
			r := make(map[string]any)

			r["name"] = *route.Name

			if props := route.Properties; props != nil {
				r["address_prefix"] = *props.AddressPrefix
				r["next_hop_type"] = string(props.NextHopType)
				if ip := props.NextHopIPAddress; ip != nil && *ip != "" {
					r["next_hop_in_ip_address"] = *ip
				}
			}

			results = append(results, r)
		}
	}

	return results
}

func flattenRouteTableSubnets(subnets *[]routetables.Subnet) []string {
	output := make([]string, 0)

	if subnets != nil {
		for _, subnet := range *subnets {
			output = append(output, *subnet.Id)
		}
	}

	return output
}
