// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mobilenetwork

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mobilenetwork/2022-11-01/attacheddatanetwork"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mobilenetwork/2022-11-01/packetcoredataplane"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type AttachedDataNetworkDataSource struct{}

type AttachedDataNetworkDataSourceModel struct {
	MobileNetworkDataNetworkName         string                             `tfschema:"mobile_network_data_network_name"`
	MobileNetworkPacketCoreDataPlaneId   string                             `tfschema:"mobile_network_packet_core_data_plane_id"`
	DnsAddresses                         []string                           `tfschema:"dns_addresses"`
	Location                             string                             `tfschema:"location"`
	NaptConfiguration                    []NaptConfigurationDataSourceModel `tfschema:"network_address_port_translation"`
	Tags                                 map[string]interface{}             `tfschema:"tags"`
	UserEquipmentAddressPoolPrefix       []string                           `tfschema:"user_equipment_address_pool_prefixes"`
	UserEquipmentStaticAddressPoolPrefix []string                           `tfschema:"user_equipment_static_address_pool_prefixes"`
	UserPlaneAccessIPv4Address           string                             `tfschema:"user_plane_access_ipv4_address"`
	UserPlaneAccessIPv4Gateway           string                             `tfschema:"user_plane_access_ipv4_gateway"`
	UserPlaneAccessIPv4Subnet            string                             `tfschema:"user_plane_access_ipv4_subnet"`
	UserPlaneAccessName                  string                             `tfschema:"user_plane_access_name"`
}

type NaptConfigurationDataSourceModel struct {
	PinholeLimits           int64                      `tfschema:"pinhole_maximum_number"`
	IcmpPinholeTimeout      int64                      `tfschema:"icmp_pinhole_timeout_in_seconds"`
	TcpPinholeTimeout       int64                      `tfschema:"tcp_pinhole_timeout_in_seconds"`
	UdpPinholeTimeout       int64                      `tfschema:"udp_pinhole_timeout_in_seconds"`
	PortRange               []PortRangeDataSourceModel `tfschema:"port_range"`
	TcpReuseMinimumHoldTime int64                      `tfschema:"tcp_port_reuse_minimum_hold_time_in_seconds"`
	UdpReuseMinimumHoldTime int64                      `tfschema:"udp_port_reuse_minimum_hold_time_in_seconds"`
}

type PortRangeDataSourceModel struct {
	Maximum int64 `tfschema:"maximum"`
	Minimum int64 `tfschema:"minimum"`
}

var _ sdk.DataSource = AttachedDataNetworkDataSource{}

func (r AttachedDataNetworkDataSource) ResourceType() string {
	return "azurerm_mobile_network_attached_data_network"
}

func (r AttachedDataNetworkDataSource) ModelObject() interface{} {
	return &AttachedDataNetworkDataSourceModel{}
}

func (r AttachedDataNetworkDataSource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return attacheddatanetwork.ValidateAttachedDataNetworkID
}

func (r AttachedDataNetworkDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"mobile_network_data_network_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"mobile_network_packet_core_data_plane_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: packetcoredataplane.ValidatePacketCoreDataPlaneID,
		},
	}
}

func (r AttachedDataNetworkDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{

		"location": commonschema.LocationComputed(),

		"dns_addresses": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"network_address_port_translation": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"pinhole_maximum_number": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
					},

					"icmp_pinhole_timeout_in_seconds": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},

					"tcp_pinhole_timeout_in_seconds": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},

					"udp_pinhole_timeout_in_seconds": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},

					"port_range": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"maximum": {
									Type:     pluginsdk.TypeInt,
									Computed: true,
								},

								"minimum": {
									Type:     pluginsdk.TypeInt,
									Computed: true,
								},
							},
						},
					},

					"tcp_port_reuse_minimum_hold_time_in_seconds": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},

					"udp_port_reuse_minimum_hold_time_in_seconds": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},
				},
			},
		},

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

		"user_equipment_address_pool_prefixes": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"user_equipment_static_address_pool_prefixes": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"tags": commonschema.TagsDataSource(),
	}
}

func (r AttachedDataNetworkDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var inputModel AttachedDataNetworkModel
			if err := metadata.Decode(&inputModel); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.MobileNetwork.AttachedDataNetworkClient

			packetCoreDataPlaneId, err := packetcoredataplane.ParsePacketCoreDataPlaneID(inputModel.MobileNetworkPacketCoreDataPlaneId)
			if err != nil {
				return err
			}

			id := attacheddatanetwork.NewAttachedDataNetworkID(packetCoreDataPlaneId.SubscriptionId, packetCoreDataPlaneId.ResourceGroupName, packetCoreDataPlaneId.PacketCoreControlPlaneName, packetCoreDataPlaneId.PacketCoreDataPlaneName, inputModel.MobileNetworkDataNetworkName)
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			state := AttachedDataNetworkDataSourceModel{
				MobileNetworkDataNetworkName:       id.AttachedDataNetworkName,
				MobileNetworkPacketCoreDataPlaneId: packetcoredataplane.NewPacketCoreDataPlaneID(id.SubscriptionId, id.ResourceGroupName, id.PacketCoreControlPlaneName, id.PacketCoreDataPlaneName).ID(),
			}

			if model := resp.Model; model != nil {
				props := model.Properties

				state.Location = location.Normalize(model.Location)
				state.DnsAddresses = props.DnsAddresses
				state.NaptConfiguration = flattenDataSourceNaptConfiguration(props.NaptConfiguration)
				state.UserEquipmentAddressPoolPrefix = pointer.From(props.UserEquipmentAddressPoolPrefix)
				state.UserEquipmentStaticAddressPoolPrefix = pointer.From(props.UserEquipmentStaticAddressPoolPrefix)
				state.UserPlaneAccessIPv4Address = pointer.From(props.UserPlaneDataInterface.IPv4Address)
				state.UserPlaneAccessIPv4Gateway = pointer.From(props.UserPlaneDataInterface.IPv4Gateway)
				state.UserPlaneAccessIPv4Subnet = pointer.From(props.UserPlaneDataInterface.IPv4Subnet)
				state.UserPlaneAccessName = pointer.From(props.UserPlaneDataInterface.Name)
				state.Tags = tags.Flatten(model.Tags)
			}

			metadata.SetID(id)

			return metadata.Encode(&state)
		},
	}
}

func flattenDataSourceNaptConfiguration(input *attacheddatanetwork.NaptConfiguration) []NaptConfigurationDataSourceModel {
	if input == nil {
		return []NaptConfigurationDataSourceModel{}
	}
	output := NaptConfigurationDataSourceModel{}

	output.PinholeLimits = pointer.From(input.PinholeLimits)

	if input.PinholeTimeouts != nil {
		output.IcmpPinholeTimeout = pointer.From(input.PinholeTimeouts.Icmp)
		output.TcpPinholeTimeout = pointer.From(input.PinholeTimeouts.Tcp)
		output.UdpPinholeTimeout = pointer.From(input.PinholeTimeouts.Udp)
	}

	output.PortRange = flattenDataSourcePortRange(input.PortRange)

	if input.PortReuseHoldTime != nil {
		output.TcpReuseMinimumHoldTime = pointer.From(input.PortReuseHoldTime.Tcp)
		output.UdpReuseMinimumHoldTime = pointer.From(input.PortReuseHoldTime.Udp)
	}

	return []NaptConfigurationDataSourceModel{output}
}

func flattenDataSourcePortRange(input *attacheddatanetwork.PortRange) []PortRangeDataSourceModel {
	if input == nil {
		return []PortRangeDataSourceModel{}
	}

	output := PortRangeDataSourceModel{}

	output.Maximum = pointer.From(input.MaxPort)

	output.Minimum = pointer.From(input.MinPort)

	return []PortRangeDataSourceModel{output}
}
