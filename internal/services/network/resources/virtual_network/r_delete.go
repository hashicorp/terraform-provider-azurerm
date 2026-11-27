// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package virtual_network

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	networkClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceVirtualNetworkDelete(d *pluginsdk.ResourceData, meta any) error {
	client := networkClient.From(meta).VirtualNetworks
	ctx, cancel := timeouts.ForDelete(account.StopContext(meta), d)
	defer cancel()

	id, err := commonids.ParseVirtualNetworkID(d.Id())
	if err != nil {
		return err
	}

	nsgNames, routeTableNames, err := expandResourcesForLocking(d)
	if err != nil {
		return fmt.Errorf("parsing Network Security Group ID's: %+v", err)
	}

	locks.ByID(id.ID())
	defer locks.UnlockByID(id.ID())

	locks.MultipleByID(&nsgNames)
	defer locks.UnlockMultipleByID(&nsgNames)

	locks.MultipleByID(&routeTableNames)
	defer locks.UnlockMultipleByID(&routeTableNames)

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
