// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package route

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/routes"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/routetables"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceRouteDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Network.Routes
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := routes.ParseRouteID(d.Id())
	if err != nil {
		return err
	}

	rtId := routetables.NewRouteTableID(id.SubscriptionId, id.ResourceGroupName, id.RouteTableName)
	locks.ByID(rtId.ID())
	defer locks.UnlockByID(rtId.ID())

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
