// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package route

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/routes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func resourceRouteFlatten(d *pluginsdk.ResourceData, id *routes.RouteId, model *routes.Route) error {
	d.Set("name", id.RouteName)
	d.Set("route_table_name", id.RouteTableName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model != nil {
		if props := model.Properties; props != nil {
			d.Set("address_prefix", props.AddressPrefix)
			d.Set("next_hop_type", string(props.NextHopType))
			d.Set("next_hop_in_ip_address", props.NextHopIPAddress)
		}
	}

	return pluginsdk.SetResourceIdentityData(d, id)
}
