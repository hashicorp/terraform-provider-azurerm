package mobilenetwork

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mobilenetwork/2022-11-01/attacheddatanetwork"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mobilenetwork/2022-11-01/packetcoredataplane"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type AttachedDataNetworkDataSource struct{}

var _ sdk.DataSource = AttachedDataNetworkDataSource{}

func (r AttachedDataNetworkDataSource) ResourceType() string {
	return "azurerm_mobile_network_attached_data_network"
}

func (r AttachedDataNetworkDataSource) ModelObject() interface{} {
	return &AttachedDataNetworkModel{}
}

func (r AttachedDataNetworkDataSource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return attacheddatanetwork.ValidateAttachedDataNetworkID
}

func (r AttachedDataNetworkDataSource) Arguments() map[string]*pluginsdk.Schema {
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

		"network_address_port_translation_configuration": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"enabled": {
						Type:     pluginsdk.TypeBool,
						Computed: true,
					},

					"pinhole_maximum_number": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
					},

					"pinhole_timeouts_in_seconds": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"icmp": {
									Type:     pluginsdk.TypeInt,
									Computed: true,
								},

								"tcp": {
									Type:     pluginsdk.TypeInt,
									Computed: true,
								},

								"udp": {
									Type:     pluginsdk.TypeInt,
									Computed: true,
								},
							},
						},
					},

					"port_range": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"max_port": {
									Type:     pluginsdk.TypeInt,
									Computed: true,
								},

								"min_port": {
									Type:     pluginsdk.TypeInt,
									Computed: true,
								},
							},
						},
					},

					"port_reuse_minimum_hold_time_in_seconds": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"tcp": {
									Type:     pluginsdk.TypeInt,
									Computed: true,
								},

								"udp": {
									Type:     pluginsdk.TypeInt,
									Computed: true,
								},
							},
						},
					},
				},
			},
		},

		"user_plane_data_interface": interfacePropertiesSchemaDataSource(),

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

			id := attacheddatanetwork.NewAttachedDataNetworkID(packetCoreDataPlaneId.SubscriptionId, packetCoreDataPlaneId.ResourceGroupName, packetCoreDataPlaneId.PacketCoreControlPlaneName, packetCoreDataPlaneId.PacketCoreDataPlaneName, inputModel.Name)
			if err != nil {
				return err
			}

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

			state := AttachedDataNetworkModel{
				Name:                               id.AttachedDataNetworkName,
				MobileNetworkPacketCoreDataPlaneId: packetcoredataplane.NewPacketCoreDataPlaneID(id.SubscriptionId, id.ResourceGroupName, id.PacketCoreControlPlaneName, id.PacketCoreDataPlaneName).ID(),
				Location:                           location.Normalize(model.Location),
			}

			properties := &model.Properties
			if properties.DnsAddresses != nil {
				state.DnsAddresses = properties.DnsAddresses
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

			state.UserPlaneDataInterface = flattenAttachedDataNetworkInterfacePropertiesModel(properties.UserPlaneDataInterface)

			if model.Tags != nil {
				state.Tags = *model.Tags
			}

			metadata.SetID(id)

			return metadata.Encode(&state)
		},
	}
}
