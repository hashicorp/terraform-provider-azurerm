// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package express_route_connection

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/expressrouteconnections"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceExpressRouteConnectionDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).VirtualWan.ExpressRouteConnections
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := expressrouteconnections.ParseExpressRouteConnectionID(d.Id())
	if err != nil {
		return err
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
