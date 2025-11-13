// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package durabletask

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/durabletask/2025-11-01/schedulers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/durabletask/2025-11-01/taskhubs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type TaskHubResourceModel struct {
	Name         string `tfschema:"name"`
	SchedulerId  string `tfschema:"scheduler_id"`
	DashboardUrl string `tfschema:"dashboard_url"`
}

type TaskHubResource struct{}

var _ sdk.Resource = TaskHubResource{}

func (r TaskHubResource) ResourceType() string {
	return "azurerm_durable_task_hub"
}

func (r TaskHubResource) ModelObject() interface{} {
	return &TaskHubResourceModel{}
}

func (r TaskHubResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return taskhubs.ValidateTaskHubID
}

func (r TaskHubResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: ValidateTaskHubName,
		},

		"scheduler_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: schedulers.ValidateSchedulerID,
		},
	}
}

func (r TaskHubResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"dashboard_url": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r TaskHubResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DurableTask.TaskHubsClient

			var model TaskHubResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			schedulerId, err := schedulers.ParseSchedulerID(model.SchedulerId)
			if err != nil {
				return fmt.Errorf("parsing scheduler ID: %+v", err)
			}

			id := taskhubs.NewTaskHubID(schedulerId.SubscriptionId, schedulerId.ResourceGroupName, schedulerId.SchedulerName, model.Name)

			metadata.Logger.Infof("Import check for %s", id)
			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			metadata.Logger.Infof("Creating %s", id)

			properties := taskhubs.TaskHub{
				Properties: &taskhubs.TaskHubProperties{},
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, properties); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r TaskHubResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DurableTask.TaskHubsClient

			id, err := taskhubs.ParseTaskHubID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Infof("Reading %s", id)
			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			model := resp.Model
			if model == nil {
				return fmt.Errorf("retrieving %s: model was nil", id)
			}

			schedulerId := schedulers.NewSchedulerID(id.SubscriptionId, id.ResourceGroupName, id.SchedulerName)

			state := TaskHubResourceModel{
				Name:        id.TaskHubName,
				SchedulerId: schedulerId.ID(),
			}

			if props := model.Properties; props != nil {
				state.DashboardUrl = pointer.From(props.DashboardUrl)
			}

			return metadata.Encode(&state)
		},
	}
}

func (r TaskHubResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DurableTask.TaskHubsClient

			id, err := taskhubs.ParseTaskHubID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Infof("Deleting %s", id)

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}
