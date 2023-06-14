package mobilenetwork

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mobilenetwork/2022-11-01/attacheddatanetwork"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mobilenetwork/2022-11-01/packetcoredataplane"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type AttachedDataNetworkModel struct {
	Name                                 string                   `tfschema:"name"`
	MobileNetworkPacketCoreDataPlaneId   string                   `tfschema:"mobile_network_packet_core_data_plane_id"`
	DnsAddresses                         []string                 `tfschema:"dns_addresses"`
	Location                             string                   `tfschema:"location"`
	NaptConfiguration                    []NaptConfigurationModel `tfschema:"network_address_port_translation_configuration"`
	Tags                                 map[string]string        `tfschema:"tags"`
	UserEquipmentAddressPoolPrefix       []string                 `tfschema:"user_equipment_address_pool_prefixes"`
	UserEquipmentStaticAddressPoolPrefix []string                 `tfschema:"user_equipment_static_address_pool_prefixes"`
	UserPlaneAccessIPv4Address           string                   `tfschema:"user_plane_access_ipv4_address"`
	UserPlaneAccessIPv4Gateway           string                   `tfschema:"user_plane_access_ipv4_gateway"`
	UserPlaneAccessIPv4Subnet            string                   `tfschema:"user_plane_access_ipv4_subnet"`
	UserPlaneAccessName                  string                   `tfschema:"user_plane_access_name"`
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
			Required: true,
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

		"user_plane_access_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"user_plane_access_ipv4_address": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.IsIPv4Address,
		},

		"user_plane_access_ipv4_subnet": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validate.CIDR,
		},

		"user_plane_access_ipv4_gateway": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.IsIPv4Address,
		},

		"user_equipment_address_pool_prefixes": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validate.CIDR,
			},
			AtLeastOneOf: []string{
				"user_equipment_address_pool_prefixes",
				"user_equipment_static_address_pool_prefixes",
			},
		},

		"user_equipment_static_address_pool_prefixes": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validate.CIDR,
			},
			AtLeastOneOf: []string{
				"user_equipment_address_pool_prefixes",
				"user_equipment_static_address_pool_prefixes",
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
		Timeout: 180 * time.Minute,
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

			attachedDataNetwork := &attacheddatanetwork.AttachedDataNetwork{
				Location: location.Normalize(model.Location),
				Properties: attacheddatanetwork.AttachedDataNetworkPropertiesFormat{
					DnsAddresses:                         model.DnsAddresses,
					UserEquipmentAddressPoolPrefix:       &model.UserEquipmentAddressPoolPrefix,
					UserEquipmentStaticAddressPoolPrefix: &model.UserEquipmentStaticAddressPoolPrefix,
					UserPlaneDataInterface:               attacheddatanetwork.InterfaceProperties{},
				},
				Tags: &model.Tags,
			}

			attachedDataNetwork.Properties.NaptConfiguration = expandNaptConfigurationModel(model.NaptConfiguration)

			if model.UserPlaneAccessName != "" {
				attachedDataNetwork.Properties.UserPlaneDataInterface.Name = &model.UserPlaneAccessName
			}

			if model.UserPlaneAccessIPv4Address != "" {
				attachedDataNetwork.Properties.UserPlaneDataInterface.IPv4Address = &model.UserPlaneAccessIPv4Address
			}

			if model.UserPlaneAccessIPv4Subnet != "" {
				attachedDataNetwork.Properties.UserPlaneDataInterface.IPv4Subnet = &model.UserPlaneAccessIPv4Subnet
			}

			if model.UserPlaneAccessIPv4Gateway != "" {
				attachedDataNetwork.Properties.UserPlaneDataInterface.IPv4Gateway = &model.UserPlaneAccessIPv4Gateway
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, *attachedDataNetwork); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r AttachedDataNetworkResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 180 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MobileNetwork.AttachedDataNetworkClient

			id, err := attacheddatanetwork.ParseAttachedDataNetworkID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var plan AttachedDataNetworkModel
			if err := metadata.Decode(&plan); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			addtachedDataNetwork := resp.Model
			if addtachedDataNetwork == nil {
				return fmt.Errorf("retrieving %s: addtachedDataNetwork was nil", id)
			}

			if metadata.ResourceData.HasChange("dns_addresses") {
				addtachedDataNetwork.Properties.DnsAddresses = plan.DnsAddresses
			}

			if metadata.ResourceData.HasChange("network_address_port_translation_configuration") {
				addtachedDataNetwork.Properties.NaptConfiguration = expandNaptConfigurationModel(plan.NaptConfiguration)
			}

			if metadata.ResourceData.HasChange("user_equipment_address_pool_prefixes") {
				addtachedDataNetwork.Properties.UserEquipmentAddressPoolPrefix = &plan.UserEquipmentAddressPoolPrefix
			}

			if metadata.ResourceData.HasChange("user_equipment_static_address_pool_prefixes") {
				addtachedDataNetwork.Properties.UserEquipmentStaticAddressPoolPrefix = &plan.UserEquipmentStaticAddressPoolPrefix
			}

			if metadata.ResourceData.HasChange("user_plane_access_name") {
				addtachedDataNetwork.Properties.UserPlaneDataInterface.Name = &plan.UserPlaneAccessName
			}

			if metadata.ResourceData.HasChange("user_plane_access_ipv4_address") {
				addtachedDataNetwork.Properties.UserPlaneDataInterface.IPv4Address = &plan.UserPlaneAccessIPv4Address
			}

			if metadata.ResourceData.HasChange("user_plane_access_ipv4_subnet") {
				addtachedDataNetwork.Properties.UserPlaneDataInterface.IPv4Subnet = &plan.UserPlaneAccessIPv4Subnet
			}

			if metadata.ResourceData.HasChange("user_plane_access_ipv4_gateway") {
				addtachedDataNetwork.Properties.UserPlaneDataInterface.IPv4Gateway = &plan.UserPlaneAccessIPv4Gateway
			}

			if metadata.ResourceData.HasChange("tags") {
				addtachedDataNetwork.Tags = &plan.Tags
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, *addtachedDataNetwork); err != nil {
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

			state := AttachedDataNetworkModel{
				Name:                               id.AttachedDataNetworkName,
				MobileNetworkPacketCoreDataPlaneId: packetcoredataplane.NewPacketCoreDataPlaneID(id.SubscriptionId, id.ResourceGroupName, id.PacketCoreControlPlaneName, id.PacketCoreDataPlaneName).ID(),
			}

			if model := resp.Model; model != nil {
				props := model.Properties

				state.Location = location.Normalize(model.Location)
				state.DnsAddresses = props.DnsAddresses
				state.NaptConfiguration = flattenNaptConfigurationModel(props.NaptConfiguration)
				state.UserEquipmentAddressPoolPrefix = pointer.From(props.UserEquipmentAddressPoolPrefix)
				state.UserEquipmentStaticAddressPoolPrefix = pointer.From(props.UserEquipmentStaticAddressPoolPrefix)
				state.UserPlaneAccessIPv4Address = pointer.From(props.UserPlaneDataInterface.IPv4Address)
				state.UserPlaneAccessIPv4Gateway = pointer.From(props.UserPlaneDataInterface.IPv4Gateway)
				state.UserPlaneAccessIPv4Subnet = pointer.From(props.UserPlaneDataInterface.IPv4Subnet)
				state.UserPlaneAccessName = pointer.From(props.UserPlaneDataInterface.Name)
				state.Tags = pointer.From(model.Tags)
			}

			return metadata.Encode(&state)
		},
	}
}

func (r AttachedDataNetworkResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 180 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MobileNetwork.AttachedDataNetworkClient

			id, err := attacheddatanetwork.ParseAttachedDataNetworkID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			if err := resourceMobileNetworkChildWaitForDeletion(ctx, id.ID(), func() (*http.Response, error) {
				resp, err := client.Get(ctx, *id)
				return resp.HttpResponse, err
			}); err != nil {
				return err
			}

			return nil
		},
	}
}

func expandNaptConfigurationModel(inputList []NaptConfigurationModel) *attacheddatanetwork.NaptConfiguration {
	if len(inputList) == 0 {
		return nil
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

	output.PinholeTimeouts = expandPinholeTimeoutsModel(input.PinholeTimeouts)

	output.PortRange = expandPortRangeModel(input.PortRange)

	output.PortReuseHoldTime = expandPortReuseHoldTimesModel(input.PortReuseHoldTime)

	return &output
}

func expandPinholeTimeoutsModel(inputList []PinholeTimeoutsModel) *attacheddatanetwork.PinholeTimeouts {
	if len(inputList) == 0 {
		return nil
	}

	input := &inputList[0]
	output := attacheddatanetwork.PinholeTimeouts{
		Icmp: &input.Icmp,
		Tcp:  &input.Tcp,
		Udp:  &input.Udp,
	}

	return &output
}

func expandPortRangeModel(inputList []PortRangeModel) *attacheddatanetwork.PortRange {
	if len(inputList) == 0 {
		return nil
	}

	input := &inputList[0]
	output := attacheddatanetwork.PortRange{
		MaxPort: &input.MaxPort,
		MinPort: &input.MinPort,
	}

	return &output
}

func expandPortReuseHoldTimesModel(inputList []PortReuseHoldTimesModel) *attacheddatanetwork.PortReuseHoldTimes {
	if len(inputList) == 0 {
		return nil
	}

	input := &inputList[0]
	output := attacheddatanetwork.PortReuseHoldTimes{
		Tcp: &input.Tcp,
		Udp: &input.Udp,
	}

	return &output
}

func flattenNaptConfigurationModel(input *attacheddatanetwork.NaptConfiguration) []NaptConfigurationModel {
	var outputList []NaptConfigurationModel
	if input == nil {
		return outputList
	}

	output := NaptConfigurationModel{
		Enabled: input.Enabled != nil && *input.Enabled == attacheddatanetwork.NaptEnabledEnabled,
	}

	if input.PinholeLimits != nil {
		output.PinholeLimits = *input.PinholeLimits
	}

	output.PinholeTimeouts = flattenPinholeTimeoutsModel(input.PinholeTimeouts)

	output.PortRange = flattenPortRangeModel(input.PortRange)

	output.PortReuseHoldTime = flattenPortReuseHoldTimesModel(input.PortReuseHoldTime)

	return append(outputList, output)
}

func flattenPinholeTimeoutsModel(input *attacheddatanetwork.PinholeTimeouts) []PinholeTimeoutsModel {
	var outputList []PinholeTimeoutsModel
	if input == nil {
		return outputList
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

	return outputList
}

func flattenPortRangeModel(input *attacheddatanetwork.PortRange) []PortRangeModel {
	var outputList []PortRangeModel
	if input == nil {
		return outputList
	}

	output := PortRangeModel{}

	if input.MaxPort != nil {
		output.MaxPort = *input.MaxPort
	}

	if input.MinPort != nil {
		output.MinPort = *input.MinPort
	}

	return append(outputList, output)
}

func flattenPortReuseHoldTimesModel(input *attacheddatanetwork.PortReuseHoldTimes) []PortReuseHoldTimesModel {
	var outputList []PortReuseHoldTimesModel
	if input == nil {
		return outputList
	}

	output := PortReuseHoldTimesModel{}

	if input.Tcp != nil {
		output.Tcp = *input.Tcp
	}

	if input.Udp != nil {
		output.Udp = *input.Udp
	}

	return append(outputList, output)
}
