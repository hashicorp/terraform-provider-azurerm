// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package route_table

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/routetables"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceRouteTableCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Network.RouteTables
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := routetables.NewRouteTableID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id, routetables.DefaultGetOperationOptions())
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_route_table", id.ID())
		}
	}

	bgpRoutePropagationEnabled := d.Get("bgp_route_propagation_enabled").(bool)

	locks.ByID(id.ID())
	defer locks.UnlockByID(id.ID())

	routeSet := routetables.RouteTable{
		Name:     &id.RouteTableName,
		Location: pointer.To(location.Normalize(d.Get("location").(string))),
		Properties: &routetables.RouteTablePropertiesFormat{
			Routes:                     expandRouteTableRoutes(d),
			DisableBgpRoutePropagation: pointer.To(!bgpRoutePropagationEnabled),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]any)),
	}

	if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, routeSet, sdk.SetIDAndIdentityCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
		return err
	}

	return resourceRouteTableRead(d, meta)
}

func expandRouteTableRoutes(d *pluginsdk.ResourceData) *[]routetables.Route {
	configs := d.Get("route").(*pluginsdk.Set).List()
	routes := make([]routetables.Route, 0, len(configs))

	for _, configRaw := range configs {
		data := configRaw.(map[string]any)

		route := routetables.Route{
			Name: pointer.To(data["name"].(string)),
			Properties: &routetables.RoutePropertiesFormat{
				AddressPrefix: pointer.To(data["address_prefix"].(string)),
				NextHopType:   routetables.RouteNextHopType(data["next_hop_type"].(string)),
			},
		}

		if v := data["next_hop_in_ip_address"].(string); v != "" {
			route.Properties.NextHopIPAddress = &v
		}

		routes = append(routes, route)
	}

	return &routes
}
