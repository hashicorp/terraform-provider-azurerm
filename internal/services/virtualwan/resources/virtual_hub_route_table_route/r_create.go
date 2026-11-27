// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package virtual_hub_route_table_route

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/virtualwans"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/virtualwan/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceVirtualHubRouteTableRouteCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).VirtualWan.VirtualWANs
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	routeTableId, err := virtualwans.ParseHubRouteTableID(d.Get("route_table_id").(string))
	if err != nil {
		return err
	}

	locks.ByName(routeTableId.VirtualHubName, virtualHubResourceName)
	defer locks.UnlockByName(routeTableId.VirtualHubName, virtualHubResourceName)

	routeTable, err := client.HubRouteTablesGet(ctx, *routeTableId)
	if err != nil {
		if !response.WasNotFound(routeTable.HttpResponse) {
			return fmt.Errorf("checking for existing %s: %+v", routeTableId, err)
		}
		return fmt.Errorf("retrieving %s: %+v", routeTableId, err)
	}

	if routeTable.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", routeTableId)
	}
	if routeTable.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", routeTableId)
	}

	props := routeTable.Model.Properties

	name := d.Get("name").(string)
	id := parse.NewHubRouteTableRouteID(routeTableId.SubscriptionId, routeTableId.ResourceGroupName, routeTableId.VirtualHubName, routeTableId.HubRouteTableName, name)

	routes := make([]virtualwans.HubRoute, 0)
	exists := false
	if hubRoutes := props.Routes; hubRoutes != nil {
		for _, r := range *hubRoutes {
			if r.Name == name {
				exists = true
				if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
					return tf.ImportAsExistsError("azurerm_virtual_hub_route_table_route", id.ID())
				}
			}
		}
		routes = *props.Routes

		result := virtualwans.HubRoute{
			Name:            d.Get("name").(string),
			DestinationType: d.Get("destinations_type").(string),
			Destinations:    pointer.From(pluginsdk.ExpandStringSlice(d.Get("destinations").(*pluginsdk.Set).List())),
			NextHopType:     d.Get("next_hop_type").(string),
			NextHop:         d.Get("next_hop").(string),
		}

		if !exists {
			routes = append(routes, result)
		}
	}

	routeTable.Model.Properties.Routes = pointer.To(routes)

	// // TODO: implement `CallbackThenPoll`, requires migrating to an ID that implements `resourceids.ResourceId`
	if err := client.HubRouteTablesCreateOrUpdateThenPoll(ctx, *routeTableId, *routeTable.Model); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceVirtualHubRouteTableRouteRead(d, meta)
}
