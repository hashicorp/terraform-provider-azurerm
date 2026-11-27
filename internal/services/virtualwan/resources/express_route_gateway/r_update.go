// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package express_route_gateway

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/expressrouteconnections"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/expressroutegateways"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	virtualwanClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/virtualwan/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceExpressRouteGatewayUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := virtualwanClient.From(meta).ExpressRouteGateways
	connectionsClient := virtualwanClient.From(meta).ExpressRouteConnections
	ctx, cancel := timeouts.ForUpdate(account.StopContext(meta), d)
	defer cancel()

	id, err := expressroutegateways.ParseExpressRouteGatewayID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", id)
	}
	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", id)
	}

	gatewayId, err := expressrouteconnections.ParseExpressRouteGatewayID(d.Id())
	if err != nil {
		return err
	}

	respConnections, err := connectionsClient.List(ctx, *gatewayId)
	if err != nil && !response.WasNotFound(respConnections.HttpResponse) {
		return fmt.Errorf("retrieving %s: %+v", gatewayId, err)
	}

	payload := existing.Model

	var connections *[]expressroutegateways.ExpressRouteConnection
	if model := respConnections.Model; model != nil {
		connections = convertConnectionsToGatewayConnections(model.Value)
	}

	payload.Properties.ExpressRouteConnections = connections

	if d.HasChange("scale_units") {
		payload.Properties.AutoScaleConfiguration = &expressroutegateways.ExpressRouteGatewayPropertiesAutoScaleConfiguration{
			Bounds: &expressroutegateways.ExpressRouteGatewayPropertiesAutoScaleConfigurationBounds{
				Min: pointer.To(int64(d.Get("scale_units").(int))),
			},
		}
	}

	if d.HasChange("allow_non_virtual_wan_traffic") {
		payload.Properties.AllowNonVirtualWanTraffic = pointer.To(d.Get("allow_non_virtual_wan_traffic").(bool))
	}

	if d.HasChange("tags") {
		payload.Tags = tags.Expand(d.Get("tags").(map[string]any))
	}

	if err := client.CreateOrUpdateThenPoll(ctx, *id, *payload); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceExpressRouteGatewayRead(d, meta)
}
