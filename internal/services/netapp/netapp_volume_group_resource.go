package netapp

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2022-01-01/snapshots"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2022-01-01/volumegroups"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	netAppValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type NetAppVolumeGroupResource struct{}

type NetAppVolumeGroupModel struct {
	Name                  string                    `tfschema:"name"`
	ResourceGroupName     string                    `tfschema:"resource_group_name"`
	Location              string                    `tfschema:"location"`
	AccountName           string                    `tfschema:"account_name"`
	GroupDescription      string                    `tfschema:"group_description"`
	ApplicationType       string                    `tfschema:"application_type"`
	ApplicationIdentifier string                    `tfschema:"application_identifier"`
	DeploymentSpecId      string                    `tfschema:"deployment_spec_id"`
	Volumes               []NetAppVolumeGroupVolume `tfschema:"volume"`
}

var _ sdk.Resource = NetAppVolumeGroupResource{}

func (r NetAppVolumeGroupResource) ModelObject() interface{} {
	return &NetAppVolumeGroupModel{}
}

func (r NetAppVolumeGroupResource) ResourceType() string {
	return "azurerm_netapp_volume_group"
}

func (r NetAppVolumeGroupResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return volumegroups.ValidateVolumeGroupID
}

func (r NetAppVolumeGroupResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"resource_group_name": azure.SchemaResourceGroupName(),

		"location": azure.SchemaLocation(),

		"account_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: netAppValidate.AccountName,
		},

		"group_description": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"application_type": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(volumegroups.ApplicationTypeSAPNegativeHANA),
			}, false),
		},

		"application_identifier": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringLenBetween(1, 3),
		},

		"deployment_spec_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsUUID,
		},

		"volume": {
			Type:     pluginsdk.TypeSet,
			Required: true,
			ForceNew: true,
			MinItems: 5,
			MaxItems: 5,
			Elem: &pluginsdk.Resource{
				//Schema: netAppVolumeGroupVolumeSchema(),
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: netAppValidate.VolumeName,
					},

					"capacity_pool_id": {
						Type:             pluginsdk.TypeString,
						Optional:         true,
						Computed:         true,
						ForceNew:         true,
						DiffSuppressFunc: suppress.CaseDifference,
						ValidateFunc:     azure.ValidateResourceID,
					},

					"proximity_placement_group_id": {
						Type:             pluginsdk.TypeString,
						Required:         true,
						ForceNew:         true,
						DiffSuppressFunc: suppress.CaseDifference,
						ValidateFunc:     azure.ValidateResourceID,
					},

					"volume_spec_name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"volume_path": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: netAppValidate.VolumePath,
					},

					"service_level": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ForceNew: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(volumegroups.ServiceLevelPremium),
							string(volumegroups.ServiceLevelStandard),
							string(volumegroups.ServiceLevelUltra),
						}, false),
					},

					"subnet_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: azure.ValidateResourceID,
					},

					"create_from_snapshot_resource_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Computed:     true,
						ForceNew:     true,
						ValidateFunc: snapshots.ValidateSnapshotID,
					},

					"network_features": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
						ForceNew: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(volumegroups.NetworkFeaturesBasic),
							string(volumegroups.NetworkFeaturesStandard),
						}, false),
					},

					"protocols": {
						Type:     pluginsdk.TypeSet,
						ForceNew: true,
						Optional: true,
						Computed: true,
						MaxItems: 2,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
							ValidateFunc: validation.StringInSlice([]string{
								"NFSv3",
								"NFSv4.1",
								"CIFS",
							}, false),
						},
					},

					"security_style": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ForceNew: true,
						Computed: true,
						ValidateFunc: validation.StringInSlice([]string{
							"Unix", // Using hardcoded values instead of SDK enum since no matter what case is passed,
							"Ntfs", // ANF changes casing to Pascal case in the backend. Please refer to https://github.com/Azure/azure-sdk-for-go/issues/14684
						}, false),
					},

					"storage_quota_in_gb": {
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntBetween(100, 102400),
					},

					"throughput_in_mibps": {
						Type:     pluginsdk.TypeFloat,
						Optional: true,
						Computed: true,
					},

					"export_policy_rule": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 5,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"rule_index": {
									Type:         pluginsdk.TypeInt,
									Required:     true,
									ValidateFunc: validation.IntBetween(1, 5),
								},

								"allowed_clients": {
									Type:     pluginsdk.TypeSet,
									Required: true,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validate.CIDR,
									},
								},

								"protocols_enabled": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Computed: true,
									MaxItems: 1,
									MinItems: 1,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
										ValidateFunc: validation.StringInSlice([]string{
											"NFSv3",
											"NFSv4.1",
											"CIFS",
										}, false),
									},
								},

								"unix_read_only": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									Computed: true,
								},

								"unix_read_write": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									Computed: true,
								},

								"root_access_enabled": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									Computed: true,
								},

								"kerberos5_read_only": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									Computed: true,
								},

								"kerberos5_read_write": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									Computed: true,
								},

								"kerberos5i_read_only": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									Computed: true,
								},

								"kerberos5i_read_write": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									Computed: true,
								},

								"kerberos5p_read_only": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									Computed: true,
								},

								"kerberos5p_read_write": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									Computed: true,
								},
							},
						},
					},

					"tags": commonschema.Tags(),

					"mount_ip_addresses": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"snapshot_directory_visible": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Computed: true,
					},

					"data_protection_replication": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						ForceNew: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"endpoint_type": {
									Type:     pluginsdk.TypeString,
									Optional: true,
									Default:  "dst",
									ValidateFunc: validation.StringInSlice([]string{
										"dst",
									}, false),
								},

								"remote_volume_location": azure.SchemaLocation(),

								"remote_volume_resource_id": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: azure.ValidateResourceID,
								},

								"replication_frequency": {
									Type:     pluginsdk.TypeString,
									Required: true,
									ValidateFunc: validation.StringInSlice([]string{
										"10minutes",
										"daily",
										"hourly",
									}, false),
								},
							},
						},
					},

					"data_protection_snapshot_policy": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"snapshot_policy_id": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: azure.ValidateResourceID,
								},
							},
						},
					},
				},
			},
		},
	}
}

func (r NetAppVolumeGroupResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		/*
			TODO - This section is for `Computed: true` only items, i.e. useful values that are returned by the
			datasource that can be used as outputs or passed programmatically to other resources or data sources.

			TODO (pmarques) - use this for first level attributes when Volume resource gets migrated to tfschema
		*/
	}
}

func (r NetAppVolumeGroupResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 90 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.NetApp.VolumeGroupClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model NetAppVolumeGroupModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := volumegroups.NewVolumeGroupID(subscriptionId, model.ResourceGroupName, model.AccountName, model.Name)

			metadata.Logger.Infof("Import check for %s", id)
			existing, err := client.VolumeGroupsGet(ctx, id)
			if err != nil && existing.HttpResponse.StatusCode != http.StatusNotFound {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}

			if existing.Model != nil && existing.Model.Id != nil && *existing.Model.Id != "" {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			applicationType := volumegroups.ApplicationType(model.ApplicationType)

			volumeList, err := expandNetAppVolumeGroupVolumes(model.Volumes, id)
			if err != nil {
				return err
			}

			parameters := volumegroups.VolumeGroupDetails{
				Location: utils.String(location.Normalize(model.Location)),
				Properties: &volumegroups.VolumeGroupProperties{
					GroupMetaData: &volumegroups.VolumeGroupMetaData{
						GroupDescription:      utils.String(model.GroupDescription),
						ApplicationType:       &applicationType,
						ApplicationIdentifier: utils.String(model.ApplicationIdentifier),
						DeploymentSpecId:      utils.String(model.DeploymentSpecId),
					},
					Volumes: volumeList,
				},
			}

			err = client.VolumeGroupsCreateThenPoll(ctx, id, parameters)
			if err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			// TODO: Check if this is necessary for volume groups
			// // Waiting for volume be completely provisioned
			// if err := waitForVolumeCreateOrUpdate(ctx, client, id); err != nil {
			// 	return err
			// }

			metadata.SetID(id)

			return nil
		},
	}
}

func (r NetAppVolumeGroupResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {

			client := metadata.Client.NetApp.VolumeGroupClient

			id, err := volumegroups.ParseVolumeGroupID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Infof("Decoding state for %s", id)
			var state NetAppVolumeGroupModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			existing, err := client.VolumeGroupsGet(ctx, *id)
			if err != nil {
				if existing.HttpResponse.StatusCode == http.StatusNotFound {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %v", id, err)
			}

			model := NetAppVolumeGroupModel{
				Name:              id.VolumeGroupName,
				AccountName:       id.AccountName,
				Location:          location.NormalizeNilable(existing.Model.Location),
				ResourceGroupName: id.ResourceGroupName,
			}

			if props := existing.Model.Properties; props != nil {
				model.GroupDescription = utils.NormalizeNilableString(props.GroupMetaData.GroupDescription)
				model.ApplicationIdentifier = utils.NormalizeNilableString(props.GroupMetaData.ApplicationIdentifier)
				model.DeploymentSpecId = utils.NormalizeNilableString(props.GroupMetaData.DeploymentSpecId)
				model.ApplicationType = string(*props.GroupMetaData.ApplicationType)
				model.DeploymentSpecId = state.DeploymentSpecId

				volumes, err := flattenNetAppVolumeGroupVolumes(props.Volumes, state.Volumes)
				if err != nil {
					return fmt.Errorf("setting `volume`: %+v", err)
				}

				model.Volumes = volumes
			}

			return metadata.Encode(&model)
		},
	}
}

func (r NetAppVolumeGroupResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			// TODO - Delete Func
			return nil
		},
	}
}
