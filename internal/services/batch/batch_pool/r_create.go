// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package batch_pool

import (
	"errors"
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-sdk/resource-manager/batch/2024-07-01/pool"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceBatchCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Batch.PoolClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := pool.NewPoolID(subscriptionId, d.Get("resource_group_name").(string), d.Get("account_name").(string), d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_batch_pool", id.ID())
		}
	}

	parameters := pool.Pool{
		Properties: &pool.PoolProperties{
			VMSize:                 pointer.To(d.Get("vm_size").(string)),
			DisplayName:            pointer.To(d.Get("display_name").(string)),
			InterNodeCommunication: pointer.ToEnum[pool.InterNodeCommunicationState](d.Get("inter_node_communication").(string)),
			TaskSlotsPerNode:       pointer.To(int64(d.Get("max_tasks_per_node").(int))),
		},
	}

	userAccounts, err := ExpandBatchPoolUserAccounts(d)
	if err != nil {
		log.Printf(`[DEBUG] expanding "user_accounts": %v`, err)
	}
	parameters.Properties.UserAccounts = userAccounts

	taskSchedulingPolicy, err := ExpandBatchPoolTaskSchedulingPolicy(d)
	if err != nil {
		log.Printf(`[DEBUG] expanding "task_scheduling_policy": %v`, err)
	}
	parameters.Properties.TaskSchedulingPolicy = taskSchedulingPolicy

	identityResult, err := identity.ExpandUserAssignedMap(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf(`expanding "identity": %v`, err)
	}
	parameters.Identity = identityResult

	scaleSettings, err := expandBatchPoolScaleSettings(d)
	if err != nil {
		return fmt.Errorf("expanding scale settings: %+v", err)
	}

	parameters.Properties.ScaleSettings = scaleSettings

	if startTaskValue, startTaskOk := d.GetOk("start_task"); startTaskOk {
		startTaskList := startTaskValue.([]interface{})
		startTask, startTaskErr := ExpandBatchPoolStartTask(startTaskList)

		if startTaskErr != nil {
			return fmt.Errorf("creating %s: %+v", id, startTaskErr)
		}

		// start task should have a user identity defined
		if userIdentityError := validateUserIdentity(startTask.UserIdentity); userIdentityError != nil {
			return fmt.Errorf("creating %s: %+v", id, userIdentityError)
		}

		parameters.Properties.StartTask = startTask
	}

	if vmDeploymentConfiguration, deploymentErr := expandBatchPoolVirtualMachineConfig(d); deploymentErr == nil {
		parameters.Properties.DeploymentConfiguration = &pool.DeploymentConfiguration{
			VirtualMachineConfiguration: vmDeploymentConfiguration,
		}
	} else {
		return deploymentErr
	}

	if err := validateBatchPoolCrossFieldRules(parameters.Properties); err != nil {
		return err
	}

	metaDataRaw := d.Get("metadata").(map[string]interface{})
	parameters.Properties.Metadata = ExpandBatchMetaData(metaDataRaw)

	mountConfiguration, err := ExpandBatchPoolMountConfigurations(d)
	if err != nil {
		log.Printf(`[DEBUG] expanding "mount": %v`, err)
	}
	parameters.Properties.MountConfiguration = mountConfiguration

	networkConfiguration := d.Get("network_configuration").([]interface{})
	parameters.Properties.NetworkConfiguration, err = ExpandBatchPoolNetworkConfiguration(networkConfiguration)
	if err != nil {
		return fmt.Errorf("expanding `network_configuration`: %+v", err)
	}

	if v, ok := d.GetOk("target_node_communication_mode"); ok {
		parameters.Properties.TargetNodeCommunicationMode = pointer.ToEnum[pool.NodeCommunicationMode](v.(string))
	}

	if _, err = client.Create(ctx, id, parameters, pool.CreateOperationOptions{}); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	read, err := client.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())
	if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
		return err
	}

	// if the pool is not Steady after the create operation, wait for it to be Steady
	if model := read.Model; model != nil {
		if props := model.Properties; props != nil && props.AllocationState != nil && *props.AllocationState != pool.AllocationStateSteady {
			if err = waitForBatchPoolPendingResizeOperation(ctx, client, id); err != nil {
				return fmt.Errorf("waiting for %s", id)
			}
		}
	}

	return resourceBatchPoolRead(d, meta)
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
