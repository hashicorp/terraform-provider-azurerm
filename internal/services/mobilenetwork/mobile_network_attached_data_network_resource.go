package mobilenetwork

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mobilenetwork/2022-04-01-preview/attacheddatanetwork"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mobilenetwork/2022-04-01-preview/packetcoredataplane"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type AttachedDataNetworkModel struct {
	Name                                 string                     `tfschema:"name"`
	MobileNetworkPacketCoreDataPlaneId   string                     `tfschema:"mobile_network_packet_core_data_plane_id"`
	DnsAddresses                         []string                   `tfschema:"dns_addresses"`
	Location                             string                     `tfschema:"location"`
	NaptConfiguration                    []NaptConfigurationModel   `tfschema:"network_address_port_translation_configuration"`
	Tags                                 map[string]string          `tfschema:"tags"`
	UserEquipmentAddressPoolPrefix       []string                   `tfschema:"user_equipment_address_pool_prefixes"`
	UserEquipmentStaticAddressPoolPrefix []string                   `tfschema:"user_equipment_static_address_pool_prefixes"`
	UserPlaneDataInterface               []InterfacePropertiesModel `tfschema:"user_plane_data_interface"`
}

type NaptConfigurationModel struct {
	Enabled           bool                      `tfschema:"enabled"`
	PinholeLimits     int64                     `tfschema:"pinhole_maximum_number"`
	PinholeTimeouts   []PinholeTimeoutsModel    `tfschema:"pinhole_timeouts_in_seconds"`
	PortRange         []PortRangeModel          `tfschema:"port_range"`
	PortReuseHoldTime []PortReuseHoldTimesModel `tfschema:"port_reuse_minimum_hold_time_in_seconds"`
}

type PinholeTimeoutsModel struct {
	Icmp int64 `tfschema:"icmp"`
	Tcp  int64 `tfschema:"tcp"`
	Udp  int64 `tfschema:"udp"`
}

type PortRangeModel struct {
	MaxPort int64 `tfschema:"max_port"`
	MinPort int64 `tfschema:"min_port"`
}

type PortReuseHoldTimesModel struct {
	Tcp int64 `tfschema:"tcp"`
	Udp int64 `tfschema:"udp"`
}

type AttachedDataNetworkResource struct{}

var _ sdk.ResourceWithUpdate = AttachedDataNetworkResource{}

func (r AttachedDataNetworkResource) ResourceType() string {
	return "azurerm_mobile_network_attached_data_network"
}

func (r AttachedDataNetworkResource) ModelObject() interface{} {
	return &AttachedDataNetworkModel{}
}

func (r AttachedDataNetworkResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return attacheddatanetwork.ValidateAttachedDataNetworkID
}

func (r AttachedDataNetworkResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"mobile_network_packet_core_data_plane_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: packetcoredataplane.ValidatePacketCoreDataPlaneID,
		},

		"location": commonschema.Location(),

		"dns_addresses": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validate.IPv4Address,
			},
		},

		"network_address_port_translation_configuration": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"enabled": {
						Type:     pluginsdk.TypeBool,
						Required: true,
					},

					"pinhole_maximum_number": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      65536,
						ValidateFunc: validation.IntBetween(1, 65536),
					},

					"pinhole_timeouts_in_seconds": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"icmp": {
									Type:         pluginsdk.TypeInt,
									Optional:     true,
									Default:      180,
									ValidateFunc: validation.IntBetween(1, 180),
								},

								"tcp": {
									Type:         pluginsdk.TypeInt,
									Optional:     true,
									Default:      180,
									ValidateFunc: validation.IntBetween(1, 180),
								},

								"udp": {
									Type:         pluginsdk.TypeInt,
									Optional:     true,
									Default:      180,
									ValidateFunc: validation.IntBetween(1, 180),
								},
							},
						},
					},

					"port_range": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"max_port": {
									Type:         pluginsdk.TypeInt,
									Optional:     true,
									Default:      49999,
									ValidateFunc: validation.IntBetween(1024, 65535),
								},

								"min_port": {
									Type:         pluginsdk.TypeInt,
									Optional:     true,
									Default:      1024,
									ValidateFunc: validation.IntBetween(1024, 65535),
								},
							},
						},
					},

					"port_reuse_minimum_hold_time_in_seconds": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"tcp": {
									Type:         pluginsdk.TypeInt,
									Optional:     true,
									Default:      120,
									ValidateFunc: validation.IntAtLeast(1),
								},

								"udp": {
									Type:         pluginsdk.TypeInt,
									Optional:     true,
									Default:      60,
									ValidateFunc: validation.IntAtLeast(60),
								},
							},
						},
					},
				},
			},
		},

		"user_plane_data_interface": interfacePropertiesSchema(),

		"user_equipment_address_pool_prefixes": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validate.CIDR,
			},
		},

		"user_equipment_static_address_pool_prefixes": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validate.CIDR,
			},
		},

		"tags": commonschema.Tags(),
	}
}

func (r AttachedDataNetworkResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r AttachedDataNetworkResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model AttachedDataNetworkModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.MobileNetwork.AttachedDataNetworkClient
			packetCoreDataPlaneId, err := packetcoredataplane.ParsePacketCoreDataPlaneID(model.MobileNetworkPacketCoreDataPlaneId)
			if err != nil {
				return err
			}

			id := attacheddatanetwork.NewAttachedDataNetworkID(packetCoreDataPlaneId.SubscriptionId, packetCoreDataPlaneId.ResourceGroupName, packetCoreDataPlaneId.PacketCoreControlPlaneName, packetCoreDataPlaneId.PacketCoreDataPlaneName, model.Name)
			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			properties := &attacheddatanetwork.AttachedDataNetwork{
				Location: location.Normalize(model.Location),
				Properties: attacheddatanetwork.AttachedDataNetworkPropertiesFormat{
					DnsAddresses:                         &model.DnsAddresses,
					UserEquipmentAddressPoolPrefix:       &model.UserEquipmentAddressPoolPrefix,
					UserEquipmentStaticAddressPoolPrefix: &model.UserEquipmentStaticAddressPoolPrefix,
				},
				Tags: &model.Tags,
			}

			naptConfigurationValue, err := expandNaptConfigurationModel(model.NaptConfiguration)
			if err != nil {
				return err
			}

			properties.Properties.NaptConfiguration = naptConfigurationValue

			userPlaneDataInterfaceValue, err := expandAttachedDataNetworkInterfacePropertiesModel(model.UserPlaneDataInterface)
			if err != nil {
				return err
			}

			if userPlaneDataInterfaceValue != nil {
				properties.Properties.UserPlaneDataInterface = *userPlaneDataInterfaceValue
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, *properties); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r AttachedDataNetworkResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MobileNetwork.AttachedDataNetworkClient

			id, err := attacheddatanetwork.ParseAttachedDataNetworkID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model AttachedDataNetworkModel
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

			if metadata.ResourceData.HasChange("dns_addresses") {
				properties.Properties.DnsAddresses = &model.DnsAddresses
			}

			if metadata.ResourceData.HasChange("network_address_port_translation_configuration") {
				naptConfigurationValue, err := expandNaptConfigurationModel(model.NaptConfiguration)
				if err != nil {
					return err
				}

				properties.Properties.NaptConfiguration = naptConfigurationValue
			}

			if metadata.ResourceData.HasChange("user_equipment_address_pool_prefixes") {
				properties.Properties.UserEquipmentAddressPoolPrefix = &model.UserEquipmentAddressPoolPrefix
			}

			if metadata.ResourceData.HasChange("user_equipment_static_address_pool_prefixes") {
				properties.Properties.UserEquipmentStaticAddressPoolPrefix = &model.UserEquipmentStaticAddressPoolPrefix
			}

			if metadata.ResourceData.HasChange("user_plane_data_interface") {
				userPlaneDataInterfaceValue, err := expandAttachedDataNetworkInterfacePropertiesModel(model.UserPlaneDataInterface)
				if err != nil {
					return err
				}

				if userPlaneDataInterfaceValue != nil {
					properties.Properties.UserPlaneDataInterface = *userPlaneDataInterfaceValue
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

func (r AttachedDataNetworkResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MobileNetwork.AttachedDataNetworkClient

			id, err := attacheddatanetwork.ParseAttachedDataNetworkID(metadata.ResourceData.Id())
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

			state := AttachedDataNetworkModel{
				Name:                               id.AttachedDataNetworkName,
				MobileNetworkPacketCoreDataPlaneId: packetcoredataplane.NewPacketCoreDataPlaneID(id.SubscriptionId, id.ResourceGroupName, id.PacketCoreControlPlaneName, id.PacketCoreDataPlaneName).ID(),
				Location:                           location.Normalize(model.Location),
			}

			properties := &model.Properties
			if properties.DnsAddresses != nil {
				state.DnsAddresses = *properties.DnsAddresses
			}

			naptConfigurationValue, err := flattenNaptConfigurationModel(properties.NaptConfiguration)
			if err != nil {
				return err
			}

			state.NaptConfiguration = naptConfigurationValue

			if properties.UserEquipmentAddressPoolPrefix != nil {
				state.UserEquipmentAddressPoolPrefix = *properties.UserEquipmentAddressPoolPrefix
			}

			if properties.UserEquipmentStaticAddressPoolPrefix != nil {
				state.UserEquipmentStaticAddressPoolPrefix = *properties.UserEquipmentStaticAddressPoolPrefix
			}

			userPlaneDataInterfaceValue, err := flattenAttachedDataNetworkInterfacePropertiesModel(&properties.UserPlaneDataInterface)
			if err != nil {
				return err
			}

			state.UserPlaneDataInterface = userPlaneDataInterfaceValue
			if model.Tags != nil {
				state.Tags = *model.Tags
			}

			return metadata.Encode(&state)
		},
	}
}

func (r AttachedDataNetworkResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MobileNetwork.AttachedDataNetworkClient

			id, err := attacheddatanetwork.ParseAttachedDataNetworkID(metadata.ResourceData.Id())
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

func expandNaptConfigurationModel(inputList []NaptConfigurationModel) (*attacheddatanetwork.NaptConfiguration, error) {
	if len(inputList) == 0 {
		return nil, nil
	}

	input := &inputList[0]
	output := attacheddatanetwork.NaptConfiguration{
		PinholeLimits: &input.PinholeLimits,
	}

	naptEnabled := attacheddatanetwork.NaptEnabledDisabled
	if input.Enabled {
		naptEnabled = attacheddatanetwork.NaptEnabledEnabled
	}
	output.Enabled = &naptEnabled

	pinholeTimeoutsValue, err := expandPinholeTimeoutsModel(input.PinholeTimeouts)
	if err != nil {
		return nil, err
	}

	output.PinholeTimeouts = pinholeTimeoutsValue

	portRangeValue, err := expandPortRangeModel(input.PortRange)
	if err != nil {
		return nil, err
	}

	output.PortRange = portRangeValue

	portReuseHoldTimeValue, err := expandPortReuseHoldTimesModel(input.PortReuseHoldTime)
	if err != nil {
		return nil, err
	}

	output.PortReuseHoldTime = portReuseHoldTimeValue

	return &output, nil
}

func expandPinholeTimeoutsModel(inputList []PinholeTimeoutsModel) (*attacheddatanetwork.PinholeTimeouts, error) {
	if len(inputList) == 0 {
		return nil, nil
	}

	input := &inputList[0]
	output := attacheddatanetwork.PinholeTimeouts{
		Icmp: &input.Icmp,
		Tcp:  &input.Tcp,
		Udp:  &input.Udp,
	}

	return &output, nil
}

func expandPortRangeModel(inputList []PortRangeModel) (*attacheddatanetwork.PortRange, error) {
	if len(inputList) == 0 {
		return nil, nil
	}

	input := &inputList[0]
	output := attacheddatanetwork.PortRange{
		MaxPort: &input.MaxPort,
		MinPort: &input.MinPort,
	}

	return &output, nil
}

func expandPortReuseHoldTimesModel(inputList []PortReuseHoldTimesModel) (*attacheddatanetwork.PortReuseHoldTimes, error) {
	if len(inputList) == 0 {
		return nil, nil
	}

	input := &inputList[0]
	output := attacheddatanetwork.PortReuseHoldTimes{
		Tcp: &input.Tcp,
		Udp: &input.Udp,
	}

	return &output, nil
}

func expandAttachedDataNetworkInterfacePropertiesModel(inputList []InterfacePropertiesModel) (*attacheddatanetwork.InterfaceProperties, error) {
	if len(inputList) == 0 {
		return nil, nil
	}
	input := &inputList[0]
	output := attacheddatanetwork.InterfaceProperties{}

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

func flattenNaptConfigurationModel(input *attacheddatanetwork.NaptConfiguration) ([]NaptConfigurationModel, error) {
	var outputList []NaptConfigurationModel
	if input == nil {
		return outputList, nil
	}

	output := NaptConfigurationModel{
		Enabled: input.Enabled != nil && *input.Enabled == attacheddatanetwork.NaptEnabledEnabled,
	}

	if input.PinholeLimits != nil {
		output.PinholeLimits = *input.PinholeLimits
	}

	pinholeTimeoutsValue, err := flattenPinholeTimeoutsModel(input.PinholeTimeouts)
	if err != nil {
		return nil, err
	}

	output.PinholeTimeouts = pinholeTimeoutsValue

	portRangeValue, err := flattenPortRangeModel(input.PortRange)
	if err != nil {
		return nil, err
	}

	output.PortRange = portRangeValue

	portReuseHoldTimeValue, err := flattenPortReuseHoldTimesModel(input.PortReuseHoldTime)
	if err != nil {
		return nil, err
	}

	output.PortReuseHoldTime = portReuseHoldTimeValue

	return append(outputList, output), nil
}

func flattenPinholeTimeoutsModel(input *attacheddatanetwork.PinholeTimeouts) ([]PinholeTimeoutsModel, error) {
	var outputList []PinholeTimeoutsModel
	if input == nil {
		return outputList, nil
	}

	output := PinholeTimeoutsModel{}

	if input.Icmp != nil {
		output.Icmp = *input.Icmp
	}

	if input.Tcp != nil {
		output.Tcp = *input.Tcp
	}

	if input.Udp != nil {
		output.Udp = *input.Udp
	}

	outputList = append(outputList, output)

	return outputList, nil
}

func flattenPortRangeModel(input *attacheddatanetwork.PortRange) ([]PortRangeModel, error) {
	var outputList []PortRangeModel
	if input == nil {
		return outputList, nil
	}

	output := PortRangeModel{}

	if input.MaxPort != nil {
		output.MaxPort = *input.MaxPort
	}

	if input.MinPort != nil {
		output.MinPort = *input.MinPort
	}

	return append(outputList, output), nil
}

func flattenPortReuseHoldTimesModel(input *attacheddatanetwork.PortReuseHoldTimes) ([]PortReuseHoldTimesModel, error) {
	var outputList []PortReuseHoldTimesModel
	if input == nil {
		return outputList, nil
	}

	output := PortReuseHoldTimesModel{}

	if input.Tcp != nil {
		output.Tcp = *input.Tcp
	}

	if input.Udp != nil {
		output.Udp = *input.Udp
	}

	return append(outputList, output), nil
}

func flattenAttachedDataNetworkInterfacePropertiesModel(input *attacheddatanetwork.InterfaceProperties) ([]InterfacePropertiesModel, error) {
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
