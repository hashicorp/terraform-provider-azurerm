// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package containers

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2019-06-01-preview/registries"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2019-06-01-preview/runs"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2019-06-01-preview/tasks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ContainerRegistryTaskScheduleResource struct{}

var _ sdk.Resource = ContainerRegistryTaskScheduleResource{}

type ContainerRegistryTaskScheduleModel struct {
	TaskId string `tfschema:"container_registry_task_id"`
}

func (r ContainerRegistryTaskScheduleResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"container_registry_task_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: tasks.ValidateTaskID,
		},
	}
}

func (r ContainerRegistryTaskScheduleResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r ContainerRegistryTaskScheduleResource) ResourceType() string {
	return "azurerm_container_registry_task_schedule_run_now"
}

func (r ContainerRegistryTaskScheduleResource) ModelObject() interface{} {
	return &ContainerRegistryTaskScheduleModel{}
}

func (r ContainerRegistryTaskScheduleResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.ContainerRegistryTaskScheduleID
}

func (r ContainerRegistryTaskScheduleResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			taskClient := metadata.Client.Containers.ContainerRegistryClient_v2019_06_01_preview.Tasks

			var model ContainerRegistryTaskScheduleModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			taskId, err := tasks.ParseTaskID(model.TaskId)
			if err != nil {
				return err
			}

			resp, err := taskClient.Get(ctx, *taskId)
			if err != nil {
				return fmt.Errorf("retrieving %q: %+v", taskId, err)
			}
			if resp.Model == nil {
				return fmt.Errorf("model was nil for %s", taskId)
			}
			if resp.Model.Properties.Step == nil {
				return fmt.Errorf("properties was nil for %s", taskId)
			}

			req := registries.TaskRunRequest{
				TaskId: taskId.ID(),
			}

			registryId := registries.NewRegistryID(taskId.SubscriptionId, taskId.ResourceGroupName, taskId.RegistryName)
			registryClient := metadata.Client.Containers.ContainerRegistryClient_v2019_06_01_preview.Registries

			_, err = registryClient.ScheduleRun(ctx, registryId, req)
			if err != nil {
				return fmt.Errorf("scheduling the task: %+v", err)

			}

			runsClient := metadata.Client.Containers.ContainerRegistryClient_v2019_06_01_preview.Runs
			run, err := runsClient.List(ctx, runs.RegistryId(registryId), runs.ListOperationOptions{})
			if err != nil {
				return fmt.Errorf("retrieving runs for %s: %+v", taskId, err)
			}

			if run.Model == nil {
				return fmt.Errorf("model was nil for %s", registryId)
			}

			runName := ""
			for _, v := range *run.Model {
				if *v.Properties.Task == taskId.TaskName {
					runName = *v.Name
				}
			}

			if runName == "" {
				return fmt.Errorf("unexpected nil scheduled run name")
			}

			runId := runs.NewRunID(registryId.SubscriptionId, registryId.ResourceGroupName, registryId.RegistryName, runName)

			timeout, _ := ctx.Deadline()
			stateConf := &pluginsdk.StateChangeConf{
				Pending: []string{string(registries.RunStatusQueued), string(registries.RunStatusStarted), string(registries.RunStatusRunning)},
				Target:  []string{string(registries.RunStatusSucceeded)},
				Refresh: func() (interface{}, string, error) {
					resp, err := runsClient.Get(ctx, runId)
					if err != nil {
						return nil, "", fmt.Errorf("getting the scheduled run: %v", err)
					}

					if resp.Model == nil {
						return nil, "", fmt.Errorf("model was nil for %s", runId)
					}

					return run, string(*resp.Model.Properties.Status), nil
				},
				ContinuousTargetOccurence: 1,
				PollInterval:              5 * time.Second,
				Timeout:                   time.Until(timeout),
			}
			if _, err := stateConf.WaitForStateContext(ctx); err != nil {
				return fmt.Errorf("waiting for scheduled task to finish: %+v", err)
			}

			metadata.SetID(parse.NewContainerRegistryTaskScheduleID(taskId.SubscriptionId, taskId.ResourceGroupName, taskId.RegistryName, taskId.TaskName, "schedule"))
			return nil
		},
	}
}

func (r ContainerRegistryTaskScheduleResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := parse.ContainerRegistryTaskScheduleID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}
			model := ContainerRegistryTaskScheduleModel{TaskId: tasks.NewTaskID(id.SubscriptionId, id.ResourceGroup, id.RegistryName, id.TaskName).ID()}
			return metadata.Encode(&model)
		},
	}
}

func (r ContainerRegistryTaskScheduleResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			return nil
		},
	}
}
