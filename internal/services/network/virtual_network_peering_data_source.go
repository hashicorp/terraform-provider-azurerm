package network

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/virtualnetworkpeerings"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.DataSource = VirtualNetworkPeeringDataSource{}

type VirtualNetworkPeeringDataSource struct{}

func (VirtualNetworkPeeringDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
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
	}
}

func (VirtualNetworkPeeringDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
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
	}
}

func (VirtualNetworkPeeringDataSource) ModelObject() interface{} {
	return nil
}

func (VirtualNetworkPeeringDataSource) ResourceType() string {
	return "azurerm_virtual_network_peering"
}

func (VirtualNetworkPeeringDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.VirtualNetworkPeerings

			subscriptionId := metadata.Client.Account.SubscriptionId
			name := metadata.ResourceData.Get("name").(string)
			resource_group_name := metadata.ResourceData.Get("resource_group_name").(string)
			virtual_network_name := metadata.ResourceData.Get("virtual_network_name").(string)
			id := virtualnetworkpeerings.NewVirtualNetworkPeeringID(subscriptionId, resource_group_name, virtual_network_name, name)

			resp, err := client.Get(ctx, id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			metadata.SetID(id)

			metadata.ResourceData.Set("name", id.VirtualNetworkPeeringName)
			metadata.ResourceData.Set("resource_group_name", id.ResourceGroupName)
			metadata.ResourceData.Set("virtual_network_name", id.VirtualNetworkName)

			if model := resp.Model; model != nil {
				if peer := model.Properties; peer != nil {
					metadata.ResourceData.Set("allow_virtual_network_access", peer.AllowVirtualNetworkAccess)
					metadata.ResourceData.Set("allow_forwarded_traffic", peer.AllowForwardedTraffic)
					metadata.ResourceData.Set("allow_gateway_transit", peer.AllowGatewayTransit)
					metadata.ResourceData.Set("peer_complete_virtual_networks_enabled", peer.PeerCompleteVnets)
					metadata.ResourceData.Set("only_ipv6_peering_enabled", peer.EnableOnlyIPv6Peering)
					metadata.ResourceData.Set("use_remote_gateways", peer.UseRemoteGateways)

					remoteVirtualNetworkId := ""
					if network := peer.RemoteVirtualNetwork; network != nil {
						parsed, err := commonids.ParseVirtualNetworkIDInsensitively(*network.Id)
						if err != nil {
							return err
						}
						remoteVirtualNetworkId = parsed.ID()
					}
					metadata.ResourceData.Set("remote_virtual_network_id", remoteVirtualNetworkId)
				}
			}
			return nil
		},
	}
}
