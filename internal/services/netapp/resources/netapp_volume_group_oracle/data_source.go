// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package netapp_volume_group_oracle

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2026-05-01/volumegroups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/models"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ sdk.DataSource = DataSource{}

type DataSource struct{}

func (r DataSource) ResourceType() string {
	return "azurerm_netapp_volume_group_oracle"
}

func (r DataSource) Description() string {
	return "Use this data source to access information about an existing Application Volume Group for Oracle application."
}

func (r DataSource) ModelObject() any {
	return &models.NetAppVolumeGroupOracleDataSourceModel{}
}

func (r DataSource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return volumegroups.ValidateVolumeGroupID
}

func (r DataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of this Application Volume Group for Oracle application.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.VolumeGroupName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"account_name": {
			Description:  "Name of the account where the application volume group belongs to.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.AccountName,
		},
	}
}

func (r DataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.LocationComputed(),

		"group_description": {
			Description: "Volume group description.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"application_identifier": {
			Description: "The application identifier.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"volume": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"id": {
						Description: "Volume ID.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},

					"name": {
						Description: "The name of this volume.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},

					"capacity_pool_id": {
						Description: "The ID of the Capacity Pool.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},

					"proximity_placement_group_id": {
						Description: "The ID of the Proximity Placement Group.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},

					"zone": commonschema.ZoneSingleComputed(),

					"volume_spec_name": {
						Description: "Volume spec name.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},

					"volume_path": {
						Description: "A unique file path for the volume.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},

					"service_level": {
						Description: "The target performance of the file system.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},

					"subnet_id": {
						Description: "The ID of the Subnet the NetApp Volume resides in.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},

					"protocols": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"security_style": {
						Description: "Volume security style.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},

					"storage_quota_in_gb": {
						Description: "The maximum Storage Quota allowed for a file system in Gigabytes.",
						Type:        pluginsdk.TypeInt,
						Computed:    true,
					},

					"throughput_in_mibps": {
						Description: "Throughput of this volume in Mibps.",
						Type:        pluginsdk.TypeFloat,
						Computed:    true,
					},

					"export_policy_rule": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"rule_index": {
									Description: "The index number of the rule.",
									Type:        pluginsdk.TypeInt,
									Computed:    true,
								},

								"allowed_clients": {
									Description: "A list of allowed clients IPv4 addresses.",
									Type:        pluginsdk.TypeString,
									Computed:    true,
								},

								"nfsv3_enabled": {
									Description: "Whether the NFSv3 protocol is enabled.",
									Type:        pluginsdk.TypeBool,
									Computed:    true,
								},

								"nfsv41_enabled": {
									Description: "Whether the NFSv4.1 is enabled.",
									Type:        pluginsdk.TypeBool,
									Computed:    true,
								},

								"unix_read_only": {
									Description: "Whether the unix filesystem is read-only.",
									Type:        pluginsdk.TypeBool,
									Computed:    true,
								},

								"unix_read_write": {
									Description: "Whether the unix filesystem is read-write.",
									Type:        pluginsdk.TypeBool,
									Computed:    true,
								},

								"root_access_enabled": {
									Description: "Whether root access is permitted to this volume.",
									Type:        pluginsdk.TypeBool,
									Computed:    true,
								},
							},
						},
					},

					"tags": commonschema.TagsDataSource(),

					"snapshot_directory_visible": {
						Description: "Whether the .snapshot (NFS clients) path of a volume is visible.",
						Type:        pluginsdk.TypeBool,
						Computed:    true,
					},

					"mount_ip_addresses": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"data_protection_replication": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"endpoint_type": {
									Description: "The endpoint type.",
									Type:        pluginsdk.TypeString,
									Computed:    true,
								},

								"remote_volume_location": commonschema.LocationComputed(),

								"remote_volume_resource_id": {
									Description: "Resource ID of the primary volume.",
									Type:        pluginsdk.TypeString,
									Computed:    true,
								},

								"replication_frequency": {
									Description: "Replication frequency.",
									Type:        pluginsdk.TypeString,
									Computed:    true,
								},
							},
						},
					},

					"data_protection_snapshot_policy": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"snapshot_policy_id": {
									Description: "Resource ID of the snapshot policy to apply to the volume.",
									Type:        pluginsdk.TypeString,
									Computed:    true,
								},
							},
						},
					},

					"encryption_key_source": {
						Description: "The encryption key source.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},

					"key_vault_private_endpoint_id": {
						Description: "The Private Endpoint ID for Key Vault when using customer managed keys.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},

					"network_features": {
						Description: "Network feature in use at the time of volume creation.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},
				},
			},
		},
	}
}

func (r DataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.NetApp.VolumeGroupClient

			var state models.NetAppVolumeGroupOracleDataSourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := volumegroups.NewVolumeGroupID(metadata.Client.Account.SubscriptionId, state.ResourceGroupName, state.AccountName, state.Name)

			resp, err := client.Get(ctx, id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}
				return fmt.Errorf("retrieving %s: %v", id, err)
			}

			if model := resp.Model; model != nil {
				state.Location = location.Normalize(pointer.From(model.Location))
				if props := model.Properties; props != nil {
					if groupMetaData := props.GroupMetaData; groupMetaData != nil {
						state.ApplicationIdentifier = pointer.From(groupMetaData.ApplicationIdentifier)
						state.GroupDescription = pointer.From(groupMetaData.GroupDescription)
					}

					volumes, err := helpers.FlattenNetAppVolumeGroupOracleVolumes(ctx, props.Volumes, metadata)
					if err != nil {
						return fmt.Errorf("setting `volume`: %+v", err)
					}
					state.Volumes = volumes
				}
			}

			metadata.SetID(id)

			return metadata.Encode(&state)
		},
	}
}
