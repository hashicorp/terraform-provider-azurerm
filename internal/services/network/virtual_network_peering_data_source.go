// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/virtualnetworkpeerings"
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
			Description:  "The name of this Virtual Network peering.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"virtual_network_id": {
			Description:  "The resource ID of the Virtual Network.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: commonids.ValidateVirtualNetworkID,
		},
	}
}

func (VirtualNetworkPeeringDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"remote_virtual_network_id": {
			Description: "The full Azure resource ID of the remote Virtual Network.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"allow_virtual_network_access": {
			Description: "Controls if the traffic from the local Virtual Network can reach the remote Virtual Network.",
			Type:        pluginsdk.TypeBool,
			Computed:    true,
		},

		"allow_forwarded_traffic": {
			Description: "Controls if forwarded traffic from VMs in the remote Virtual Network is allowed.",
			Type:        pluginsdk.TypeBool,
			Computed:    true,
		},

		"allow_gateway_transit": {
			Description: "Controls gatewayLinks can be used in the remote Virtual Network's link to the local Virtual Network.",
			Type:        pluginsdk.TypeBool,
			Computed:    true,
		},

		"only_ipv6_peering_enabled": {
			Description: "Whether only IPv6 address space is peered for Subnet peering.",
			Type:        pluginsdk.TypeBool,
			Computed:    true,
		},

		"peer_complete_virtual_networks_enabled": {
			Description: "Whether complete Virtual Network address space is peered.",
			Type:        pluginsdk.TypeBool,
			Computed:    true,
		},

		"use_remote_gateways": {
			Description: "Controls if remote gateways can be used on the local Virtual Network.",
			Type:        pluginsdk.TypeBool,
			Computed:    true,
		},
	}
}

func (VirtualNetworkPeeringDataSource) ModelObject() any {
	return &VirtualNetworkPeeringDataSourceModel{}
}

func (VirtualNetworkPeeringDataSource) ResourceType() string {
	return "azurerm_virtual_network_peering"
}

func (VirtualNetworkPeeringDataSource) Description() string {
	return "Use this data source to access information about an existing Virtual Network peering."
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
