// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2024-05-01/virtualnetworkpeerings"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.DataSource = VirtualNetworkPeeringDataSource{}

type VirtualNetworkPeeringDataSource struct{}

type VirtualNetworkPeeringDataSourceModel struct {
	Name                      string `tfschema:"name"`
	VirtualNetworkId          string `tfschema:"virtual_network_id"`
	RemoteVirtualNetworkId    string `tfschema:"remote_virtual_network_id"`
	AllowVirtualNetworkAccess bool   `tfschema:"allow_virtual_network_access"`
	AllowForwardedTraffic     bool   `tfschema:"allow_forwarded_traffic"`
	AllowGatewayTransit       bool   `tfschema:"allow_gateway_transit"`
	OnlyIPv6PeeringEnabled    bool   `tfschema:"only_ipv6_peering_enabled"`
	PeerCompleteVnetsEnabled  bool   `tfschema:"peer_complete_virtual_networks_enabled"`
	UseRemoteGateways         bool   `tfschema:"use_remote_gateways"`
}

func (VirtualNetworkPeeringDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"virtual_network_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: commonids.ValidateVirtualNetworkID,
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
	return &VirtualNetworkPeeringDataSourceModel{}
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

			var state VirtualNetworkPeeringDataSourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			virtualNetworkId, err := commonids.ParseVirtualNetworkID(state.VirtualNetworkId)
			if err != nil {
				return err
			}

			id := virtualnetworkpeerings.NewVirtualNetworkPeeringID(subscriptionId, virtualNetworkId.ResourceGroupName, virtualNetworkId.VirtualNetworkName, state.Name)

			resp, err := client.Get(ctx, id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			metadata.SetID(id)

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					state.AllowVirtualNetworkAccess = pointer.From(props.AllowVirtualNetworkAccess)
					state.AllowForwardedTraffic = pointer.From(props.AllowForwardedTraffic)
					state.AllowGatewayTransit = pointer.From(props.AllowGatewayTransit)
					state.OnlyIPv6PeeringEnabled = pointer.From(props.EnableOnlyIPv6Peering)
					state.PeerCompleteVnetsEnabled = pointer.From(props.PeerCompleteVnets)
					state.UseRemoteGateways = pointer.From(props.UseRemoteGateways)

					remoteVirtualNetworkId := ""
					if network := props.RemoteVirtualNetwork; network != nil {
						parsed, err := commonids.ParseVirtualNetworkIDInsensitively(*network.Id)
						if err != nil {
							return err
						}
						remoteVirtualNetworkId = parsed.ID()
					}
					state.RemoteVirtualNetworkId = remoteVirtualNetworkId
				}
			}
			return metadata.Encode(&state)
		},
	}
}
