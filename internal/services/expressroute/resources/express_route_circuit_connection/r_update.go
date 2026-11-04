// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package express_route_circuit_connection

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/expressroutecircuitconnections"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/expressroutecircuits"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceExpressRouteCircuitConnectionUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).ExpressRoute.ExpressRouteCircuitConnections
	circuitClient := meta.(*clients.Client).ExpressRoute.ExpressRouteCircuits
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := expressroutecircuitconnections.ParsePeeringConnectionID(d.Id())
	if err != nil {
		return err
	}

	circuitPeeringId, err := commonids.ParseExpressRouteCircuitPeeringID(d.Get("peering_id").(string))
	if err != nil {
		return err
	}

	circuitPeerPeeringId, err := commonids.ParseExpressRouteCircuitPeeringID(d.Get("peer_peering_id").(string))
	if err != nil {
		return err
	}

	expressRouteCircuitConnectionParameters := expressroutecircuitconnections.ExpressRouteCircuitConnection{
		Name: pointer.To(id.ConnectionName),
		Properties: &expressroutecircuitconnections.ExpressRouteCircuitConnectionPropertiesFormat{
			AddressPrefix: pointer.To(d.Get("address_prefix_ipv4").(string)),
			ExpressRouteCircuitPeering: &expressroutecircuitconnections.SubResource{
				Id: pointer.To(circuitPeeringId.ID()),
			},
			PeerExpressRouteCircuitPeering: &expressroutecircuitconnections.SubResource{
				Id: pointer.To(circuitPeerPeeringId.ID()),
			},
		},
	}

	if v, ok := d.GetOk("authorization_key"); ok {
		expressRouteCircuitConnectionParameters.Properties.AuthorizationKey = pointer.To(v.(string))
	}

	if v, ok := d.GetOk("address_prefix_ipv6"); ok {
		circuitId := expressroutecircuits.NewExpressRouteCircuitID(circuitPeeringId.SubscriptionId, circuitPeeringId.ResourceGroupName, circuitPeeringId.CircuitName)

		circuit, err := circuitClient.Get(ctx, circuitId)
		if err != nil {
			return fmt.Errorf("retrieving %s: %+v", circuitId, err)
		}

		if circuit.Model != nil && circuit.Model.Properties != nil && circuit.Model.Properties.ExpressRoutePort != nil {
			return fmt.Errorf("`address_prefix_ipv6` cannot be set when ExpressRoute Circuit Connection with ExpressRoute Circuit based on ExpressRoute Port")
		} else {
			expressRouteCircuitConnectionParameters.Properties.IPv6CircuitConnectionConfig = &expressroutecircuitconnections.IPv6CircuitConnectionConfig{
				AddressPrefix: pointer.To(v.(string)),
			}
		}
	}

	if err := client.CreateOrUpdateThenPoll(ctx, *id, expressRouteCircuitConnectionParameters); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceExpressRouteCircuitConnectionRead(d, meta)
}
