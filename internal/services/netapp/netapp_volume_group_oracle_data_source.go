// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package netapp

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2025-01-01/volumegroups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	netAppModels "github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/models"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ sdk.DataSource = NetAppVolumeGroupOracleDataSource{}

type NetAppVolumeGroupOracleDataSource struct{}

func (r NetAppVolumeGroupOracleDataSource) ResourceType() string {
	return "azurerm_netapp_volume_group_oracle"
}

func (r NetAppVolumeGroupOracleDataSource) ModelObject() interface{} {
	return &netAppModels.NetAppVolumeGroupOracleDataSourceModel{}
}

func (r NetAppVolumeGroupOracleDataSource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return volumegroups.ValidateVolumeGroupID
}

func (r NetAppVolumeGroupOracleDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.VolumeGroupName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"account_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.AccountName,
		},
	}
}

func (r NetAppVolumeGroupOracleDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.LocationComputed(),

		"group_description": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"application_identifier": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"volume": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"name": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"capacity_pool_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"proximity_placement_group_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"zone": commonschema.ZoneSingleComputed(),

					"volume_spec_name": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"volume_path": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"service_level": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"subnet_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"protocols": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"security_style": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"storage_quota_in_gb": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},

					"throughput_in_mibps": {
						Type:     pluginsdk.TypeFloat,
						Required: true,
					},

					"export_policy_rule": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"rule_index": {
									Type:     pluginsdk.TypeInt,
									Computed: true,
								},

								"allowed_clients": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},

								"nfsv3_enabled": {
									Type:     pluginsdk.TypeBool,
									Computed: true,
								},

								"nfsv41_enabled": {
									Type:     pluginsdk.TypeBool,
									Computed: true,
								},

								"unix_read_only": {
									Type:     pluginsdk.TypeBool,
									Computed: true,
								},

								"unix_read_write": {
									Type:     pluginsdk.TypeBool,
									Computed: true,
								},

								"root_access_enabled": {
									Type:     pluginsdk.TypeBool,
									Computed: true,
								},
							},
						},
					},

					"tags": commonschema.TagsDataSource(),

					"snapshot_directory_visible": {
						Type:     pluginsdk.TypeBool,
						Computed: true,
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
									Type:     pluginsdk.TypeString,
									Computed: true,
								},

								"remote_volume_location": commonschema.LocationComputed(),

								"remote_volume_resource_id": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},

								"replication_frequency": {
									Type:     pluginsdk.TypeString,
									Computed: true,
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
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
							},
						},
					},

					"encryption_key_source": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"key_vault_private_endpoint_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"network_features": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},
	}
}

func (r NetAppVolumeGroupOracleDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.NetApp.VolumeGroupClient

			var state netAppModels.NetAppVolumeGroupOracleDataSourceModel
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

					volumes, err := flattenNetAppVolumeGroupOracleVolumes(ctx, props.Volumes, metadata)
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
