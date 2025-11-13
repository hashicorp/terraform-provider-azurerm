// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package durabletask

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/durabletask/2025-11-01/schedulers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type SchedulerResourceModel struct {
	Name              string            `tfschema:"name"`
	ResourceGroupName string            `tfschema:"resource_group_name"`
	Location          string            `tfschema:"location"`
	SkuName           string            `tfschema:"sku_name"`
	IpAllowList       []string          `tfschema:"ip_allow_list"`
	Capacity          int64             `tfschema:"capacity"`
	Tags              map[string]string `tfschema:"tags"`
	Endpoint          string            `tfschema:"endpoint"`
	RedundancyState   string            `tfschema:"redundancy_state"`
}

type SchedulerResource struct{}

var (
	_ sdk.Resource           = SchedulerResource{}
	_ sdk.ResourceWithUpdate = SchedulerResource{}
)

func (r SchedulerResource) ResourceType() string {
	return "azurerm_durable_task_scheduler"
}

func (r SchedulerResource) ModelObject() interface{} {
	return &SchedulerResourceModel{}
}

func (r SchedulerResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return schedulers.ValidateSchedulerID
}

func (r SchedulerResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: ValidateSchedulerName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"sku_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(schedulers.SkuNameConsumption),
				string(schedulers.SkuNameDedicated),
			}, false),
		},

		"ip_allow_list": {
			Type:     pluginsdk.TypeList,
			Required: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},

		"capacity": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntAtLeast(1),
		},

		"tags": commonschema.Tags(),
	}
}

func (r SchedulerResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"endpoint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"redundancy_state": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r SchedulerResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DurableTask.SchedulersClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model SchedulerResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := schedulers.NewSchedulerID(subscriptionId, model.ResourceGroupName, model.Name)

			metadata.Logger.Infof("Import check for %s", id)
			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			metadata.Logger.Infof("Creating %s", id)

			properties := schedulers.Scheduler{
				Location: location.Normalize(model.Location),
				Properties: &schedulers.SchedulerProperties{
					Sku: schedulers.SchedulerSku{
						Name: schedulers.SkuName(model.SkuName),
					},
					IPAllowList: &model.IpAllowList,
				},
				Tags: tags.Expand(model.Tags),
			}

			if model.Capacity != 0 {
				properties.Properties.Sku.Capacity = pointer.To(model.Capacity)
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, properties); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r SchedulerResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DurableTask.SchedulersClient

			id, err := schedulers.ParseSchedulerID(metadata.ResourceData.Id())
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

			state := SchedulerResourceModel{
				Name:              id.SchedulerName,
				ResourceGroupName: id.ResourceGroupName,
				Location:          location.Normalize(model.Location),
				Tags:              tags.Flatten(model.Tags),
			}

			if props := model.Properties; props != nil {
				state.SkuName = string(props.Sku.Name)

				if props.Sku.Capacity != nil {
					state.Capacity = *props.Sku.Capacity
				}

				if props.IPAllowList != nil {
					state.IpAllowList = *props.IPAllowList
				}

				state.Endpoint = pointer.From(props.Endpoint)
				
				if props.Sku.RedundancyState != nil {
					state.RedundancyState = string(*props.Sku.RedundancyState)
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r SchedulerResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DurableTask.SchedulersClient

			id, err := schedulers.ParseSchedulerID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model SchedulerResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			metadata.Logger.Infof("Updating %s", id)

			properties := schedulers.SchedulerUpdate{
				Properties: &schedulers.SchedulerProperties{
					Sku: schedulers.SchedulerSku{
						Name: schedulers.SkuName(model.SkuName),
					},
					IPAllowList: &model.IpAllowList,
				},
				Tags: tags.Expand(model.Tags),
			}

			if model.Capacity != 0 {
				properties.Properties.Sku.Capacity = pointer.To(model.Capacity)
			}

			if err := client.UpdateThenPoll(ctx, *id, properties); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r SchedulerResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DurableTask.SchedulersClient

			id, err := schedulers.ParseSchedulerID(metadata.ResourceData.Id())
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
