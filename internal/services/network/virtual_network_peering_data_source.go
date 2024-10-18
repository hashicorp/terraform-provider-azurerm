package network

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/virtualnetworkpeerings"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceVirtualNetworkPeering() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceVirtualNetworkPeeringRead,
		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"virtual_network_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"remote_virtual_network_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"allow_virtual_network_access": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"allow_forwarded_traffic": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"allow_gateway_transit": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"only_ipv6_peering_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"peer_complete_virtual_networks_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"use_remote_gateways": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceVirtualNetworkPeeringRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VirtualNetworkPeerings
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := virtualnetworkpeerings.NewVirtualNetworkPeeringID(subscriptionId, d.Get("resource_group_name").(string), d.Get("virtual_network_name").(string), d.Get("name").(string))

	resp, err := client.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())

	d.Set("name", id.VirtualNetworkPeeringName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("virtual_network_name", id.VirtualNetworkName)

	if model := resp.Model; model != nil {
		if peer := model.Properties; peer != nil {
			d.Set("allow_virtual_network_access", peer.AllowVirtualNetworkAccess)
			d.Set("allow_forwarded_traffic", peer.AllowForwardedTraffic)
			d.Set("allow_gateway_transit", peer.AllowGatewayTransit)
			d.Set("peer_complete_virtual_networks_enabled", peer.PeerCompleteVnets)
			d.Set("only_ipv6_peering_enabled", peer.EnableOnlyIPv6Peering)
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

	return nil
}
