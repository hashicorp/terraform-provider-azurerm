// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package express_route_port

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/expressrouteports"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceArmExpressRoutePortDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).ExpressRoute.ExpressRoutePorts
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := expressrouteports.ParseExpressRoutePortID(d.Id())
	if err != nil {
		return err
	}

	// a lock is needed here for subresource express_route_port_authorization needs a lock.
	locks.ByID(id.ID())
	defer locks.UnlockByID(id.ID())

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}
