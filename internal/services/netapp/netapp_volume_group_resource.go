package netapp

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2022-01-01/volumegroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2022-01-01/volumes"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2022-01-01/volumesreplication"
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
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
			Default:  "20542149-bfca-5618-1879-9863dc6767f1",
			ValidateFunc: validation.StringInSlice([]string{
				"20542149-bfca-5618-1879-9863dc6767f1",
			}, false),
		},

		"volume": {
			Type:     pluginsdk.TypeSet,
			Required: true,
			MinItems: 5,
			MaxItems: 5,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

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
						ValidateFunc: validation.StringInSlice([]string{
							"data",
							"data-backup",
							"log",
							"log-backup",
							"shared",
						}, false),
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
						MaxItems: 1,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
							ValidateFunc: validation.StringInSlice([]string{
								"NFSv3",
								"NFSv4.1",
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
			replicationClient := metadata.Client.NetApp.VolumeReplicationClient

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

			// Parse volume list to set secondary volumes for CRR
			for i, volumeCrr := range *volumeList {
				if volumeCrr.Properties.DataProtection != nil &&
					volumeCrr.Properties.DataProtection.Replication != nil &&
					*volumeCrr.Properties.DataProtection.Replication.EndpointType == volumegroups.EndpointTypeDst {

					// Modify volumeType as data protection type on main volumeList
					// so it gets created correctly as data protection volume
					(*volumeList)[i].Properties.VolumeType = utils.String("DataProtection")
				}
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

			// CRR - Authorizing secondaries from primary volumes
			for _, volumeCrr := range *volumeList {
				if volumeCrr.Properties.DataProtection != nil &&
					volumeCrr.Properties.DataProtection.Replication != nil &&
					*volumeCrr.Properties.DataProtection.Replication.EndpointType == volumegroups.EndpointTypeDst {

					// Getting secondary volume resource id
					secondaryId := volumes.NewVolumeID(subscriptionId,
						model.ResourceGroupName,
						model.AccountName,
						getResourceNameString(volumeCrr.Properties.CapacityPoolResourceId),
						getResourceNameString(volumeCrr.Name),
					)

					// Getting primary resource id
					primaryId, err := volumesreplication.ParseVolumeID(volumeCrr.Properties.DataProtection.Replication.RemoteVolumeResourceId)
					if err != nil {
						return err
					}

					// Authorizing
					if err = replicationClient.VolumesAuthorizeReplicationThenPoll(ctx, *primaryId, volumesreplication.AuthorizeRequest{
						RemoteVolumeResourceId: utils.String(secondaryId.ID()),
					},
					); err != nil {
						return fmt.Errorf("cannot authorize volume replication: %v", err)
					}

					// Wait for volume replication authorization to complete
					log.Printf("[DEBUG] Waiting for replication authorization on %s to complete", id)
					if err := waitForReplAuthorization(ctx, replicationClient, *primaryId); err != nil {
						return err
					}
				}
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (r NetAppVolumeGroupResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 120 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			volumeClient := metadata.Client.NetApp.VolumeClient

			id, err := volumes.ParseVolumeID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Infof("Decoding state for %s", id)
			var state NetAppVolumeGroupModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			metadata.Logger.Infof("Updating %s", id)

			if metadata.ResourceData.HasChange("volume") {

				// Iterating over each volume and performing individual patch
				for i := 0; i < metadata.ResourceData.Get("volume.#").(int); i++ {

					// Checking if individual volume has a change
					volumeItem := fmt.Sprintf("volume.%v", i)

					if metadata.ResourceData.HasChange(volumeItem) {

						update := volumes.VolumePatch{
							Properties: &volumes.VolumePatchProperties{},
						}

						if metadata.ResourceData.HasChange(fmt.Sprintf("%v.storage_quota_in_gb", volumeItem)) {
							storageQuotaInBytes := int64(metadata.ResourceData.Get(fmt.Sprintf("%v.storage_quota_in_gb", volumeItem)).(int) * 1073741824)
							update.Properties.UsageThreshold = utils.Int64(storageQuotaInBytes)
						}

						if metadata.ResourceData.HasChange(fmt.Sprintf("%v.export_policy_rule", volumeItem)) {
							exportPolicyRuleRaw := metadata.ResourceData.Get(fmt.Sprintf("%v.export_policy_rule", volumeItem)).([]interface{})
							exportPolicyRule := expandNetAppVolumeGroupVolumeExportPolicyRulePatch(exportPolicyRuleRaw)
							update.Properties.ExportPolicy = exportPolicyRule
						}

						if metadata.ResourceData.HasChange(fmt.Sprintf("%v.data_protection_snapshot_policy", volumeItem)) {
							// Validating that snapshot policies are not being created in a data protection volume
							dataProtectionReplicationRaw := metadata.ResourceData.Get(fmt.Sprintf("%v.data_protection_replication", volumeItem)).([]interface{})
							dataProtectionReplication := expandNetAppVolumeDataProtectionReplication(dataProtectionReplicationRaw)

							if dataProtectionReplication != nil && dataProtectionReplication.Replication != nil && dataProtectionReplication.Replication.EndpointType != nil && strings.ToLower(string(*dataProtectionReplication.Replication.EndpointType)) == "dst" {
								return fmt.Errorf("snapshot policy cannot be enabled on a data protection volume, %s", id)
							}

							dataProtectionSnapshotPolicyRaw := metadata.ResourceData.Get(fmt.Sprintf("%v.data_protection_snapshot_policy", volumeItem)).([]interface{})
							dataProtectionSnapshotPolicy := expandNetAppVolumeDataProtectionSnapshotPolicyPatch(dataProtectionSnapshotPolicyRaw)
							update.Properties.DataProtection = dataProtectionSnapshotPolicy
						}

						if metadata.ResourceData.HasChange(fmt.Sprintf("%v.throughput_in_mibps", volumeItem)) {
							throughputMibps := metadata.ResourceData.Get(fmt.Sprintf("%v.throughput_in_mibps", volumeItem))
							update.Properties.ThroughputMibps = utils.Float(throughputMibps.(float64))
						}

						if metadata.ResourceData.HasChange(fmt.Sprintf("%v.tags", volumeItem)) {
							tagsRaw := metadata.ResourceData.Get(fmt.Sprintf("%v.tags", volumeItem)).(map[string]interface{})
							update.Tags = tags.Expand(tagsRaw)
						}

						if err = volumeClient.UpdateThenPoll(ctx, *id, update); err != nil {
							return fmt.Errorf("updating %s: %+v", id, err)
						}

						// Wait for volume to complete update
						if err := waitForVolumeCreateOrUpdate(ctx, volumeClient, *id); err != nil {
							return err
						}

					}
				}
			}

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
