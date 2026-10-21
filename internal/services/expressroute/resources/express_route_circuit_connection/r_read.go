// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package express_route_circuit_connection

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/expressroutecircuitconnections"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceExpressRouteCircuitConnectionRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).ExpressRoute.ExpressRouteCircuitConnections
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := expressroutecircuitconnections.ParsePeeringConnectionID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s does not exist - removing from state", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.ConnectionName)
	d.Set("peering_id", commonids.NewExpressRouteCircuitPeeringID(id.SubscriptionId, id.ResourceGroupName, id.ExpressRouteCircuitName, id.PeeringName).ID())

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("address_prefix_ipv4", props.AddressPrefix)

			// The ExpressRoute Circuit Connection API returns "*****************" for AuthorizationKey when it's changed from a valid value to `nil`
			// See more details from https://github.com/Azure/azure-rest-api-specs/issues/15030
			authorizationKey := ""
			if props.AuthorizationKey != nil && *props.AuthorizationKey != "*****************" {
				authorizationKey = *props.AuthorizationKey
			}
			d.Set("authorization_key", authorizationKey)

			addressPrefixIPv6 := ""
			if props.IPv6CircuitConnectionConfig != nil && props.IPv6CircuitConnectionConfig.AddressPrefix != nil {
				addressPrefixIPv6 = *props.IPv6CircuitConnectionConfig.AddressPrefix
			}
			d.Set("address_prefix_ipv6", addressPrefixIPv6)

			if props.PeerExpressRouteCircuitPeering != nil && props.PeerExpressRouteCircuitPeering.Id != nil {
				circuitPeerPeeringId, err := commonids.ParseExpressRouteCircuitPeeringIDInsensitively(*props.PeerExpressRouteCircuitPeering.Id)
				if err != nil {
					return err
				}
				d.Set("peer_peering_id", circuitPeerPeeringId.ID())
			}
		}
	}

	return nil
}
