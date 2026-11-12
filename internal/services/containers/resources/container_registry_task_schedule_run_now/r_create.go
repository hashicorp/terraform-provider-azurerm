// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package container_registry_task_schedule_run_now

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2019-06-01-preview/registries"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2019-06-01-preview/runs"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2019-06-01-preview/tasks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Create() sdk.ResourceFunc {
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

			scheduleResp, err := registryClient.ScheduleRun(ctx, registryId, req)
			if err != nil {
				return fmt.Errorf("scheduling the task: %+v", err)
			}
			if scheduleResp.Model == nil {
				// If the SDK didn't parse the response body, try parsing it on our side.
				if scheduleResp.HttpResponse != nil {
					scheduleRunModel := registries.Run{}
					if err = json.
						NewDecoder(scheduleResp.HttpResponse.Body).
						Decode(&scheduleRunModel); err != nil {
						return fmt.Errorf("can't decode ScheduleRun model, err: %w for taskID %s", err, taskId)
					}

					scheduleResp.Model = &scheduleRunModel
				}

				// If parsing on our side didn't work as well - throw error.
				if scheduleResp.Model == nil {
					return fmt.Errorf("ScheduleRun model was nil (status: %d) for taskID %s", scheduleResp.HttpResponse.StatusCode, taskId)
				}
			}

			runName := pointer.From(scheduleResp.Model.Name)
			if runName == "" {
				return fmt.Errorf("unexpected nil scheduled run name")
			}

			runId := runs.NewRunID(registryId.SubscriptionId, registryId.ResourceGroupName, registryId.RegistryName, runName)

			runsClient := metadata.Client.Containers.ContainerRegistryClient_v2019_06_01_preview.Runs

			timeout, _ := ctx.Deadline()
			stateConf := &pluginsdk.StateChangeConf{
				Pending: []string{string(registries.RunStatusQueued), string(registries.RunStatusStarted), string(registries.RunStatusRunning)},
				Target:  []string{string(registries.RunStatusSucceeded)},
				Refresh: func() (any, string, error) {
					resp, err := runsClient.Get(ctx, runId)
					if err != nil {
						return nil, "", fmt.Errorf("getting the scheduled run: %v", err)
					}

					if resp.Model == nil {
						return nil, "", fmt.Errorf("model was nil for %s", runId)
					}

					return resp, string(*resp.Model.Properties.Status), nil
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
