// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package batch_pool

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-sdk/resource-manager/batch/2024-07-01/pool"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceBatchUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Batch.PoolClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := pool.ParsePoolID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil && props.AllocationState != nil && *props.AllocationState != pool.AllocationStateSteady {
			log.Printf("[INFO] there is a pending resize operation on this pool...")
			stopPendingResizeOperation := d.Get("stop_pending_resize_operation").(bool)
			if !stopPendingResizeOperation {
				return fmt.Errorf("updating %s because of pending resize operation. Set flag `stop_pending_resize_operation` to true to force update", *id)
			}

			log.Printf("[INFO] stopping the pending resize operation on this pool...")
			if _, err = client.StopResize(ctx, *id); err != nil {
				return fmt.Errorf("stopping resize operation for %s: %+v", *id, err)
			}

			// waiting for the pool to be in steady state
			if err = waitForBatchPoolPendingResizeOperation(ctx, client, *id); err != nil {
				return fmt.Errorf("waiting for %s", *id)
			}
		}
	}

	parameters := pool.Pool{
		Properties: &pool.PoolProperties{},
	}

	identity, err := identity.ExpandUserAssignedMap(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf(`expanding "identity": %v`, err)
	}
	parameters.Identity = identity

	scaleSettings, err := expandBatchPoolScaleSettings(d)
	if err != nil {
		return fmt.Errorf("expanding scale settings: %+v", err)
	}

	parameters.Properties.ScaleSettings = scaleSettings

	taskSchedulingPolicy, err := ExpandBatchPoolTaskSchedulingPolicy(d)
	if err != nil {
		log.Printf(`[DEBUG] expanding "task_scheduling_policy": %v`, err)
	}
	parameters.Properties.TaskSchedulingPolicy = taskSchedulingPolicy

	userAccounts, err := ExpandBatchPoolUserAccounts(d)
	if err != nil {
		log.Printf(`[DEBUG] expanding "user_accounts": %v`, err)
	}
	parameters.Properties.UserAccounts = userAccounts

	if startTaskValue, startTaskOk := d.GetOk("start_task"); startTaskOk {
		startTaskList := startTaskValue.([]interface{})
		startTask, startTaskErr := ExpandBatchPoolStartTask(startTaskList)

		if startTaskErr != nil {
			return fmt.Errorf("updating %s: %+v", *id, startTaskErr)
		}

		// start task should have a user identity defined
		if userIdentityError := validateUserIdentity(startTask.UserIdentity); userIdentityError != nil {
			return fmt.Errorf("creating %s: %+v", *id, userIdentityError)
		}

		parameters.Properties.StartTask = startTask
	}
	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			// when updating `data_disks`, it has to include additional properties such as `NodeAgentSkuId`, `ImageReference` and `OsDisk`, otherwise API request will fail.
			parameters.Properties.DeploymentConfiguration = props.DeploymentConfiguration
			if d.HasChange("data_disks") {
				parameters.Properties.DeploymentConfiguration.VirtualMachineConfiguration.DataDisks = expandBatchPoolDataDisks(d.Get("data_disks").([]interface{}))
			}
		}
	}

	if err := validateBatchPoolCrossFieldRules(parameters.Properties); err != nil {
		return err
	}

	if d.HasChange("metadata") {
		metaDataRaw := d.Get("metadata").(map[string]interface{})

		parameters.Properties.Metadata = ExpandBatchMetaData(metaDataRaw)
	}

	mountConfiguration, err := ExpandBatchPoolMountConfigurations(d)
	if err != nil {
		log.Printf(`[DEBUG] expanding "mount": %v`, err)
	}
	parameters.Properties.MountConfiguration = mountConfiguration

	if d.HasChange("target_node_communication_mode") {
		parameters.Properties.TargetNodeCommunicationMode = pointer.ToEnum[pool.NodeCommunicationMode](d.Get("target_node_communication_mode").(string))
	}

	result, err := client.Update(ctx, *id, parameters, pool.UpdateOperationOptions{})
	if err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	// if the pool is not Steady after the update, wait for it to be Steady
	if model := result.Model; model != nil {
		if props := model.Properties; props != nil && props.AllocationState != nil && *props.AllocationState != pool.AllocationStateSteady {
			if err := waitForBatchPoolPendingResizeOperation(ctx, client, *id); err != nil {
				return fmt.Errorf("waiting for %s", *id)
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
