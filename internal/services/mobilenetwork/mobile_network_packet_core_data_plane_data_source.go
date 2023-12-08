// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mobilenetwork

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mobilenetwork/2022-11-01/packetcorecontrolplane"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mobilenetwork/2022-11-01/packetcoredataplane"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type PacketCoreDataPlaneDataSource struct{}

var _ sdk.DataSource = PacketCoreDataPlaneDataSource{}

func (r PacketCoreDataPlaneDataSource) ResourceType() string {
	return "azurerm_mobile_network_packet_core_data_plane"
}

func (r PacketCoreDataPlaneDataSource) ModelObject() interface{} {
	return &PacketCoreDataPlaneModel{}
}

func (r PacketCoreDataPlaneDataSource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return packetcoredataplane.ValidatePacketCoreDataPlaneID
}

func (r PacketCoreDataPlaneDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"mobile_network_packet_core_control_plane_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: packetcorecontrolplane.ValidatePacketCoreControlPlaneID,
		},
	}
}

func (r PacketCoreDataPlaneDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.LocationComputed(),

		"tags": commonschema.TagsDataSource(),

		"user_plane_access_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"user_plane_access_ipv4_address": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"user_plane_access_ipv4_subnet": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"user_plane_access_ipv4_gateway": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r PacketCoreDataPlaneDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var metaModel PacketCoreDataPlaneModel
			if err := metadata.Decode(&metaModel); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.MobileNetwork.PacketCoreDataPlaneClient
			packetCoreControlPlaneId, err := packetcorecontrolplane.ParsePacketCoreControlPlaneID(metaModel.MobileNetworkPacketCoreControlPlaneId)
			if err != nil {
				return err
			}

			id := packetcoredataplane.NewPacketCoreDataPlaneID(packetCoreControlPlaneId.SubscriptionId, packetCoreControlPlaneId.ResourceGroupName, packetCoreControlPlaneId.PacketCoreControlPlaneName, metaModel.Name)

			resp, err := client.Get(ctx, id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			state := PacketCoreDataPlaneModel{
				Name:                                  id.PacketCoreDataPlaneName,
				MobileNetworkPacketCoreControlPlaneId: packetcorecontrolplane.NewPacketCoreControlPlaneID(id.SubscriptionId, id.ResourceGroupName, id.PacketCoreControlPlaneName).ID(),
			}

			if resp.Model != nil {
				state.Location = location.Normalize(resp.Model.Location)

				props := resp.Model.Properties

				state.UserPlaneAccessIPv4Address = pointer.From(props.UserPlaneAccessInterface.IPv4Address)
				state.UserPlaneAccessIPv4Gateway = pointer.From(props.UserPlaneAccessInterface.IPv4Gateway)
				state.UserPlaneAccessIPv4Subnet = pointer.From(props.UserPlaneAccessInterface.IPv4Subnet)
				state.UserPlaneAccessName = pointer.From(props.UserPlaneAccessInterface.Name)

				if resp.Model.Tags != nil {
					state.Tags = *resp.Model.Tags
				}
			}

			metadata.SetID(id)

			return metadata.Encode(&state)
		},
	}
}
