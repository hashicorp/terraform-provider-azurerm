package mobilenetwork

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mobilenetwork/2022-04-01-preview/packetcorecontrolplane"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mobilenetwork/2022-04-01-preview/packetcoredataplane"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type PacketCoreDataPlaneModel struct {
	Name                                  string                     `tfschema:"name"`
	MobileNetworkPacketCoreControlPlaneId string                     `tfschema:"mobile_network_packet_core_control_plane_id"`
	Location                              string                     `tfschema:"location"`
	Tags                                  map[string]string          `tfschema:"tags"`
	UserPlaneAccessInterface              []InterfacePropertiesModel `tfschema:"user_plane_access_interface"`
}

type PacketCoreDataPlaneResource struct{}

var _ sdk.ResourceWithUpdate = PacketCoreDataPlaneResource{}

func (r PacketCoreDataPlaneResource) ResourceType() string {
	return "azurerm_mobile_network_packet_core_data_plane"
}

func (r PacketCoreDataPlaneResource) ModelObject() interface{} {
	return &PacketCoreDataPlaneModel{}
}

func (r PacketCoreDataPlaneResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return packetcoredataplane.ValidatePacketCoreDataPlaneID
}

func (r PacketCoreDataPlaneResource) Arguments() map[string]*pluginsdk.Schema {
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

		"location": commonschema.Location(),

		"tags": commonschema.Tags(),

		"user_plane_access_interface": interfacePropertiesSchema(),
	}
}

func (r PacketCoreDataPlaneResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r PacketCoreDataPlaneResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model PacketCoreDataPlaneModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.MobileNetwork.PacketCoreDataPlaneClient
			packetCoreControlPlaneId, err := packetcorecontrolplane.ParsePacketCoreControlPlaneID(model.MobileNetworkPacketCoreControlPlaneId)
			if err != nil {
				return err
			}

			id := packetcoredataplane.NewPacketCoreDataPlaneID(packetCoreControlPlaneId.SubscriptionId, packetCoreControlPlaneId.ResourceGroupName, packetCoreControlPlaneId.PacketCoreControlPlaneName, model.Name)
			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			properties := &packetcoredataplane.PacketCoreDataPlane{
				Location:   location.Normalize(model.Location),
				Properties: packetcoredataplane.PacketCoreDataPlanePropertiesFormat{},
				Tags:       &model.Tags,
			}

			userPlaneAccessInterfaceValue, err := expandPacketCoreDataPlaneInterfacePropertiesModel(model.UserPlaneAccessInterface)
			if err != nil {
				return err
			}

			if userPlaneAccessInterfaceValue != nil {
				properties.Properties.UserPlaneAccessInterface = *userPlaneAccessInterfaceValue
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, *properties); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r PacketCoreDataPlaneResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MobileNetwork.PacketCoreDataPlaneClient

			id, err := packetcoredataplane.ParsePacketCoreDataPlaneID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model PacketCoreDataPlaneModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			properties := resp.Model
			if properties == nil {
				return fmt.Errorf("retrieving %s: properties was nil", id)
			}

			if metadata.ResourceData.HasChange("user_plane_access_interface") {
				userPlaneAccessInterfaceValue, err := expandPacketCoreDataPlaneInterfacePropertiesModel(model.UserPlaneAccessInterface)
				//userPlaneAccessInterfaceValue, err := expandPacketCoreDataPlaneInterfacePropertiesModel(model.UserPlaneAccessInterface)
				if err != nil {
					return err
				}

				if userPlaneAccessInterfaceValue != nil {
					properties.Properties.UserPlaneAccessInterface = *userPlaneAccessInterfaceValue
				}
			}

			properties.SystemData = nil

			if metadata.ResourceData.HasChange("tags") {
				properties.Tags = &model.Tags
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, *properties); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r PacketCoreDataPlaneResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MobileNetwork.PacketCoreDataPlaneClient

			id, err := packetcoredataplane.ParsePacketCoreDataPlaneID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
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
			userPlaneAccessInterfaceValue, err := flattenPacketCoreDataPlaneInterfacePropertiesModel(&properties.UserPlaneAccessInterface)
			if err != nil {
				return err
			}

			state.UserPlaneAccessInterface = userPlaneAccessInterfaceValue
			if model.Tags != nil {
				state.Tags = *model.Tags
			}

			return metadata.Encode(&state)
		},
	}
}

func (r PacketCoreDataPlaneResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MobileNetwork.PacketCoreDataPlaneClient

			id, err := packetcoredataplane.ParsePacketCoreDataPlaneID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}

func expandPacketCoreDataPlaneInterfacePropertiesModel(inputList []InterfacePropertiesModel) (*packetcoredataplane.InterfaceProperties, error) {
	if len(inputList) == 0 {
		return nil, nil
	}

	input := &inputList[0]
	output := packetcoredataplane.InterfaceProperties{}

	if input.IPv4Address != "" {
		output.IPv4Address = &input.IPv4Address
	}

	if input.IPv4Gateway != "" {
		output.IPv4Gateway = &input.IPv4Gateway
	}

	if input.IPv4Subnet != "" {
		output.IPv4Subnet = &input.IPv4Subnet
	}

	if input.Name != "" {
		output.Name = &input.Name
	}

	return &output, nil
}

func flattenPacketCoreDataPlaneInterfacePropertiesModel(input *packetcoredataplane.InterfaceProperties) ([]InterfacePropertiesModel, error) {
	var outputList []InterfacePropertiesModel
	if input == nil {
		return outputList, nil
	}

	output := InterfacePropertiesModel{}

	if input.IPv4Address != nil {
		output.IPv4Address = *input.IPv4Address
	}

	if input.IPv4Gateway != nil {
		output.IPv4Gateway = *input.IPv4Gateway
	}

	if input.IPv4Subnet != nil {
		output.IPv4Subnet = *input.IPv4Subnet
	}

	if input.Name != nil {
		output.Name = *input.Name
	}

	return append(outputList, output), nil
}
