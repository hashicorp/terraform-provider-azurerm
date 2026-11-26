// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package express_route_circuit_authorization

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/expressroutecircuitauthorizations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceExpressRouteCircuitAuthorizationRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).ExpressRoute.ExpressRouteCircuitAuthorizations
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := expressroutecircuitauthorizations.ParseAuthorizationID(d.Id())
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

	d.Set("name", id.AuthorizationName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("express_route_circuit_name", id.ExpressRouteCircuitName)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("authorization_key", props.AuthorizationKey)
			d.Set("authorization_use_status", pointer.FromEnum(props.AuthorizationUseStatus))
		}
	}

	return nil
}
