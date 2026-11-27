// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package virtual_hub_route_table_route

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/virtualwans"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	virtualwanClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/virtualwan/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/virtualwan/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceVirtualHubRouteTableRouteUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := virtualwanClient.From(meta).VirtualWANs
	ctx, cancel := timeouts.ForUpdate(account.StopContext(meta), d)
	defer cancel()

	routeTableId, err := virtualwans.ParseHubRouteTableID(d.Get("route_table_id").(string))
	if err != nil {
		return err
	}

	locks.ByName(routeTableId.VirtualHubName, virtualHubResourceName)
	defer locks.UnlockByName(routeTableId.VirtualHubName, virtualHubResourceName)

	routeTable, err := client.HubRouteTablesGet(ctx, *routeTableId)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", routeTableId, err)
	}

	if routeTable.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", routeTableId)
	}
	if routeTable.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", routeTableId)
	}

	props := routeTable.Model.Properties

	id, err := parse.HubRouteTableRouteID(d.Id())
	if err != nil {
		return err
	}

	routes := *props.Routes
	for i := range routes {
		if routes[i].Name == id.RouteName {
			if d.HasChange("destinations_type") {
				routes[i].DestinationType = d.Get("destinations_type").(string)
			}
			if d.HasChange("destinations") {
				routes[i].Destinations = pointer.From(pluginsdk.ExpandStringSlice(d.Get("destinations").(*pluginsdk.Set).List()))
			}
			if d.HasChange("next_hop_type") {
				routes[i].NextHopType = d.Get("next_hop_type").(string)
			}
			if d.HasChange("next_hop") {
				routes[i].NextHop = d.Get("next_hop").(string)
			}
			break
		}
	}

	routeTable.Model.Properties.Routes = pointer.To(routes)

	if err := client.HubRouteTablesCreateOrUpdateThenPoll(ctx, *routeTableId, *routeTable.Model); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	d.SetId(id.ID())

	return resourceVirtualHubRouteTableRouteRead(d, meta)
}
