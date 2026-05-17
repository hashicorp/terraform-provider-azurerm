// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package durabletask

//go:generate go run ../../tools/generator-tests resourceidentity -resource-name durable_task_scheduler -service-package-name durabletask -properties "resource_group_name,name" -known-values "subscription_id:data.Subscriptions.Primary"

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/durabletask/2025-11-01/schedulers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/durabletask/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type SchedulerResourceModel struct {
	Name              string            `tfschema:"name"`
	ResourceGroupName string            `tfschema:"resource_group_name"`
	Location          string            `tfschema:"location"`
	IpAllowList       []string          `tfschema:"ip_allow_list"`
	SkuName           string            `tfschema:"sku_name"`
	Capacity          int64             `tfschema:"capacity"`
	Tags              map[string]string `tfschema:"tags"`
	Endpoint          string            `tfschema:"endpoint"`
}

type SchedulerResource struct{}

var (
	_ sdk.Resource                  = SchedulerResource{}
	_ sdk.ResourceWithUpdate        = SchedulerResource{}
	_ sdk.ResourceWithCustomizeDiff = SchedulerResource{}
	_ sdk.ResourceWithIdentity      = SchedulerResource{}
)

func (r SchedulerResource) Identity() resourceids.ResourceId {
	return &schedulers.SchedulerId{}
}

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
			ValidateFunc: validate.SchedulerName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"ip_allow_list": {
			Type:     pluginsdk.TypeList,
			Required: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.Any(validation.IsIPAddress, validation.IsCIDR),
			},
		},

		"sku_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(schedulers.PossibleValuesForSchedulerSkuName(), false),
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
	}
}

func (r SchedulerResource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			if metadata.ResourceDiff == nil {
				return nil
			}

			rawConfig := metadata.ResourceDiff.GetRawConfig().AsValueMap()
			rawCapacity := rawConfig["capacity"]

			if !rawCapacity.IsNull() {
				skuName := metadata.ResourceDiff.Get("sku_name").(string)
				if skuName != string(schedulers.SchedulerSkuNameDedicated) {
					return errors.New("`capacity` can only be configured when `sku_name` is set to `Dedicated`")
				}
			}

			return nil
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
						Name: schedulers.SchedulerSkuName(model.SkuName),
					},
					IPAllowlist: model.IpAllowList,
				},
				Tags: &model.Tags,
			}

			if model.Capacity != 0 {
				properties.Properties.Sku.Capacity = pointer.To(model.Capacity)
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, properties); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, &id); err != nil {
				return err
			}
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
			}

			state.Tags = pointer.From(model.Tags)

			if props := model.Properties; props != nil {
				state.SkuName = string(props.Sku.Name)
				state.Capacity = pointer.From(props.Sku.Capacity)
				state.IpAllowList = props.IPAllowlist
				state.Endpoint = pointer.From(props.Endpoint)
			}

			if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
				return err
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

			if metadata.ResourceData.HasChange("capacity") {
				rawCapacity := metadata.ResourceData.GetRawConfig().AsValueMap()["capacity"]
				if rawCapacity.IsNull() {
					properties := schedulers.Scheduler{
						Location: location.Normalize(model.Location),
						Properties: &schedulers.SchedulerProperties{
							Sku: schedulers.SchedulerSku{
								Name: schedulers.SchedulerSkuName(model.SkuName),
							},
							IPAllowlist: model.IpAllowList,
						},
						Tags: &model.Tags,
					}

					if err := client.CreateOrUpdateThenPoll(ctx, *id, properties); err != nil {
						return fmt.Errorf("updating %s: %+v", id, err)
					}

					return nil
				}
			}

			properties := schedulers.SchedulerUpdate{
				Properties: &schedulers.SchedulerPropertiesUpdate{},
			}

			if metadata.ResourceData.HasChange("ip_allow_list") {
				properties.Properties.IPAllowlist = &model.IpAllowList
			}

			if metadata.ResourceData.HasChange("capacity") {
				if properties.Properties.Sku == nil {
					properties.Properties.Sku = &schedulers.SchedulerSkuUpdate{}
				}
				properties.Properties.Sku.Capacity = pointer.To(model.Capacity)
			}

			if metadata.ResourceData.HasChange("tags") {
				properties.Tags = &model.Tags
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
