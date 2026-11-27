// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package route

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/routes"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/routetables"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	networkClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceRouteDelete(d *pluginsdk.ResourceData, meta any) error {
	client := networkClient.From(meta).Routes
	ctx, cancel := timeouts.ForDelete(account.StopContext(meta), d)
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
