// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package batch_pool

import (
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

func expandBatchPoolScaleSettings(d *pluginsdk.ResourceData) (*pool.ScaleSettings, error) {
	scaleSettings := &pool.ScaleSettings{}

	autoScaleValue, autoScaleOk := d.GetOk("auto_scale")
	fixedScaleValue, fixedScaleOk := d.GetOk("fixed_scale")

	if !autoScaleOk && !fixedScaleOk {
		return nil, fmt.Errorf("auto_scale block or fixed_scale block need to be specified")
	}

	if autoScaleOk && fixedScaleOk {
		return nil, fmt.Errorf("auto_scale and fixed_scale blocks cannot be specified at the same time")
	}

	if autoScaleOk {
		autoScale := autoScaleValue.([]interface{})
		if len(autoScale) == 0 {
			return nil, fmt.Errorf("when scale mode is Auto, auto_scale block is required")
		}

		autoScaleSettings := autoScale[0].(map[string]interface{})

		autoScaleFormula := autoScaleSettings["formula"].(string)

		scaleSettings.AutoScale = &pool.AutoScaleSettings{
			EvaluationInterval: pointer.To(autoScaleSettings["evaluation_interval"].(string)),
			Formula:            autoScaleFormula,
		}
	} else if fixedScaleOk {
		fixedScale := fixedScaleValue.([]interface{})
		if len(fixedScale) == 0 {
			return nil, fmt.Errorf("when scale mode is Fixed, fixed_scale block is required")
		}

		fixedScaleSettings := fixedScale[0].(map[string]interface{})
		targetDedicatedNodes := int32(fixedScaleSettings["target_dedicated_nodes"].(int))
		targetLowPriorityNodes := int32(fixedScaleSettings["target_low_priority_nodes"].(int))

		scaleSettings.FixedScale = &pool.FixedScaleSettings{
			NodeDeallocationOption: pointer.ToEnum[pool.ComputeNodeDeallocationOption](fixedScaleSettings["node_deallocation_method"].(string)),
			ResizeTimeout:          pointer.To(fixedScaleSettings["resize_timeout"].(string)),
			TargetDedicatedNodes:   pointer.To(int64(targetDedicatedNodes)),
			TargetLowPriorityNodes: pointer.To(int64(targetLowPriorityNodes)),
		}
	}

	return scaleSettings, nil
}
