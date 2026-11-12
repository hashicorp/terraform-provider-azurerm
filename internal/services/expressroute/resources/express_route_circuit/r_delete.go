// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package express_route_circuit

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/expressroutecircuits"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceExpressRouteCircuitDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).ExpressRoute.ExpressRouteCircuits
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := expressroutecircuits.ParseExpressRouteCircuitID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.ExpressRouteCircuitName, expressRouteCircuitResourceName)
	defer locks.UnlockByName(id.ExpressRouteCircuitName, expressRouteCircuitResourceName)

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s : %+v", *id, err)
	}

	return err
}
