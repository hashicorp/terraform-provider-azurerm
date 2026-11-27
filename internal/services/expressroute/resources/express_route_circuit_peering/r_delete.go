// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package express_route_circuit_peering

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	expressrouteClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/expressroute/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceExpressRouteCircuitPeeringDelete(d *pluginsdk.ResourceData, meta any) error {
	client := expressrouteClient.From(meta).ExpressRouteCircuitPeerings
	ctx, cancel := timeouts.ForDelete(account.StopContext(meta), d)
	defer cancel()

	id, err := commonids.ParseExpressRouteCircuitPeeringID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.CircuitName, expressRouteCircuitResourceName)
	defer locks.UnlockByName(id.CircuitName, expressRouteCircuitResourceName)

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return err
}
