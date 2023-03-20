package mobilenetwork

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mobilenetwork/2022-11-01/packetcorecontrolplane"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type PacketCoreControlPlaneDataSource struct{}

var _ sdk.DataSource = PacketCoreControlPlaneDataSource{}

func (r PacketCoreControlPlaneDataSource) ResourceType() string {
	return "azurerm_mobile_network_packet_core_control_plane"
}

func (r PacketCoreControlPlaneDataSource) ModelObject() interface{} {
	return &PacketCoreControlPlaneModel{}
}

func (r PacketCoreControlPlaneDataSource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return packetcorecontrolplane.ValidatePacketCoreControlPlaneID
}

func (r PacketCoreControlPlaneDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),
	}
}

func (r PacketCoreControlPlaneDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{

		"location": commonschema.LocationComputed(),

		"control_plane_access_interface": interfacePropertiesSchemaDataSource(),

		"site_ids": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"sku": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"user_equipment_mtu_in_bytes": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"core_network_technology": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"platform": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"edge_device_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"azure_stack_hci_cluster_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"azure_arc_connected_cluster_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"custom_location_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"type": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"identity": commonschema.UserAssignedIdentityComputed(),

		"interop_json": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"local_diagnostics_access": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*schema.Schema{
					"authentication_type": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"https_server_certificate_url": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"tags": commonschema.TagsDataSource(),

		"version": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r PacketCoreControlPlaneDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var metaModel PacketCoreControlPlaneModel
			if err := metadata.Decode(&metaModel); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.MobileNetwork.PacketCoreControlPlaneClient
			subscriptionId := metadata.Client.Account.SubscriptionId
			id := packetcorecontrolplane.NewPacketCoreControlPlaneID(subscriptionId, metaModel.ResourceGroupName, metaModel.Name)

			resp, err := client.Get(ctx, id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if resp.Model == nil {
				return fmt.Errorf("retrieving %s: model was nil", id)
			}
			model := *resp.Model

			state := PacketCoreControlPlaneModel{
				Name:              id.PacketCoreControlPlaneName,
				ResourceGroupName: id.ResourceGroupName,
				Location:          location.Normalize(model.Location),
			}

			if model.Properties.UeMtu != nil {
				state.UeMtu = *model.Properties.UeMtu
			}

			identityValue, err := identity.FlattenLegacySystemAndUserAssignedMap(model.Identity)
			if err != nil {
				return fmt.Errorf("flattening `identity`: %+v", err)
			}

			if err := metadata.ResourceData.Set("identity", identityValue); err != nil {
				return fmt.Errorf("setting `identity`: %+v", err)
			}

			properties := &model.Properties
			state.ControlPlaneAccessInterface = flattenPacketCoreControlPlaneInterfacePropertiesModel(properties.ControlPlaneAccessInterface)

			if properties.CoreNetworkTechnology != nil {
				state.CoreNetworkTechnology = string(*properties.CoreNetworkTechnology)
			}

			if properties.InteropSettings != nil && *properties.InteropSettings != nil {
				interopSettingsValue, err := json.Marshal(*properties.InteropSettings)
				if err != nil {
					return err
				}

				state.InteropSettings = string(interopSettingsValue)
			}

			state.LocalDiagnosticsAccess = flattenLocalPacketCoreControlLocalDiagnosticsAccessConfiguration(properties.LocalDiagnosticsAccess)

			state.SiteIds = flattenPacketCoreControlPlaneSites(properties.Sites)

			platformValue := flattenPlatformConfigurationModel(properties.Platform)
			state.Platform = platformValue

			state.Sku = string(properties.Sku)

			if properties.Version != nil {
				state.Version = *properties.Version
			}
			if model.Tags != nil {
				state.Tags = *model.Tags
			}
			metadata.SetID(id)

			return metadata.Encode(&state)
		},
	}
}
