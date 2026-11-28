// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package virtual_hub_route_table

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/virtualwans"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceVirtualHubRouteTableRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).VirtualWan.VirtualWANs
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := virtualwans.ParseHubRouteTableID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.HubRouteTablesGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s does not exist - removing from state", id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.HubRouteTableName)
	d.Set("virtual_hub_id", virtualwans.NewVirtualHubID(id.SubscriptionId, id.ResourceGroupName, id.VirtualHubName).ID())

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("labels", pluginsdk.FlattenSlice(props.Labels))

			if err := d.Set("route", flattenVirtualHubRouteTableHubRoutes(props.Routes)); err != nil {
				return fmt.Errorf("setting `route`: %+v", err)
			}
		}
	}

	return pluginsdk.SetResourceIdentityData(d, id)
}

func flattenVirtualHubRouteTableHubRoutes(input *[]virtualwans.HubRoute) []any {
	results := make([]any, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		v := map[string]any{
			"name":              item.Name,
			"destinations":      pluginsdk.FlattenSlice(&item.Destinations),
			"destinations_type": item.DestinationType,
			"next_hop":          item.NextHop,
			"next_hop_type":     item.NextHopType,
		}

		results = append(results, v)
	}

	return results
}
