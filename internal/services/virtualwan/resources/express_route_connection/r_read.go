// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package express_route_connection

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/expressrouteconnections"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/expressroutegateways"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceExpressRouteConnectionRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).VirtualWan.ExpressRouteConnections
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := expressrouteconnections.ParseExpressRouteConnectionID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.ExpressRouteConnectionName)
	d.Set("express_route_gateway_id", expressroutegateways.NewExpressRouteGatewayID(id.SubscriptionId, id.ResourceGroupName, id.ExpressRouteGatewayName).ID())

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("routing_weight", props.RoutingWeight)
			d.Set("authorization_key", props.AuthorizationKey)

			d.Set("internet_security_enabled", props.EnableInternetSecurity)

			if props.ExpressRouteGatewayBypass != nil {
				d.Set("express_route_gateway_bypass_enabled", props.ExpressRouteGatewayBypass)
			}

			circuitPeeringID := pointer.From(props.ExpressRouteCircuitPeering.Id)
			peeringId, err := commonids.ParseExpressRouteCircuitPeeringIDInsensitively(circuitPeeringID)
			if err != nil {
				return err
			}
			d.Set("express_route_circuit_peering_id", peeringId.ID())

			routing, err := flattenExpressRouteConnectionRouting(props.RoutingConfiguration)
			if err != nil {
				return err
			}
			if err := d.Set("routing", routing); err != nil {
				return fmt.Errorf("setting `routing`: %+v", err)
			}
		}
	}

	return nil
}
