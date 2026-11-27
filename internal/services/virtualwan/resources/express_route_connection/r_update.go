// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package express_route_connection

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/expressrouteconnections"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceExpressRouteConnectionUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).VirtualWan.ExpressRouteConnections
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := expressrouteconnections.ParseExpressRouteConnectionID(d.Id())
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
	props := existing.Model.Properties

	if d.HasChange("authorization_key") {
		props.AuthorizationKey = nil
		if authKey := d.Get("authorization_key").(string); authKey != "" {
			props.AuthorizationKey = pointer.To(authKey)
		}
	}

	if d.HasChange("internet_security_enabled") {
		props.EnableInternetSecurity = pointer.To(d.Get("internet_security_enabled").(bool))
	}

	if d.HasChange("express_route_gateway_bypass_enabled") {
		props.ExpressRouteGatewayBypass = pointer.To(d.Get("express_route_gateway_bypass_enabled").(bool))
	}

	if d.HasChange("routing") {
		props.RoutingConfiguration = expandExpressRouteConnectionRouting(d.Get("routing").([]any))
	}

	if d.HasChange("routing_weight") {
		props.RoutingWeight = pointer.To(int64(d.Get("routing_weight").(int)))
	}

	if err := client.CreateOrUpdateThenPoll(ctx, *id, *existing.Model); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceExpressRouteConnectionRead(d, meta)
}
