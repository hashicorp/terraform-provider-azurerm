package containers

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/containerregistry/mgmt/2021-08-01-preview/containerregistry" // nolint: staticcheck
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
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
			ValidateFunc: validate.ContainerRegistryTaskID,
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
			taskClient := metadata.Client.Containers.TasksClient

			var model ContainerRegistryTaskScheduleModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			taskId, err := parse.ContainerRegistryTaskID(model.TaskId)
			if err != nil {
				return fmt.Errorf("parsing container registry task ID: %v", err)
			}

			resp, err := taskClient.Get(ctx, taskId.ResourceGroup, taskId.RegistryName, taskId.TaskName)
			if err != nil {
				return fmt.Errorf("retrieving %q: %+v", taskId, err)
			}
			if resp.TaskProperties == nil {
				return fmt.Errorf("unexpected nil `taskProperties` of %q", taskId)
			}
			if resp.TaskProperties.Step == nil {
				return fmt.Errorf("unexpected nil `taskProperties.step` of %q", taskId)
			}

			req := containerregistry.TaskRunRequest{
				TaskID: utils.String(taskId.ID()),
			}
			switch resp.TaskProperties.Step.(type) {
			case containerregistry.DockerBuildStep:
				req.Type = containerregistry.TypeDockerBuildRequest
			case containerregistry.FileTaskStep:
				req.Type = containerregistry.TypeFileTaskRunRequest
			case containerregistry.EncodedTaskStep:
				req.Type = containerregistry.TypeEncodedTaskRunRequest
			default:
				return fmt.Errorf("unexpected container registry task step type: %T", resp.TaskProperties.Step)
			}

			registryClient := metadata.Client.Containers.RegistriesClient
			future, err := registryClient.ScheduleRun(ctx, taskId.ResourceGroup, taskId.RegistryName, req)
			if err != nil {
				return fmt.Errorf("scheduling the task: %v", err)
			}
			if err := future.WaitForCompletionRef(ctx, registryClient.Client); err != nil {
				return fmt.Errorf("waiting for schedule: %v", err)
			}

			run, err := future.Result(*registryClient)
			if err != nil {
				return fmt.Errorf("getting the scheduled run: %v", err)
			}

			if run.Name == nil {
				return fmt.Errorf("unexpected nil scheduled run name")
			}

			runsClient := metadata.Client.Containers.RunsClient

			timeout, _ := ctx.Deadline()
			stateConf := &pluginsdk.StateChangeConf{
				Pending: []string{string(containerregistry.RunStatusQueued), string(containerregistry.RunStatusStarted), string(containerregistry.RunStatusRunning)},
				Target:  []string{string(containerregistry.RunStatusSucceeded)},
				Refresh: func() (interface{}, string, error) {
					resp, err := runsClient.Get(ctx, taskId.ResourceGroup, taskId.RegistryName, *run.Name)
					if err != nil {
						return nil, "", fmt.Errorf("getting the scheduled run: %v", err)
					}

					if resp.RunProperties == nil {
						return nil, "", fmt.Errorf("unexpected nil properties of the scheduled run")
					}

					return run, string(resp.RunProperties.Status), nil
				},
				ContinuousTargetOccurence: 1,
				PollInterval:              5 * time.Second,
				Timeout:                   time.Until(timeout),
			}
			if _, err := stateConf.WaitForStateContext(ctx); err != nil {
				return fmt.Errorf("waiting for scheduled task to finish: %+v", err)
			}

			metadata.SetID(parse.NewContainerRegistryTaskScheduleID(taskId.SubscriptionId, taskId.ResourceGroup, taskId.RegistryName, taskId.TaskName, "schedule"))
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
			model := ContainerRegistryTaskScheduleModel{TaskId: parse.NewContainerRegistryTaskID(id.SubscriptionId, id.ResourceGroup, id.RegistryName, id.TaskName).ID()}
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
