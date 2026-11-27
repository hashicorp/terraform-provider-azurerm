// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package virtual_hub_route_table_route

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/virtualwans"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/virtualwan/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceVirtualHubRouteTableRouteDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).VirtualWan.VirtualWANs
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.HubRouteTableRouteID(d.Id())
	if err != nil {
		return err
	}

	routeTableId := virtualwans.NewHubRouteTableID(id.SubscriptionId, id.ResourceGroup, id.VirtualHubName, id.HubRouteTableName)

	locks.ByName(id.VirtualHubName, virtualHubResourceName)
	defer locks.UnlockByName(id.VirtualHubName, virtualHubResourceName)

	// get latest list of routes
	routeTable, err := client.HubRouteTablesGet(ctx, routeTableId)
	if err != nil {
		if !response.WasNotFound(routeTable.HttpResponse) {
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if routeTable.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", routeTableId)
	}
	if routeTable.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", routeTableId)
	}

	props := routeTable.Model.Properties

	if props.Routes != nil {
		routes := *props.Routes

		newRoutes := make([]virtualwans.HubRoute, 0)
		for _, r := range routes {
			if r.Name != id.RouteName {
				newRoutes = append(newRoutes, r)
			}
		}
		props.Routes = &newRoutes
	}

	routeTable.Model.Properties = props

	if err := client.HubRouteTablesCreateOrUpdateThenPoll(ctx, routeTableId, *routeTable.Model); err != nil {
		return fmt.Errorf("removing %s: %+v", id, err)
	}

	return nil
}
