package netapp

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2022-01-01/volumegroups"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	netAppValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
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
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsUUID,
			Default:      "20542149-bfca-5618-1879-9863dc6767f1",
		},

		"volume": {
			Type:     pluginsdk.TypeSet,
			Required: true,
			ForceNew: true,
			MinItems: 5,
			MaxItems: 5,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: netAppValidate.VolumeName,
					},

					"capacity_pool_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: azure.ValidateResourceID,
					},

					"proximity_placement_group_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: azure.ValidateResourceID,
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

					"protocols": {
						Type:     pluginsdk.TypeSet,
						ForceNew: true,
						Required: true,
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
						Required: true,
						ForceNew: true,
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
						Required: true,
					},

					"export_policy_rule": {
						Type:     pluginsdk.TypeSet,
						Required: true,
						MaxItems: 5,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"rule_index": {
									Type:         pluginsdk.TypeInt,
									Required:     true,
									ValidateFunc: validation.IntBetween(1, 5),
								},

								"allowed_clients": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},

								"cifs_enabled": {
									Type:     pluginsdk.TypeBool,
									Required: true,
								},

								"nfsv3_enabled": {
									Type:     pluginsdk.TypeBool,
									Required: true,
								},

								"nfsv41_enabled": {
									Type:     pluginsdk.TypeBool,
									Required: true,
								},

								"unix_read_only": {
									Type:     pluginsdk.TypeBool,
									Required: true,
								},

								"unix_read_write": {
									Type:     pluginsdk.TypeBool,
									Required: true,
								},

								"root_access_enabled": {
									Type:     pluginsdk.TypeBool,
									Required: true,
								},

								"kerberos5_read_only": {
									Type:     pluginsdk.TypeBool,
									Required: true,
								},

								"kerberos5_read_write": {
									Type:     pluginsdk.TypeBool,
									Required: true,
								},

								"kerberos5i_read_only": {
									Type:     pluginsdk.TypeBool,
									Required: true,
								},

								"kerberos5i_read_write": {
									Type:     pluginsdk.TypeBool,
									Required: true,
								},

								"kerberos5p_read_only": {
									Type:     pluginsdk.TypeBool,
									Required: true,
								},

								"kerberos5p_read_write": {
									Type:     pluginsdk.TypeBool,
									Required: true,
								},
							},
						},
					},

					"tags": commonschema.Tags(),

					"snapshot_directory_visible": {
						Type:     pluginsdk.TypeBool,
						Required: true,
						ForceNew: true,
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
			Set: resourceVolumeGroupVolumeListHash,
		},
	}
}

func resourceVolumeGroupVolumeListHash(v interface{}) int {
	// Computed = true items must be out of this

	var buf bytes.Buffer

	if m, ok := v.(map[string]interface{}); ok {
		buf.WriteString(fmt.Sprintf("%s-", m["name"].(string)))
		buf.WriteString(fmt.Sprintf("%s-", m["proximity_placement_group_id"].(string)))
		buf.WriteString(fmt.Sprintf("%s-", m["volume_spec_name"].(string)))
		buf.WriteString(fmt.Sprintf("%s-", m["volume_path"].(string)))
		buf.WriteString(fmt.Sprintf("%s-", m["service_level"].(string)))

		if protocols, ok := m["protocols"].([]interface{}); ok {
			for _, item := range protocols {
				v := item.(string)
				buf.WriteString(fmt.Sprintf("%s-", v))
			}
		}

		buf.WriteString(fmt.Sprintf("%s-", m["security_style"].(string)))
		buf.WriteString(fmt.Sprintf("%d-", m["storage_quota_in_gb"].(int)))
		buf.WriteString(fmt.Sprintf("%f-", m["throughput_in_mibps"].(float64)))
		buf.WriteString(fmt.Sprintf("%t-", m["snapshot_directory_visible"].(bool)))

		if exportPolicies, ok := m["export_policy_rule"].([]interface{}); ok {
			for _, item := range exportPolicies {
				v := item.(map[string]interface{})
				if ruleIndex, ok := v["rule_index"].(int); ok {
					buf.WriteString(fmt.Sprintf("%d-", ruleIndex))
				}
				if allowedClients, ok := v["allowed_clients"].(string); ok {
					buf.WriteString(fmt.Sprintf("%s-", allowedClients))
				}
				if cifsEnabled, ok := v["cifs_enabled"].(bool); ok {
					buf.WriteString(fmt.Sprintf("%t-", cifsEnabled))
				}
				if nfsv3Enabled, ok := v["nfsv3_enabled"].(bool); ok {
					buf.WriteString(fmt.Sprintf("%t-", nfsv3Enabled))
				}
				if nfsv41Enabled, ok := v["nfsv41_enabled"].(bool); ok {
					buf.WriteString(fmt.Sprintf("%t-", nfsv41Enabled))
				}
				if unixReadOnly, ok := v["unix_read_only"].(bool); ok {
					buf.WriteString(fmt.Sprintf("%t-", unixReadOnly))
				}
				if unixReadWrite, ok := v["unix_read_write"].(bool); ok {
					buf.WriteString(fmt.Sprintf("%t-", unixReadWrite))
				}
				if rootAccessEnabled, ok := v["root_access_enabled"].(bool); ok {
					buf.WriteString(fmt.Sprintf("%t-", rootAccessEnabled))
				}
				if kerberos5ReadOnly, ok := v["kerberos5_read_only"].(bool); ok {
					buf.WriteString(fmt.Sprintf("%t-", kerberos5ReadOnly))
				}
				if kerberos5ReadWrite, ok := v["kerberos5_read_write"].(bool); ok {
					buf.WriteString(fmt.Sprintf("%t-", kerberos5ReadWrite))
				}
				if kerberos5iReadOnly, ok := v["kerberos5i_read_only"].(bool); ok {
					buf.WriteString(fmt.Sprintf("%t-", kerberos5iReadOnly))
				}
				if kerberos5iReadWrite, ok := v["kerberos5i_read_write"].(bool); ok {
					buf.WriteString(fmt.Sprintf("%t-", kerberos5iReadWrite))
				}
				if kerberos5pReadOnly, ok := v["kerberos5p_read_only"].(bool); ok {
					buf.WriteString(fmt.Sprintf("%t-", kerberos5pReadOnly))
				}
				if kerberos5pReadWrite, ok := v["kerberos5p_read_write"].(bool); ok {
					buf.WriteString(fmt.Sprintf("%t-", kerberos5pReadWrite))
				}
			}
		}

		if tags, ok := m["tags"].([]interface{}); ok {
			for _, item := range tags {
				i := item.(map[string]interface{})
				for k, v := range i {
					buf.WriteString(fmt.Sprintf("%s-%s-", k, v))
				}
			}
		}

		if dpReplication, ok := m["data_protection_replication"].([]interface{}); ok {
			for _, item := range dpReplication {
				v := item.(map[string]interface{})
				if endpointType, ok := v["endpoint_type"].(string); ok {
					buf.WriteString(fmt.Sprintf("%s-", endpointType))
				}
				if remoteVolumeLocation, ok := v["remote_volume_location"].(string); ok {
					buf.WriteString(fmt.Sprintf("%s-", remoteVolumeLocation))
				}
				if remoteVolumeResourceId, ok := v["remote_volume_resource_id"].(string); ok {
					buf.WriteString(fmt.Sprintf("%s-", remoteVolumeResourceId))
				}
				if replicationFrequency, ok := v["replication_frequency"].(string); ok {
					buf.WriteString(fmt.Sprintf("%s-", replicationFrequency))
				}
			}
		}

		if dpSnapshotPolicy, ok := m["data_protection_snapshot_policy"].([]interface{}); ok {
			for _, item := range dpSnapshotPolicy {
				v := item.(map[string]interface{})
				if snapshotPolicyId, ok := v["snapshot_policy_id"].(string); ok {
					buf.WriteString(fmt.Sprintf("%s-", snapshotPolicyId))
				}
			}
		}
	}

	return pluginsdk.HashString(buf.String())
}

func (r NetAppVolumeGroupResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		/*
			TODO - This section is for `Computed: true` only items, i.e. useful values that are returned by the
			datasource that can be used as outputs or passed programmatically to other resources or data sources.

			NOTE: Not applicable for this resource type
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

				if state.DeploymentSpecId != "" {
					model.DeploymentSpecId = state.DeploymentSpecId
				} else {
					// Setting a default value here to overcome issue with SDK
					// not returning this value back from Azure
					// This is the only supported value for the time being and
					// will be fixed by ANF team if it introduces a new SpecId
					// option.
					model.DeploymentSpecId = "20542149-bfca-5618-1879-9863dc6767f1"
				}

				volumes, err := flattenNetAppVolumeGroupVolumes(ctx, props.Volumes, metadata)
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
		Timeout: 120 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {

			client := metadata.Client.NetApp.VolumeGroupClient

			id, err := volumegroups.ParseVolumeGroupID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.VolumeGroupsGet(ctx, *id)
			if err != nil {
				if existing.HttpResponse.StatusCode == http.StatusNotFound {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %v", id, err)
			}

			// Removing volumes before deleting volume group
			if props := existing.Model.Properties; props != nil {
				if volumeList := props.Volumes; volumeList != nil {
					for _, volume := range *volumeList {
						if err := deleteVolume(ctx, metadata, *volume.Id); err != nil {
							return fmt.Errorf("deleting `volume`: %+v", err)
						}
					}
				}
			}

			// Removing Volume Group
			if err = client.VolumeGroupsDeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}
