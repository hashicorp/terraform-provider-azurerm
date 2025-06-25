// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-11-01/virtualmachinescalesets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/standbypool/2025-03-01/standbyvirtualmachinepools"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type VirtualMachineScaleSetStandbyPoolModel struct {
	Name                             string                                                    `tfschema:"name"`
	ResourceGroupName                string                                                    `tfschema:"resource_group_name"`
	Location                         string                                                    `tfschema:"location"`
	AttachedVirtualMachineScaleSetId string                                                    `tfschema:"attached_virtual_machine_scale_set_id"`
	ElasticityProfile                []VirtualMachineScaleSetStandbyPoolElasticityProfileModel `tfschema:"elasticity_profile"`
	VirtualMachineState              standbyvirtualmachinepools.VirtualMachineState            `tfschema:"virtual_machine_state"`
	Tags                             map[string]string                                         `tfschema:"tags"`
}

type VirtualMachineScaleSetStandbyPoolElasticityProfileModel struct {
	MaxReadyCapacity int64 `tfschema:"max_ready_capacity"`
	MinReadyCapacity int64 `tfschema:"min_ready_capacity"`
}

type VirtualMachineScaleSetStandbyPoolResource struct{}

var (
	_ sdk.ResourceWithUpdate        = VirtualMachineScaleSetStandbyPoolResource{}
	_ sdk.ResourceWithCustomizeDiff = VirtualMachineScaleSetStandbyPoolResource{}
)

func (r VirtualMachineScaleSetStandbyPoolResource) ResourceType() string {
	return "azurerm_virtual_machine_scale_set_standby_pool"
}

func (r VirtualMachineScaleSetStandbyPoolResource) ModelObject() interface{} {
	return &VirtualMachineScaleSetStandbyPoolModel{}
}

func (r VirtualMachineScaleSetStandbyPoolResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return standbyvirtualmachinepools.ValidateStandbyVirtualMachinePoolID
}

func (r VirtualMachineScaleSetStandbyPoolResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile("^[a-zA-Z0-9-]{3,24}$"),
				"name must be between 3 and 24 characters in length and may contain only letters, numbers and hyphens (-).",
			),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"attached_virtual_machine_scale_set_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: virtualmachinescalesets.ValidateVirtualMachineScaleSetID,
		},

		"elasticity_profile": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"max_ready_capacity": {
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntBetween(0, 2000),
					},

					"min_ready_capacity": {
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntBetween(0, 2000),
					},
				},
			},
		},

		"virtual_machine_state": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(standbyvirtualmachinepools.PossibleValuesForVirtualMachineState(), false),
		},

		"tags": commonschema.Tags(),
	}
}

func (r VirtualMachineScaleSetStandbyPoolResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r VirtualMachineScaleSetStandbyPoolResource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var config VirtualMachineScaleSetStandbyPoolModel
			if err := metadata.DecodeDiff(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			if len(config.ElasticityProfile) > 0 && config.ElasticityProfile[0].MaxReadyCapacity < config.ElasticityProfile[0].MinReadyCapacity {
				return fmt.Errorf("`min_ready_capacity` cannot exceed `max_ready_capacity`")
			}

			return nil
		},
	}
}

func (r VirtualMachineScaleSetStandbyPoolResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Compute.StandbyVirtualMachinePoolsClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model VirtualMachineScaleSetStandbyPoolModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := standbyvirtualmachinepools.NewStandbyVirtualMachinePoolID(subscriptionId, model.ResourceGroupName, model.Name)
			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			properties := &standbyvirtualmachinepools.StandbyVirtualMachinePoolResource{
				Location: location.Normalize(model.Location),
				Properties: &standbyvirtualmachinepools.StandbyVirtualMachinePoolResourceProperties{
					AttachedVirtualMachineScaleSetId: pointer.To(model.AttachedVirtualMachineScaleSetId),
					ElasticityProfile:                expandStandbyVirtualMachinePoolElasticityProfileModel(model.ElasticityProfile),
					VirtualMachineState:              model.VirtualMachineState,
				},
				Tags: &model.Tags,
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, *properties); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r VirtualMachineScaleSetStandbyPoolResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Compute.StandbyVirtualMachinePoolsClient

			id, err := standbyvirtualmachinepools.ParseStandbyVirtualMachinePoolID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model VirtualMachineScaleSetStandbyPoolModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			properties := resp.Model
			if properties == nil {
				return fmt.Errorf("retrieving %s: properties was nil", id)
			}

			if metadata.ResourceData.HasChange("attached_virtual_machine_scale_set_id") {
				properties.Properties.AttachedVirtualMachineScaleSetId = pointer.To(model.AttachedVirtualMachineScaleSetId)
			}

			if metadata.ResourceData.HasChange("elasticity_profile") {
				properties.Properties.ElasticityProfile = expandStandbyVirtualMachinePoolElasticityProfileModel(model.ElasticityProfile)
			}

			if metadata.ResourceData.HasChange("virtual_machine_state") {
				properties.Properties.VirtualMachineState = model.VirtualMachineState
			}

			if metadata.ResourceData.HasChange("tags") {
				properties.Tags = &model.Tags
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, *properties); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r VirtualMachineScaleSetStandbyPoolResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Compute.StandbyVirtualMachinePoolsClient

			id, err := standbyvirtualmachinepools.ParseStandbyVirtualMachinePoolID(metadata.ResourceData.Id())
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

			state := VirtualMachineScaleSetStandbyPoolModel{
				Name:              id.StandbyVirtualMachinePoolName,
				ResourceGroupName: id.ResourceGroupName,
			}

			if model := resp.Model; model != nil {
				state.Location = location.Normalize(model.Location)
				if properties := model.Properties; properties != nil {
					parsedAttachedVirtualMachineScaleSetId, err := virtualmachinescalesets.ParseVirtualMachineScaleSetIDInsensitively(pointer.From(properties.AttachedVirtualMachineScaleSetId))
					if err != nil {
						return fmt.Errorf("parsing `attached_virtual_machine_scale_set_id` for %s: %+v", *id, err)
					}

					state.AttachedVirtualMachineScaleSetId = parsedAttachedVirtualMachineScaleSetId.ID()
					state.ElasticityProfile = flattenStandbyVirtualMachinePoolElasticityProfileModel(properties.ElasticityProfile)
					state.VirtualMachineState = properties.VirtualMachineState
				}

				state.Tags = pointer.From(model.Tags)
			}

			return metadata.Encode(&state)
		},
	}
}

func (r VirtualMachineScaleSetStandbyPoolResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Compute.StandbyVirtualMachinePoolsClient

			id, err := standbyvirtualmachinepools.ParseStandbyVirtualMachinePoolID(metadata.ResourceData.Id())
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

func expandStandbyVirtualMachinePoolElasticityProfileModel(inputList []VirtualMachineScaleSetStandbyPoolElasticityProfileModel) *standbyvirtualmachinepools.StandbyVirtualMachinePoolElasticityProfile {
	if len(inputList) == 0 {
		return nil
	}

	input := &inputList[0]
	output := standbyvirtualmachinepools.StandbyVirtualMachinePoolElasticityProfile{
		MaxReadyCapacity: input.MaxReadyCapacity,
		MinReadyCapacity: pointer.To(input.MinReadyCapacity),
	}

	return &output
}

func flattenStandbyVirtualMachinePoolElasticityProfileModel(input *standbyvirtualmachinepools.StandbyVirtualMachinePoolElasticityProfile) []VirtualMachineScaleSetStandbyPoolElasticityProfileModel {
	outputList := make([]VirtualMachineScaleSetStandbyPoolElasticityProfileModel, 0)
	if input == nil {
		return outputList
	}

	output := VirtualMachineScaleSetStandbyPoolElasticityProfileModel{
		MaxReadyCapacity: input.MaxReadyCapacity,
		MinReadyCapacity: pointer.From(input.MinReadyCapacity),
	}

	return append(outputList, output)
}
