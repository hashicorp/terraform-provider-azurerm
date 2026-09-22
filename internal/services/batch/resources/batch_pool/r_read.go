// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package batch_pool

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-sdk/resource-manager/batch/2024-07-01/pool"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceBatchPoolRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Batch.PoolClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := pool.ParsePoolID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.PoolName)
	d.Set("account_name", id.BatchAccountName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		identityResult, err := identity.FlattenUserAssignedMap(model.Identity)
		if err != nil {
			return fmt.Errorf("flattening `identity`: %+v", err)
		}
		if err := d.Set("identity", identityResult); err != nil {
			return fmt.Errorf("setting `identity`: %+v", err)
		}

		if props := model.Properties; props != nil {
			d.Set("display_name", props.DisplayName)
			d.Set("vm_size", props.VMSize)
			d.Set("inter_node_communication", string(pointer.From(props.InterNodeCommunication)))

			if scaleSettings := props.ScaleSettings; scaleSettings != nil {
				if err := d.Set("auto_scale", flattenBatchPoolAutoScaleSettings(scaleSettings.AutoScale)); err != nil {
					return fmt.Errorf("flattening `auto_scale`: %+v", err)
				}
				if err := d.Set("fixed_scale", flattenBatchPoolFixedScaleSettings(d, scaleSettings.FixedScale)); err != nil {
					return fmt.Errorf("flattening `fixed_scale `: %+v", err)
				}
			}

			if props.TaskSchedulingPolicy != nil && props.TaskSchedulingPolicy.NodeFillType != "" {
				taskSchedulingPolicy := make([]interface{}, 0)
				nodeFillType := make(map[string]interface{})
				nodeFillType["node_fill_type"] = string(props.TaskSchedulingPolicy.NodeFillType)
				taskSchedulingPolicy = append(taskSchedulingPolicy, nodeFillType)
				d.Set("task_scheduling_policy", taskSchedulingPolicy)
			}

			if props.UserAccounts != nil {
				userAccounts := make([]interface{}, 0)
				for _, userAccount := range *props.UserAccounts {
					userAccounts = append(userAccounts, flattenBatchPoolUserAccount(d, &userAccount))
				}
				d.Set("user_accounts", userAccounts)
			}

			d.Set("max_tasks_per_node", props.TaskSlotsPerNode)

			if props.DeploymentConfiguration != nil {
				if props.DeploymentConfiguration.VirtualMachineConfiguration != nil {
					config := props.DeploymentConfiguration.VirtualMachineConfiguration
					if config.ContainerConfiguration != nil {
						d.Set("container_configuration", flattenBatchPoolContainerConfiguration(d, config.ContainerConfiguration))
					}
					if config.DataDisks != nil {
						dataDisks := make([]interface{}, 0)
						for _, item := range *config.DataDisks {
							dataDisk := make(map[string]interface{})
							dataDisk["lun"] = item.Lun
							dataDisk["disk_size_gb"] = item.DiskSizeGB

							caching := ""
							if item.Caching != nil {
								caching = string(*item.Caching)
							}
							dataDisk["caching"] = caching

							storageAccountType := ""
							if item.StorageAccountType != nil {
								storageAccountType = string(*item.StorageAccountType)
							}
							dataDisk["storage_account_type"] = storageAccountType

							dataDisks = append(dataDisks, dataDisk)
						}
						d.Set("data_disks", dataDisks)
					}
					if config.DiskEncryptionConfiguration != nil {
						diskEncryptionConfiguration := make([]interface{}, 0)
						if config.DiskEncryptionConfiguration.Targets != nil {
							for _, item := range *config.DiskEncryptionConfiguration.Targets {
								target := make(map[string]interface{})
								target["disk_encryption_target"] = string(item)
								diskEncryptionConfiguration = append(diskEncryptionConfiguration, target)
							}
						}
						d.Set("disk_encryption", diskEncryptionConfiguration)
					}
					if config.Extensions != nil {
						extensions := make([]interface{}, 0)
						n := len(*config.Extensions)
						for _, item := range *config.Extensions {
							extension := make(map[string]interface{})
							extension["name"] = item.Name
							extension["publisher"] = item.Publisher
							extension["type"] = item.Type
							if item.TypeHandlerVersion != nil {
								extension["type_handler_version"] = *item.TypeHandlerVersion
							}
							if item.AutoUpgradeMinorVersion != nil {
								extension["auto_upgrade_minor_version"] = *item.AutoUpgradeMinorVersion
							}
							if item.EnableAutomaticUpgrade != nil {
								extension["automatic_upgrade_enabled"] = *item.EnableAutomaticUpgrade
							}
							if item.Settings != nil {
								settingValue, err := json.Marshal((*item.Settings).(map[string]interface{}))
								if err != nil {
									return fmt.Errorf("flattening `settings_json`: %+v", err)
								}
								extension["settings_json"] = string(settingValue)
							}

							for i := range n {
								if v, ok := d.GetOk(fmt.Sprintf("extensions.%d.name", i)); ok && v == item.Name {
									extension["protected_settings"] = d.Get(fmt.Sprintf("extensions.%d.protected_settings", i))
									break
								}
							}

							if item.ProvisionAfterExtensions != nil {
								extension["provision_after_extensions"] = *item.ProvisionAfterExtensions
							}
							extensions = append(extensions, extension)
						}
						d.Set("extensions", extensions)
					}

					d.Set("storage_image_reference", flattenBatchPoolImageReference(&config.ImageReference))
					d.Set("license_type", config.LicenseType)
					d.Set("node_agent_sku_id", config.NodeAgentSkuId)

					if config.NodePlacementConfiguration != nil {
						nodePlacementConfiguration := make([]interface{}, 0)
						nodePlacementConfig := make(map[string]interface{})
						nodePlacementConfig["policy"] = string(*config.NodePlacementConfiguration.Policy)
						nodePlacementConfiguration = append(nodePlacementConfiguration, nodePlacementConfig)
						d.Set("node_placement", nodePlacementConfiguration)
					}

					osDiskPlacement := ""
					if config.OsDisk != nil && config.OsDisk.EphemeralOSDiskSettings != nil && config.OsDisk.EphemeralOSDiskSettings.Placement != nil {
						osDiskPlacement = string(*config.OsDisk.EphemeralOSDiskSettings.Placement)
					}
					d.Set("os_disk_placement", osDiskPlacement)

					if config.SecurityProfile != nil {
						d.Set("security_profile", flattenBatchPoolSecurityProfile(config.SecurityProfile))
					}

					if config.WindowsConfiguration != nil {
						windowsConfig := []interface{}{
							map[string]interface{}{
								"enable_automatic_updates": *config.WindowsConfiguration.EnableAutomaticUpdates,
							},
						}
						d.Set("windows", windowsConfig)
					}
				}
			}

			d.Set("start_task", flattenBatchPoolStartTask(d, props.StartTask))
			d.Set("metadata", FlattenBatchMetaData(props.Metadata))

			if props.MountConfiguration != nil {
				mountConfigs := make([]interface{}, 0)
				for _, mountConfig := range *props.MountConfiguration {
					mountConfigs = append(mountConfigs, flattenBatchPoolMountConfig(d, &mountConfig))
				}
				d.Set("mount", mountConfigs)
			}

			targetNodeCommunicationMode := ""
			if props.TargetNodeCommunicationMode != nil {
				targetNodeCommunicationMode = string(*props.TargetNodeCommunicationMode)
			}
			d.Set("target_node_communication_mode", targetNodeCommunicationMode)

			if err := d.Set("network_configuration", flattenBatchPoolNetworkConfiguration(props.NetworkConfiguration)); err != nil {
				return fmt.Errorf("setting `network_configuration`: %v", err)
			}
		}
	}

	return pluginsdk.SetResourceIdentityData(d, id)
}

func waitForBatchPoolPendingResizeOperation(ctx context.Context, client *pool.PoolClient, id pool.PoolId) error {
	// waiting for the pool to be in steady state
	log.Printf("[INFO] waiting for the pending resize operation on this pool to be stopped...")
	isSteady := false
	for !isSteady {
		resp, err := client.Get(ctx, id)
		if err != nil {
			return fmt.Errorf("retrieving %s: %+v", id, err)
		}

		if resp.Model != nil && resp.Model.Properties != nil && resp.Model.Properties.AllocationState != nil {
			isSteady = *resp.Model.Properties.AllocationState == pool.AllocationStateSteady
			if isSteady {
				break
			}
		}
		time.Sleep(time.Second * 30)
		log.Printf("[INFO] waiting for the pending resize operation on this pool to be stopped... New try in 30 seconds...")
	}
	return nil
}

// validateUserIdentity validates that the user identity for a start task has been well specified
// it should have a auto_user block or a user_name defined, but not both at the same time.
func validateUserIdentity(userIdentity *pool.UserIdentity) error {
	if userIdentity == nil {
		return errors.New("user_identity block needs to be specified")
	}

	if userIdentity.AutoUser == nil && userIdentity.UserName == nil {
		return errors.New("auto_user or user_name needs to be specified in the user_identity block")
	}

	if userIdentity.AutoUser != nil && userIdentity.UserName != nil && *userIdentity.UserName != "" {
		return errors.New("auto_user and user_name cannot be specified in the user_identity at the same time")
	}

	return nil
}

func validateBatchPoolCrossFieldRules(pool *pool.PoolProperties) error {
	// Perform validation across multiple fields as per https://docs.microsoft.com/rest/api/batchmanagement/pool/create#resourcefile

	if pool.StartTask != nil {
		startTask := *pool.StartTask
		if startTask.ResourceFiles != nil {
			for _, referenceFile := range *startTask.ResourceFiles {
				// Must specify exactly one of AutoStorageContainerName, StorageContainerURL or HttpUrl
				sourceCount := 0
				if referenceFile.AutoStorageContainerName != nil {
					sourceCount++
				}
				if referenceFile.StorageContainerURL != nil {
					sourceCount++
				}
				if referenceFile.HTTPURL != nil {
					sourceCount++
				}
				if sourceCount != 1 {
					return fmt.Errorf("exactly one of auto_storage_container_name, storage_container_url and http_url must be specified")
				}

				if referenceFile.BlobPrefix != nil {
					if referenceFile.AutoStorageContainerName == nil && referenceFile.StorageContainerURL == nil {
						return fmt.Errorf("auto_storage_container_name or storage_container_url must be specified when using blob_prefix")
					}
				}

				if referenceFile.HTTPURL != nil {
					if referenceFile.FilePath == nil {
						return fmt.Errorf("file_path must be specified when using http_url")
					}
				}
			}
		}
	}

	return nil
}

func containerRegistry() map[string]*schema.Schema {
	return map[string]*pluginsdk.Schema{
		"registry_server": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"user_assigned_identity_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateUserAssignedIdentityID,
			Description:  "The User Assigned Identity to use for Container Registry access.",
		},
		"user_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"password": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}
