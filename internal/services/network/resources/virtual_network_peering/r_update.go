// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_network_peering

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/virtualnetworkpeerings"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceVirtualNetworkPeeringUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Network.VirtualNetworkPeerings
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
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

	existing, err := client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	if existing.Model == nil || existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", *id)
	}

	if d.HasChange("allow_forwarded_traffic") {
		existing.Model.Properties.AllowForwardedTraffic = pointer.To(d.Get("allow_forwarded_traffic").(bool))
	}
	if d.HasChange("allow_gateway_transit") {
		existing.Model.Properties.AllowGatewayTransit = pointer.To(d.Get("allow_gateway_transit").(bool))
	}
	if d.HasChange("allow_virtual_network_access") {
		existing.Model.Properties.AllowVirtualNetworkAccess = pointer.To(d.Get("allow_virtual_network_access").(bool))
	}
	if d.HasChange("local_subnet_names") {
		existing.Model.Properties.LocalSubnetNames = pluginsdk.ExpandStringSlice(d.Get("local_subnet_names").([]any))
	}
	if d.HasChange("remote_subnet_names") {
		existing.Model.Properties.RemoteSubnetNames = pluginsdk.ExpandStringSlice(d.Get("remote_subnet_names").([]any))
	}
	if d.HasChange("use_remote_gateways") {
		existing.Model.Properties.UseRemoteGateways = pointer.To(d.Get("use_remote_gateways").(bool))
	}
	if d.HasChange("remote_virtual_network_id") {
		existing.Model.Properties.RemoteVirtualNetwork = &virtualnetworkpeerings.SubResource{
			Id: pointer.To(d.Get("remote_virtual_network_id").(string)),
		}
	}

	if err := client.CreateOrUpdateThenPoll(ctx, *id, *existing.Model, virtualnetworkpeerings.CreateOrUpdateOperationOptions{SyncRemoteAddressSpace: pointer.To(virtualnetworkpeerings.SyncRemoteAddressSpaceTrue)}); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	return resourceVirtualNetworkPeeringRead(d, meta)
}
