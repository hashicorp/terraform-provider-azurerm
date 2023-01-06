package mobilenetwork

import (
	"context"
	"fmt"
	"time"

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

		"user_plane_access_interface": interfacePropertiesSchemaDataSource(),
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
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			model := resp.Model
			if model == nil {
				return fmt.Errorf("retrieving %s: model was nil", id)
			}

			state := PacketCoreDataPlaneModel{
				Name:                                  id.PacketCoreDataPlaneName,
				MobileNetworkPacketCoreControlPlaneId: packetcorecontrolplane.NewPacketCoreControlPlaneID(id.SubscriptionId, id.ResourceGroupName, id.PacketCoreControlPlaneName).ID(),
				Location:                              location.Normalize(model.Location),
			}

			properties := &model.Properties
			state.UserPlaneAccessInterface = flattenPacketCoreDataPlaneInterfacePropertiesModel(properties.UserPlaneAccessInterface)

			if model.Tags != nil {
				state.Tags = *model.Tags
			}
			metadata.SetID(id)

			return metadata.Encode(&state)
		},
	}
}
