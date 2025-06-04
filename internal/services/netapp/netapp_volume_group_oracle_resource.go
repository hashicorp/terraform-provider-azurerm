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
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2025-01-01/capacitypools"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2025-01-01/volumegroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2025-01-01/volumes"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	netAppModels "github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/models"
	netAppValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type NetAppVolumeGroupOracleResource struct{}

var _ sdk.Resource = NetAppVolumeGroupOracleResource{}

func (r NetAppVolumeGroupOracleResource) ModelObject() interface{} {
	return &netAppModels.NetAppVolumeGroupOracleModel{}
}

func (r NetAppVolumeGroupOracleResource) ResourceType() string {
	return "azurerm_netapp_volume_group_oracle"
}

func (r NetAppVolumeGroupOracleResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return volumegroups.ValidateVolumeGroupID
}

func (r NetAppVolumeGroupOracleResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: netAppValidate.VolumeGroupName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"account_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: netAppValidate.AccountName,
		},

		"group_description": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"application_identifier": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringLenBetween(1, 3),
		},

		"volume": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MinItems: 2,
			MaxItems: 12,
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
						Type:     pluginsdk.TypeString,
						Optional: true,
						ForceNew: true,
					},

					"zone": commonschema.ZoneSingleOptionalForceNew(),

					"volume_spec_name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringInSlice(netAppValidate.PossibleValuesForVolumeSpecNameOracle(), false),
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
						Type:     pluginsdk.TypeList,
						ForceNew: true,
						Required: true,
						MinItems: 1,
						MaxItems: 1,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringInSlice(netAppValidate.PossibleValuesForProtocolTypeVolumeGroupOracle(), false),
						},
					},

					"security_style": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringInSlice(volumegroups.PossibleValuesForSecurityStyle(), false),
					},

					"storage_quota_in_gb": {
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntBetween(100, 102400),
					},

					"throughput_in_mibps": {
						Type:         pluginsdk.TypeFloat,
						Required:     true,
						ValidateFunc: validation.FloatAtLeast(0.1),
					},

					"export_policy_rule": {
						Type:     pluginsdk.TypeList,
						Required: true,
						MinItems: 1,
						MaxItems: 5,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"rule_index": {
									Type:         pluginsdk.TypeInt,
									Required:     true,
									ValidateFunc: validation.IntBetween(1, 5),
								},

								"allowed_clients": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
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
									Optional: true,
									Default:  false,
								},

								"unix_read_write": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									Default:  true,
								},

								"root_access_enabled": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									Default:  true,
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

					"network_features": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Computed:     true, // O+C - This is Optional/Computed because the service team is changing network features on the backend to upgrade everyone from Basic to Standard and there is a feature that allows customers to change network features from portal but not the API. This could cause drift that forces data loss that we want to avoid
						ValidateFunc: validation.StringInSlice(volumegroups.PossibleValuesForNetworkFeatures(), false),
					},

					"mount_ip_addresses": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"data_protection_snapshot_policy": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Computed: true, // O+C - Adding this because Terraform is not being able to build proper deletion graph, it is trying to delete the snapshot policy before the volume because this is in a deeper level within the schema inside an array of volumes
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

					"encryption_key_source": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						Computed:     true, // O+C - This is computed/optional since there is a feature coming up that will allow customers to change the encryption key source from portal but not the API. This could cause drift if configuration is not updated
						ValidateFunc: validation.StringInSlice(volumes.PossibleValuesForEncryptionKeySource(), false),
					},

					"key_vault_private_endpoint_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						Computed:     true, // O+C - This is computed/optional since there is a feature coming up that will allow customers to change the encryption key source from portal but not the API (tied to encryption_key_source). This could cause drift if configuration is not updated
						ValidateFunc: azure.ValidateResourceID,
					},
				},
			},
		},
	}
}

func (r NetAppVolumeGroupOracleResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r NetAppVolumeGroupOracleResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 90 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.NetApp.VolumeGroupClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model netAppModels.NetAppVolumeGroupOracleModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := volumegroups.NewVolumeGroupID(subscriptionId, model.ResourceGroupName, model.AccountName, model.Name)

			metadata.Logger.Infof("Import check for %s", id)
			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}

			if existing.Model != nil && existing.Model.Id != nil && *existing.Model.Id != "" {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			volumeList, err := expandNetAppVolumeGroupOracleVolumes(model.Volumes)
			if err != nil {
				return err
			}

			// Performing some basic validations that are not possible in the schema
			if errorList := netAppValidate.ValidateNetAppVolumeGroupOracleVolumes(volumeList); len(errorList) > 0 {
				return fmt.Errorf("one or more issues found while performing deeper validations for %s:\n%+v", id, errorList)
			}

			parameters := volumegroups.VolumeGroupDetails{
				Location: utils.String(location.Normalize(model.Location)),
				Properties: &volumegroups.VolumeGroupProperties{
					GroupMetaData: &volumegroups.VolumeGroupMetaData{
						GroupDescription:      utils.String(model.GroupDescription),
						ApplicationType:       pointer.To(volumegroups.ApplicationTypeORACLE),
						ApplicationIdentifier: utils.String(model.ApplicationIdentifier),
					},
					Volumes: volumeList,
				},
			}

			if err = client.CreateThenPoll(ctx, id, parameters); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			// Waiting for volume group be completely provisioned
			if err := waitForVolumeGroupCreateOrUpdate(ctx, client, id); err != nil {
				return fmt.Errorf("waiting creation %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (r NetAppVolumeGroupOracleResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 120 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			volumeClient := metadata.Client.NetApp.VolumeClient

			id, err := volumegroups.ParseVolumeGroupID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Infof("Decoding state for %s", id)
			var state netAppModels.NetAppVolumeGroupOracleModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			metadata.Logger.Infof("Updating %s", id)

			if metadata.ResourceData.HasChange("volume") {
				for i := 0; i < metadata.ResourceData.Get("volume.#").(int); i++ {
					// Checking if individual volume has a change
					volumeItem := fmt.Sprintf("volume.%v", i)

					capacityPoolId, err := capacitypools.ParseCapacityPoolID(metadata.ResourceData.Get(fmt.Sprintf("%v.capacity_pool_id", volumeItem)).(string))
					if err != nil {
						return err
					}

					if metadata.ResourceData.HasChange(volumeItem) {
						volumeId := volumes.NewVolumeID(id.SubscriptionId,
							id.ResourceGroupName,
							id.NetAppAccountName,
							capacityPoolId.CapacityPoolName,
							metadata.ResourceData.Get(fmt.Sprintf("%v.name", volumeItem)).(string))

						update := volumes.VolumePatch{
							Properties: &volumes.VolumePatchProperties{},
						}

						if metadata.ResourceData.HasChange(fmt.Sprintf("%v.storage_quota_in_gb", volumeItem)) {
							storageQuotaInBytes := int64(metadata.ResourceData.Get(fmt.Sprintf("%v.storage_quota_in_gb", volumeItem)).(int) * 1073741824)
							update.Properties.UsageThreshold = utils.Int64(storageQuotaInBytes)
						}

						if metadata.ResourceData.HasChange(fmt.Sprintf("%v.export_policy_rule", volumeItem)) {
							exportPolicyRuleRaw := metadata.ResourceData.Get(fmt.Sprintf("%v.export_policy_rule", volumeItem)).([]interface{})

							// Validating export policy rules
							volumeProtocolRaw := (metadata.ResourceData.Get(fmt.Sprintf("%v.protocols", volumeItem)).([]interface{}))[0]
							volumeProtocol := volumeProtocolRaw.(string)

							errors := make([]error, 0)
							for _, ruleRaw := range exportPolicyRuleRaw {
								if ruleRaw != nil {
									rule := volumegroups.ExportPolicyRule{}

									v := ruleRaw.(map[string]interface{})
									rule.Nfsv3 = utils.Bool(v["nfsv3_enabled"].(bool))
									rule.Nfsv41 = utils.Bool(v["nfsv41_enabled"].(bool))

									errors = append(errors, netAppValidate.ValidateNetAppVolumeGroupExportPolicyRule(rule, volumeProtocol)...)
								}
							}

							if len(errors) > 0 {
								return fmt.Errorf("one or more issues found while performing export policies validations for %s:\n%+v", id, errors)
							}

							exportPolicyRule := expandNetAppVolumeGroupVolumeExportPolicyRulePatch(exportPolicyRuleRaw)
							update.Properties.ExportPolicy = exportPolicyRule
						}

						if metadata.ResourceData.HasChange(fmt.Sprintf("%v.data_protection_snapshot_policy", volumeItem)) {
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

						if err = volumeClient.UpdateThenPoll(ctx, volumeId, update); err != nil {
							return fmt.Errorf("updating %s: %+v", volumeId, err)
						}

						// Waiting for volume to fully complete an update
						if err := waitForVolumeCreateOrUpdate(ctx, volumeClient, volumeId); err != nil {
							return fmt.Errorf("waiting update %s: %+v", volumeId, err)
						}
					}
				}
			}

			return nil
		},
	}
}

func (r NetAppVolumeGroupOracleResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.NetApp.VolumeGroupClient

			id, err := volumegroups.ParseVolumeGroupID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Infof("Decoding state for %s", id)
			var state netAppModels.NetAppVolumeGroupOracleModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.Get(ctx, pointer.From(id))
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %v", id, err)
			}

			model := netAppModels.NetAppVolumeGroupOracleModel{
				Name:              id.VolumeGroupName,
				AccountName:       id.NetAppAccountName,
				Location:          location.NormalizeNilable(existing.Model.Location),
				ResourceGroupName: id.ResourceGroupName,
			}

			if props := existing.Model.Properties; props != nil {
				model.GroupDescription = pointer.From(props.GroupMetaData.GroupDescription)
				model.ApplicationIdentifier = pointer.From(props.GroupMetaData.ApplicationIdentifier)

				volumes, err := flattenNetAppVolumeGroupOracleVolumes(ctx, props.Volumes, metadata)
				if err != nil {
					return fmt.Errorf("setting `volume`: %+v", err)
				}

				model.Volumes = volumes
			}

			metadata.SetID(id)

			return metadata.Encode(&model)
		},
	}
}

func (r NetAppVolumeGroupOracleResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 120 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.NetApp.VolumeGroupClient

			id, err := volumegroups.ParseVolumeGroupID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.Get(ctx, pointer.From(id))
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %v", id, err)
			}

			// Removing volumes before deleting volume group
			if props := existing.Model.Properties; props != nil {
				if volumeList := props.Volumes; volumeList != nil {
					for _, volume := range *volumeList {
						if err := deleteVolume(ctx, metadata, pointer.From(volume.Id)); err != nil {
							return fmt.Errorf("deleting `volume`: %+v", err)
						}
					}
				}
			}

			if err = client.DeleteThenPoll(ctx, pointer.From(id)); err != nil {
				return fmt.Errorf("deleting %s: %+v", pointer.From(id), err)
			}

			// Waiting for volume group be completely deleted
			if err := waitForVolumeGroupDelete(ctx, client, pointer.From(id)); err != nil {
				return fmt.Errorf("waiting delete %s: %+v", pointer.From(id), err)
			}

			return nil
		},
	}
}
