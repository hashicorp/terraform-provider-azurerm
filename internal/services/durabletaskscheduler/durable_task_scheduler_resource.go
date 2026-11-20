// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package durabletaskscheduler

//go:generate go run ../../tools/generator-tests resourceidentity

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
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/durabletaskscheduler/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type SchedulerResourceModel struct {
	Name              string              `tfschema:"name"`
	ResourceGroupName string              `tfschema:"resource_group_name"`
	Location          string              `tfschema:"location"`
	Sku               []SchedulerSkuModel `tfschema:"sku"`
	IpAllowList       []string            `tfschema:"ip_allowlist"`
	Tags              map[string]string   `tfschema:"tags"`
	Endpoint          string              `tfschema:"endpoint"`
}

type SchedulerSkuModel struct {
	Name     string `tfschema:"name"`
	Capacity int64  `tfschema:"capacity"`
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
			ValidateFunc: validate.DurableTaskName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"sku": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringInSlice(schedulers.PossibleValuesForSchedulerSkuName(), false),
					},
					"capacity": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
						// The Dedicated scheduler capacity bounds are taken from the Azure portal UX
						// because the current Durable Task SDK model does not surface the limit.
						ValidateFunc: validation.IntBetween(1, 3),
					},
				},
			},
		},

		"ip_allowlist": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			// NOTE: O+C SDKv2 list schemas cannot define defaults. Azure persists `0.0.0.0/0`
			// when the property is omitted, so Computed avoids a post-create diff.
			Computed: true,
			MinItems: 1,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.Any(validation.IsIPAddress, validation.IsCIDR),
			},
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

			rawSku := metadata.ResourceDiff.GetRawConfig().AsValueMap()["sku"]
			if rawSku.IsNull() || !rawSku.IsKnown() {
				return nil
			}

			sku := rawSku.AsValueSlice()
			if len(sku) != 1 || !sku[0].IsKnown() {
				return nil
			}

			skuDetails := sku[0].AsValueMap()
			rawName := skuDetails["name"]
			rawCapacity := skuDetails["capacity"]
			if rawName.IsNull() || !rawName.IsKnown() || !rawCapacity.IsKnown() {
				return nil
			}

			skuName := rawName.AsString()
			hasCapacity := !rawCapacity.IsNull()

			if skuName == string(schedulers.SchedulerSkuNameDedicated) && !hasCapacity {
				return errors.New("`sku.0.capacity` must be configured when `sku.0.name` is set to `Dedicated`")
			}

			if hasCapacity && skuName != string(schedulers.SchedulerSkuNameDedicated) {
				return errors.New("`sku.0.capacity` can only be configured when `sku.0.name` is set to `Dedicated`")
			}

			return nil
		},
	}
}

func (r SchedulerResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DurableTaskScheduler.SchedulersClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model SchedulerResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}
			if len(model.Sku) != 1 {
				return errors.New("decoding: expected one `sku` block")
			}

			id := schedulers.NewSchedulerID(subscriptionId, model.ResourceGroupName, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			ipAllowList := model.IpAllowList
			if len(ipAllowList) == 0 {
				ipAllowList = []string{"0.0.0.0/0"}
			}
			sku := model.Sku[0]

			properties := schedulers.Scheduler{
				Location: location.Normalize(model.Location),
				Properties: &schedulers.SchedulerProperties{
					Sku: schedulers.SchedulerSku{
						Name:     schedulers.SchedulerSkuName(sku.Name),
						Capacity: pointer.ToOrNil(sku.Capacity),
					},
					IPAllowlist: ipAllowList,
				},
				Tags: &model.Tags,
			}

			if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, properties, metadata.SetIDAndIdentityCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return pluginsdk.SetResourceIdentityData(metadata.ResourceData, &id)
		},
	}
}

func (r SchedulerResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DurableTaskScheduler.SchedulersClient

			id, err := schedulers.ParseSchedulerID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := flattenScheduler(*id, resp.Model)

			if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
				return err
			}

			return metadata.Encode(&state)
		},
	}
}

func flattenScheduler(id schedulers.SchedulerId, model *schedulers.Scheduler) SchedulerResourceModel {
	state := SchedulerResourceModel{
		Name:              id.SchedulerName,
		ResourceGroupName: id.ResourceGroupName,
	}

	if model != nil {
		state.Location = location.Normalize(model.Location)
		state.Tags = pointer.From(model.Tags)

		if props := model.Properties; props != nil {
			state.Sku = []SchedulerSkuModel{{
				Name:     string(props.Sku.Name),
				Capacity: pointer.From(props.Sku.Capacity),
			}}
			state.IpAllowList = props.IPAllowlist
			state.Endpoint = pointer.From(props.Endpoint)
		}
	}

	return state
}

func (r SchedulerResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DurableTaskScheduler.SchedulersClient

			id, err := schedulers.ParseSchedulerID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model SchedulerResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}
			if len(model.Sku) != 1 {
				return errors.New("decoding: expected one `sku` block")
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}
			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: model was nil", *id)
			}
			if existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: properties was nil", *id)
			}

			payload := *existing.Model
			sku := model.Sku[0]

			payload.Properties.IPAllowlist = model.IpAllowList
			payload.Properties.Sku.Capacity = pointer.ToOrNil(sku.Capacity)
			payload.Tags = &model.Tags

			if err := client.CreateOrUpdateThenPoll(ctx, *id, payload); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r SchedulerResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DurableTaskScheduler.SchedulersClient

			id, err := schedulers.ParseSchedulerID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}
