package netapp

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2022-05-01/volumegroups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type NetAppVolumeGroupDataSourceModel struct {
	Name                  string                    `tfschema:"name"`
	ResourceGroupName     string                    `tfschema:"resource_group_name"`
	Location              string                    `tfschema:"location"`
	AccountName           string                    `tfschema:"account_name"`
	GroupDescription      string                    `tfschema:"group_description"`
	ApplicationType       string                    `tfschema:"application_type"`
	ApplicationIdentifier string                    `tfschema:"application_identifier"`
	Volumes               []NetAppVolumeGroupVolume `tfschema:"volume"`
}

var _ sdk.DataSource = NetAppVolumeGroupDataSource{}

type NetAppVolumeGroupDataSource struct{}

func (r NetAppVolumeGroupDataSource) ResourceType() string {
	return "azurerm_netapp_volume_group"
}

func (r NetAppVolumeGroupDataSource) ModelObject() interface{} {
	return &NetAppVolumeGroupDataSourceModel{}
}

func (r NetAppVolumeGroupDataSource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return volumegroups.ValidateVolumeGroupID
}

func (r NetAppVolumeGroupDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"account_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},
	}
}

func (r NetAppVolumeGroupDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.LocationComputed(),

		"group_description": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"application_type": {
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
				},
			},
		},
	}
}

func (r NetAppVolumeGroupDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {

			client := metadata.Client.NetApp.VolumeGroupClient

			var state NetAppVolumeGroupDataSourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := volumegroups.NewVolumeGroupID(metadata.Client.Account.SubscriptionId, state.ResourceGroupName, state.AccountName, state.Name)

			resp, err := client.VolumeGroupsGet(ctx, id)
			if err != nil {
				if resp.HttpResponse.StatusCode == http.StatusNotFound {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %v", id, err)
			}

			model := resp.Model
			if model == nil {
				return fmt.Errorf("retrieving %s: model was nil", id)
			}

			state.Location = location.Normalize(*model.Location)
			state.ApplicationIdentifier = *model.Properties.GroupMetaData.ApplicationIdentifier
			state.ApplicationType = string(*model.Properties.GroupMetaData.ApplicationType)
			state.GroupDescription = *model.Properties.GroupMetaData.GroupDescription

			volumes, err := flattenNetAppVolumeGroupVolumes(ctx, model.Properties.Volumes, metadata)
			if err != nil {
				return fmt.Errorf("setting `volume`: %+v", err)
			}

			state.Volumes = volumes

			metadata.SetID(id)

			return metadata.Encode(&state)
		},
	}
}
