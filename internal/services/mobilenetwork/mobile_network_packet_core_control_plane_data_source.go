// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mobilenetwork

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
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

		"control_plane_access_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"control_plane_access_ipv4_address": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"control_plane_access_ipv4_subnet": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"control_plane_access_ipv4_gateway": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

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

					"stack_hci_cluster_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"arc_kubernetes_cluster_id": {
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

		"interoperability_settings_json": {
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

		"software_version": {
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

			state := PacketCoreControlPlaneModel{
				Name:              id.PacketCoreControlPlaneName,
				ResourceGroupName: id.ResourceGroupName,
			}

			if model := resp.Model; model != nil {
				state.Location = location.Normalize(model.Location)

				state.Identity, err = flattenMobileNetworkUserAssignedToNetworkLegacyIdentity(model.Identity)
				if err != nil {
					return fmt.Errorf("flattening `identity`: %+v", err)
				}

				properties := model.Properties

				state.UeMtu = pointer.From(model.Properties.UeMtu)

				state.ControlPlaneAccessIPv4Address = pointer.From(properties.ControlPlaneAccessInterface.IPv4Address)

				state.ControlPlaneAccessIPv4Gateway = pointer.From(properties.ControlPlaneAccessInterface.IPv4Gateway)

				state.ControlPlaneAccessIPv4Subnet = pointer.From(properties.ControlPlaneAccessInterface.IPv4Subnet)

				state.ControlPlaneAccessName = pointer.From(properties.ControlPlaneAccessInterface.Name)

				// it still needs a nil check because it needs to do type conversion
				if properties.CoreNetworkTechnology != nil {
					state.CoreNetworkTechnology = string(pointer.From(properties.CoreNetworkTechnology))
				}

				// Marshal on a nil interface{} may get random result.
				if properties.InteropSettings != nil && *properties.InteropSettings != nil {
					interopSettingsValue, err := json.Marshal(*properties.InteropSettings)
					if err != nil {
						return err
					}

					state.InteropSettings = string(interopSettingsValue)
				}

				state.LocalDiagnosticsAccess = flattenLocalPacketCoreControlLocalDiagnosticsAccessConfiguration(properties.LocalDiagnosticsAccess)

				state.SiteIds = flattenPacketCoreControlPlaneSites(properties.Sites)

				state.Platform = flattenPlatformConfigurationModel(properties.Platform)

				state.Sku = string(properties.Sku)

				state.SoftwareVersion = pointer.From(properties.Version)
				state.Tags = pointer.From(model.Tags)
			}

			metadata.SetID(id)

			return metadata.Encode(&state)
		},
	}
}
