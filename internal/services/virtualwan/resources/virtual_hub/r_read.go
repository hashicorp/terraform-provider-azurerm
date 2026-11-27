// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package virtual_hub

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/virtualwans"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceVirtualHubRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).VirtualWan.VirtualWANs
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := virtualwans.ParseVirtualHubID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.VirtualHubsGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s does not exist - removing from state", id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.VirtualHubName)
	d.Set("resource_group_name", id.ResourceGroupName)

	defaultRouteTable := virtualwans.NewHubRouteTableID(id.SubscriptionId, id.ResourceGroupName, id.VirtualHubName, "defaultRouteTable")
	d.Set("default_route_table_id", defaultRouteTable.ID())

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))
		if props := model.Properties; props != nil {
			d.Set("address_prefix", props.AddressPrefix)
			d.Set("branch_to_branch_traffic_enabled", pointer.From(props.AllowBranchToBranchTraffic))
			d.Set("sku", props.Sku)

			if err := d.Set("route", flattenVirtualHubRoute(props.RouteTable)); err != nil {
				return fmt.Errorf("setting `route`: %+v", err)
			}

			d.Set("hub_routing_preference", pointer.FromEnum(props.HubRoutingPreference))

			var virtualWanId *string
			if props.VirtualWAN != nil {
				virtualWanId = props.VirtualWAN.Id
			}
			d.Set("virtual_wan_id", virtualWanId)

			var virtualRouterAsn *int64
			if props.VirtualRouterAsn != nil {
				virtualRouterAsn = props.VirtualRouterAsn
			}
			d.Set("virtual_router_asn", virtualRouterAsn)

			var virtualRouterIps *[]string
			if props.VirtualRouterIPs != nil {
				virtualRouterIps = props.VirtualRouterIPs
			}
			d.Set("virtual_router_ips", virtualRouterIps)

			d.Set("virtual_router_auto_scale_min_capacity", props.VirtualRouterAutoScaleConfiguration.MinCapacity)
		}
		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	return pluginsdk.SetResourceIdentityData(d, id)
}

func flattenVirtualHubRoute(input *virtualwans.VirtualHubRouteTable) []any {
	results := make([]any, 0)
	if input == nil || input.Routes == nil {
		return results
	}

	for _, item := range *input.Routes {
		addressPrefixes := pluginsdk.FlattenSlice(item.AddressPrefixes)

		results = append(results, map[string]any{
			"address_prefixes":    addressPrefixes,
			"next_hop_ip_address": pointer.From(item.NextHopIPAddress),
		})
	}

	return results
}

func virtualHubCreateRefreshFunc(ctx context.Context, client *virtualwans.VirtualWANsClient, id virtualwans.VirtualHubId) pluginsdk.StateRefreshFunc {
	return func() (any, string, error) {
		res, err := client.VirtualHubsGet(ctx, id)
		if err != nil {
			if response.WasNotFound(res.HttpResponse) {
				return nil, "", fmt.Errorf("%s doesn't exist", id)
			}
			return nil, "", fmt.Errorf("retrieving %s: %+v", id, err)
		}
		if res.Model == nil {
			return nil, "", fmt.Errorf("retrieving %s: `model` was nil", id)
		}
		if res.Model.Properties == nil {
			return nil, "", fmt.Errorf("retrieving %s: `properties` was nil", id)
		}

		state := res.Model.Properties.RoutingState
		if state != nil && *state == "Failed" {
			return nil, "", fmt.Errorf("failed to provision routing on %s", id)
		}
		return res, string(*res.Model.Properties.RoutingState), nil
	}
}
