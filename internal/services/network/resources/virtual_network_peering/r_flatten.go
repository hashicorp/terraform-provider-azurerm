// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_network_peering

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/virtualnetworkpeerings"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func resourceVirtualNetworkPeeringFlatten(d *pluginsdk.ResourceData, id *virtualnetworkpeerings.VirtualNetworkPeeringId, model *virtualnetworkpeerings.VirtualNetworkPeering) error {
	d.Set("name", id.VirtualNetworkPeeringName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("virtual_network_name", id.VirtualNetworkName)

	if model != nil {
		if peer := model.Properties; peer != nil {
			d.Set("allow_virtual_network_access", peer.AllowVirtualNetworkAccess)
			d.Set("allow_forwarded_traffic", peer.AllowForwardedTraffic)
			d.Set("allow_gateway_transit", peer.AllowGatewayTransit)
			d.Set("peer_complete_virtual_networks_enabled", pointer.From(peer.PeerCompleteVnets))
			d.Set("only_ipv6_peering_enabled", pointer.From(peer.EnableOnlyIPv6Peering))
			d.Set("local_subnet_names", pointer.From(peer.LocalSubnetNames))
			d.Set("remote_subnet_names", pointer.From(peer.RemoteSubnetNames))
			d.Set("use_remote_gateways", peer.UseRemoteGateways)

			remoteVirtualNetworkId := ""
			if network := peer.RemoteVirtualNetwork; network != nil {
				parsed, err := commonids.ParseVirtualNetworkIDInsensitively(*network.Id)
				if err != nil {
					return err
				}
				remoteVirtualNetworkId = parsed.ID()
			}
			d.Set("remote_virtual_network_id", remoteVirtualNetworkId)
		}
	}

	return pluginsdk.SetResourceIdentityData(d, id)
}
