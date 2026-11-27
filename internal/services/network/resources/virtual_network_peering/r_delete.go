// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package virtual_network_peering

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/virtualnetworkpeerings"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	networkClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceVirtualNetworkPeeringDelete(d *pluginsdk.ResourceData, meta any) error {
	client := networkClient.From(meta).VirtualNetworkPeerings
	ctx, cancel := timeouts.ForDelete(account.StopContext(meta), d)
	defer cancel()

	id, err := virtualnetworkpeerings.ParseVirtualNetworkPeeringID(d.Id())
	if err != nil {
		return err
	}

	vnetId := commonids.NewVirtualNetworkID(id.SubscriptionId, id.ResourceGroupName, id.VirtualNetworkName)
	remoteVnetId, err := commonids.ParseVirtualNetworkID(d.Get("remote_virtual_network_id").(string))
	if err != nil {
		return err
	}

	vnetIDsToLock := []string{vnetId.ID(), remoteVnetId.ID()}
	locks.MultipleByID(&vnetIDsToLock)
	defer locks.UnlockMultipleByID(&vnetIDsToLock)

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return err
}
