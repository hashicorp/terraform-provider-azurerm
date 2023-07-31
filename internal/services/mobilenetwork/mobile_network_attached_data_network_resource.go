package mobilenetwork

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mobilenetwork/2022-11-01/attacheddatanetwork"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mobilenetwork/2022-11-01/packetcoredataplane"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type AttachedDataNetworkModel struct {
	MobileNetworkDataNetworkName         string                   `tfschema:"mobile_network_data_network_name"`
	MobileNetworkPacketCoreDataPlaneId   string                   `tfschema:"mobile_network_packet_core_data_plane_id"`
	DnsAddresses                         []string                 `tfschema:"dns_addresses"`
	Location                             string                   `tfschema:"location"`
	NaptConfiguration                    []NaptConfigurationModel `tfschema:"network_address_port_translation"`
	Tags                                 map[string]interface{}   `tfschema:"tags"`
	UserEquipmentAddressPoolPrefix       []string                 `tfschema:"user_equipment_address_pool_prefixes"`
	UserEquipmentStaticAddressPoolPrefix []string                 `tfschema:"user_equipment_static_address_pool_prefixes"`
	UserPlaneAccessIPv4Address           string                   `tfschema:"user_plane_access_ipv4_address"`
	UserPlaneAccessIPv4Gateway           string                   `tfschema:"user_plane_access_ipv4_gateway"`
	UserPlaneAccessIPv4Subnet            string                   `tfschema:"user_plane_access_ipv4_subnet"`
	UserPlaneAccessName                  string                   `tfschema:"user_plane_access_name"`
}

type NaptConfigurationModel struct {
	PinholeLimits           int64            `tfschema:"pinhole_maximum_number"`
	IcmpPinholeTimeout      int64            `tfschema:"icmp_pinhole_timeout_in_seconds"`
	TcpPinholeTimeout       int64            `tfschema:"tcp_pinhole_timeout_in_seconds"`
	UdpPinholeTimeout       int64            `tfschema:"udp_pinhole_timeout_in_seconds"`
	PortRange               []PortRangeModel `tfschema:"port_range"`
	TcpReuseMinimumHoldTime int64            `tfschema:"tcp_port_reuse_minimum_hold_time_in_seconds"`
	UdpReuseMinimumHoldTime int64            `tfschema:"udp_port_reuse_minimum_hold_time_in_seconds"`
}

type PortRangeModel struct {
	Maximum int64 `tfschema:"maximum"`
	Minimum int64 `tfschema:"minimum"`
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
		"mobile_network_data_network_name": {
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

		"network_address_port_translation": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"pinhole_maximum_number": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      65536,
						ValidateFunc: validation.IntBetween(1, 65536),
					},

					"icmp_pinhole_timeout_in_seconds": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      180,
						ValidateFunc: validation.IntBetween(1, 180),
					},

					"tcp_pinhole_timeout_in_seconds": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      180,
						ValidateFunc: validation.IntBetween(1, 180),
					},

					"udp_pinhole_timeout_in_seconds": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      180,
						ValidateFunc: validation.IntBetween(1, 180),
					},

					"port_range": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"maximum": {
									Type:         pluginsdk.TypeInt,
									Optional:     true,
									Default:      49999,
									ValidateFunc: validation.IntBetween(1024, 65535),
								},

								"minimum": {
									Type:         pluginsdk.TypeInt,
									Optional:     true,
									Default:      1024,
									ValidateFunc: validation.IntBetween(1024, 65535),
								},
							},
						},
					},

					"tcp_port_reuse_minimum_hold_time_in_seconds": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      120,
						ValidateFunc: validation.IntAtLeast(1),
					},

					"udp_port_reuse_minimum_hold_time_in_seconds": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      60,
						ValidateFunc: validation.IntAtLeast(60),
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

			id := attacheddatanetwork.NewAttachedDataNetworkID(packetCoreDataPlaneId.SubscriptionId, packetCoreDataPlaneId.ResourceGroupName, packetCoreDataPlaneId.PacketCoreControlPlaneName, packetCoreDataPlaneId.PacketCoreDataPlaneName, model.MobileNetworkDataNetworkName)
			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			attachedDataNetwork := attacheddatanetwork.AttachedDataNetwork{
				Location: location.Normalize(model.Location),
				Properties: attacheddatanetwork.AttachedDataNetworkPropertiesFormat{
					DnsAddresses:           model.DnsAddresses,
					UserPlaneDataInterface: attacheddatanetwork.InterfaceProperties{},
					NaptConfiguration:      expandNaptConfiguration(model.NaptConfiguration),
				},
				Tags: tags.Expand(model.Tags),
			}

			// if we pass an empty array the service will return an error
			// Array is too short (0), minimum 1.
			if len(model.UserEquipmentStaticAddressPoolPrefix) > 0 {
				attachedDataNetwork.Properties.UserEquipmentStaticAddressPoolPrefix = &model.UserEquipmentStaticAddressPoolPrefix
			}

			if len(model.UserEquipmentAddressPoolPrefix) > 0 {
				attachedDataNetwork.Properties.UserEquipmentAddressPoolPrefix = &model.UserEquipmentAddressPoolPrefix
			}

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

			if err := client.CreateOrUpdateThenPoll(ctx, id, attachedDataNetwork); err != nil {
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

			if resp.Model == nil {
				return fmt.Errorf("retrieving %s: Model was nil", id)
			}

			attachedDataNetwork := *resp.Model

			if metadata.ResourceData.HasChange("dns_addresses") {
				attachedDataNetwork.Properties.DnsAddresses = plan.DnsAddresses
			}

			if metadata.ResourceData.HasChange("network_address_port_translation") {
				attachedDataNetwork.Properties.NaptConfiguration = expandNaptConfiguration(plan.NaptConfiguration)
			}

			if metadata.ResourceData.HasChange("user_equipment_address_pool_prefixes") {
				// if we pass an empty array the service will return an error
				// Array is too short (0), minimum 1.
				if len(plan.UserEquipmentAddressPoolPrefix) > 0 {
					attachedDataNetwork.Properties.UserEquipmentAddressPoolPrefix = &plan.UserEquipmentAddressPoolPrefix
				} else {
					attachedDataNetwork.Properties.UserEquipmentAddressPoolPrefix = nil
				}
			}

			if metadata.ResourceData.HasChange("user_equipment_static_address_pool_prefixes") {
				if len(plan.UserEquipmentStaticAddressPoolPrefix) > 0 {
					attachedDataNetwork.Properties.UserEquipmentStaticAddressPoolPrefix = &plan.UserEquipmentStaticAddressPoolPrefix
				} else {
					attachedDataNetwork.Properties.UserEquipmentStaticAddressPoolPrefix = nil
				}
			}

			if metadata.ResourceData.HasChange("user_plane_access_name") {
				if plan.UserPlaneAccessName != "" {
					attachedDataNetwork.Properties.UserPlaneDataInterface.Name = &plan.UserPlaneAccessName
				} else {
					attachedDataNetwork.Properties.UserPlaneDataInterface.Name = nil
				}
			}

			if metadata.ResourceData.HasChange("user_plane_access_ipv4_address") {
				if plan.UserPlaneAccessIPv4Address != "" {
					attachedDataNetwork.Properties.UserPlaneDataInterface.IPv4Address = &plan.UserPlaneAccessIPv4Address
				} else {
					attachedDataNetwork.Properties.UserPlaneDataInterface.IPv4Address = nil
				}
			}

			if metadata.ResourceData.HasChange("user_plane_access_ipv4_subnet") {
				if plan.UserPlaneAccessIPv4Subnet != "" {
					attachedDataNetwork.Properties.UserPlaneDataInterface.IPv4Subnet = &plan.UserPlaneAccessIPv4Subnet
				} else {
					attachedDataNetwork.Properties.UserPlaneDataInterface.IPv4Subnet = nil
				}
			}

			if metadata.ResourceData.HasChange("user_plane_access_ipv4_gateway") {
				if plan.UserPlaneAccessIPv4Gateway != "" {
					attachedDataNetwork.Properties.UserPlaneDataInterface.IPv4Gateway = &plan.UserPlaneAccessIPv4Gateway
				} else {
					attachedDataNetwork.Properties.UserPlaneDataInterface.IPv4Gateway = nil
				}
			}

			if metadata.ResourceData.HasChange("tags") {
				// pass empty array instead of nil to remove all tags
				attachedDataNetwork.Tags = tags.Expand(plan.Tags)
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, attachedDataNetwork); err != nil {
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
				MobileNetworkDataNetworkName:       id.AttachedDataNetworkName,
				MobileNetworkPacketCoreDataPlaneId: packetcoredataplane.NewPacketCoreDataPlaneID(id.SubscriptionId, id.ResourceGroupName, id.PacketCoreControlPlaneName, id.PacketCoreDataPlaneName).ID(),
			}

			if model := resp.Model; model != nil {
				props := model.Properties

				state.Location = location.Normalize(model.Location)
				state.DnsAddresses = props.DnsAddresses
				state.NaptConfiguration = flattenNaptConfiguration(props.NaptConfiguration)
				state.UserEquipmentAddressPoolPrefix = pointer.From(props.UserEquipmentAddressPoolPrefix)
				state.UserEquipmentStaticAddressPoolPrefix = pointer.From(props.UserEquipmentStaticAddressPoolPrefix)
				state.UserPlaneAccessIPv4Address = pointer.From(props.UserPlaneDataInterface.IPv4Address)
				state.UserPlaneAccessIPv4Gateway = pointer.From(props.UserPlaneDataInterface.IPv4Gateway)
				state.UserPlaneAccessIPv4Subnet = pointer.From(props.UserPlaneDataInterface.IPv4Subnet)
				state.UserPlaneAccessName = pointer.From(props.UserPlaneDataInterface.Name)
				state.Tags = tags.Flatten(model.Tags)
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

			return nil
		},
	}
}

func expandNaptConfiguration(inputList []NaptConfigurationModel) *attacheddatanetwork.NaptConfiguration {
	if len(inputList) == 0 {
		return nil
	}

	input := inputList[0]
	output := attacheddatanetwork.NaptConfiguration{
		PinholeLimits: &input.PinholeLimits,
	}

	output.Enabled = pointer.To(attacheddatanetwork.NaptEnabledEnabled)

	output.PinholeTimeouts = &attacheddatanetwork.PinholeTimeouts{
		Icmp: &input.IcmpPinholeTimeout,
		Tcp:  &input.TcpPinholeTimeout,
		Udp:  &input.UdpPinholeTimeout,
	}

	output.PortRange = expandPortRange(input.PortRange)

	output.PortReuseHoldTime = &attacheddatanetwork.PortReuseHoldTimes{
		Tcp: &input.TcpReuseMinimumHoldTime,
		Udp: &input.UdpReuseMinimumHoldTime,
	}

	return &output
}

func expandPortRange(inputList []PortRangeModel) *attacheddatanetwork.PortRange {
	if len(inputList) == 0 {
		return nil
	}

	input := &inputList[0]
	output := attacheddatanetwork.PortRange{
		MaxPort: &input.Maximum,
		MinPort: &input.Minimum,
	}

	return &output
}

func flattenNaptConfiguration(input *attacheddatanetwork.NaptConfiguration) []NaptConfigurationModel {
	if input == nil {
		return []NaptConfigurationModel{}
	}
	output := NaptConfigurationModel{}

	output.PinholeLimits = pointer.From(input.PinholeLimits)

	if input.PinholeTimeouts != nil {
		output.IcmpPinholeTimeout = pointer.From(input.PinholeTimeouts.Icmp)
		output.TcpPinholeTimeout = pointer.From(input.PinholeTimeouts.Tcp)
		output.UdpPinholeTimeout = pointer.From(input.PinholeTimeouts.Udp)
	}

	output.PortRange = flattenPortRange(input.PortRange)

	if input.PortReuseHoldTime != nil {
		output.TcpReuseMinimumHoldTime = pointer.From(input.PortReuseHoldTime.Tcp)
		output.UdpReuseMinimumHoldTime = pointer.From(input.PortReuseHoldTime.Udp)
	}

	return []NaptConfigurationModel{output}
}

func flattenPortRange(input *attacheddatanetwork.PortRange) []PortRangeModel {
	if input == nil {
		return []PortRangeModel{}
	}

	output := PortRangeModel{}

	output.Maximum = pointer.From(input.MaxPort)

	output.Minimum = pointer.From(input.MinPort)

	return []PortRangeModel{output}
}
