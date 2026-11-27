// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package virtual_hub_route_table_route

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/virtualwans"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	virtualwanClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/virtualwan/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/virtualwan/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceVirtualHubRouteTableRouteRead(d *pluginsdk.ResourceData, meta any) error {
	client := virtualwanClient.From(meta).VirtualWANs
	ctx, cancel := timeouts.ForRead(account.StopContext(meta), d)
	defer cancel()

	id, err := parse.HubRouteTableRouteID(d.Id())
	if err != nil {
		return err
	}

	routeTableId := virtualwans.NewHubRouteTableID(id.SubscriptionId, id.ResourceGroup, id.VirtualHubName, id.HubRouteTableName)

	resp, err := client.HubRouteTablesGet(ctx, routeTableId)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s does not exist - removing from state", routeTableId.ID())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.RouteName)

	routeTableID := virtualwans.NewHubRouteTableID(id.SubscriptionId, id.ResourceGroup, id.VirtualHubName, id.HubRouteTableName)
	d.Set("route_table_id", routeTableID.ID())

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			found := false
			for _, r := range *props.Routes {
				if r.Name == id.RouteName {
					found = true
					d.Set("destinations_type", r.DestinationType)
					d.Set("destinations", r.Destinations)
					d.Set("next_hop_type", r.NextHopType)
					d.Set("next_hop", r.NextHop)
					break
				}
			}

			if !found {
				// could not find existing id by name
				d.SetId("")
				return nil
			}
		}
	}

	return nil
}
