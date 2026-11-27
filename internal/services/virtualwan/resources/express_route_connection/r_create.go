// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package express_route_connection

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/expressrouteconnections"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/expressroutegateways"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceExpressRouteConnectionCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).VirtualWan.ExpressRouteConnections
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	expressRouteGatewayId, err := expressroutegateways.ParseExpressRouteGatewayID(d.Get("express_route_gateway_id").(string))
	if err != nil {
		return err
	}

	id := expressrouteconnections.NewExpressRouteConnectionID(expressRouteGatewayId.SubscriptionId, expressRouteGatewayId.ResourceGroupName, expressRouteGatewayId.ExpressRouteGatewayName, d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_express_route_connection", id.ID())
		}
	}

	parameters := expressrouteconnections.ExpressRouteConnection{
		Name: id.ExpressRouteConnectionName,
		Properties: &expressrouteconnections.ExpressRouteConnectionProperties{
			ExpressRouteCircuitPeering: expressrouteconnections.ExpressRouteCircuitPeeringId{
				Id: pointer.To(d.Get("express_route_circuit_peering_id").(string)),
			},
			EnableInternetSecurity:    pointer.To(d.Get("internet_security_enabled").(bool)),
			RoutingConfiguration:      expandExpressRouteConnectionRouting(d.Get("routing").([]any)),
			RoutingWeight:             pointer.To(int64(d.Get("routing_weight").(int))),
			ExpressRouteGatewayBypass: pointer.To(d.Get("express_route_gateway_bypass_enabled").(bool)),
		},
	}

	if v, ok := d.GetOk("authorization_key"); ok {
		parameters.Properties.AuthorizationKey = pointer.To(v.(string))
	}

	if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, parameters, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceExpressRouteConnectionRead(d, meta)
}
