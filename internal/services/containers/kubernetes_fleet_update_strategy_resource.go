// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package containers

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2024-04-01/fleetupdatestrategies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.Resource = KubernetesFleetUpdateStrategyResource{}
var _ sdk.ResourceWithUpdate = KubernetesFleetUpdateStrategyResource{}

type KubernetesFleetUpdateStrategyResource struct{}

func (r KubernetesFleetUpdateStrategyResource) ModelObject() interface{} {
	return &KubernetesFleetUpdateStrategyResourceSchema{}
}

type KubernetesFleetUpdateStrategyResourceSchema struct {
	KubernetesFleetManagerId string                                                   `tfschema:"kubernetes_fleet_manager_id"`
	Name                     string                                                   `tfschema:"name"`
	Stage                    []KubernetesFleetUpdateStrategyResourceUpdateStageSchema `tfschema:"stage"`
}

type KubernetesFleetUpdateStrategyResourceUpdateGroupSchema struct {
	Name string `tfschema:"name"`
}

type KubernetesFleetUpdateStrategyResourceUpdateStageSchema struct {
	AfterStageWaitInSeconds int64                                                    `tfschema:"after_stage_wait_in_seconds"`
	Group                   []KubernetesFleetUpdateStrategyResourceUpdateGroupSchema `tfschema:"group"`
	Name                    string                                                   `tfschema:"name"`
}

func (r KubernetesFleetUpdateStrategyResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return fleetupdatestrategies.ValidateUpdateStrategyID
}

func (r KubernetesFleetUpdateStrategyResource) ResourceType() string {
	return "azurerm_kubernetes_fleet_update_strategy"
}

func (r KubernetesFleetUpdateStrategyResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			ForceNew: true,
			Required: true,
			Type:     pluginsdk.TypeString,
		},

		"kubernetes_fleet_manager_id": commonschema.ResourceIDReferenceRequiredForceNew(&fleetupdatestrategies.FleetId{}),

		"stage": {
			Required: true,
			Type:     pluginsdk.TypeList,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Required:     true,
						Type:         pluginsdk.TypeString,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"group": {
						Required: true,
						Type:     pluginsdk.TypeList,
						MinItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Required:     true,
									Type:         pluginsdk.TypeString,
									ValidateFunc: validation.StringIsNotEmpty,
								},
							},
						},
					},

					"after_stage_wait_in_seconds": {
						Optional: true,
						Type:     pluginsdk.TypeInt,
					},
				},
			},
		},
	}
}

func (r KubernetesFleetUpdateStrategyResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r KubernetesFleetUpdateStrategyResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Containers.FleetUpdateStrategiesClient

			var config KubernetesFleetUpdateStrategyResourceSchema
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			fleetId, err := commonids.ParseKubernetesFleetID(config.KubernetesFleetManagerId)
			if err != nil {
				return err
			}

			id := fleetupdatestrategies.NewUpdateStrategyID(fleetId.SubscriptionId, fleetId.ResourceGroupName, fleetId.FleetName, config.Name)

			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for the presence of an existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			payload := fleetupdatestrategies.FleetUpdateStrategy{
				Properties: &fleetupdatestrategies.FleetUpdateStrategyProperties{
					Strategy: fleetupdatestrategies.UpdateRunStrategy{
						Stages: expandKubernetesFleetUpdateStrategyStage(config.Stage),
					},
				},
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, payload, fleetupdatestrategies.DefaultCreateOrUpdateOperationOptions()); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r KubernetesFleetUpdateStrategyResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Containers.FleetUpdateStrategiesClient

			id, err := fleetupdatestrategies.ParseUpdateStrategyID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var config KubernetesFleetUpdateStrategyResourceSchema
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving existing %s: %+v", *id, err)
			}
			if existing.Model == nil {
				return fmt.Errorf("retrieving existing %s: properties was nil", *id)
			}
			payload := *existing.Model

			if metadata.ResourceData.HasChange("stage") {
				payload.Properties.Strategy.Stages = expandKubernetesFleetUpdateStrategyStage(config.Stage)
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, payload, fleetupdatestrategies.DefaultCreateOrUpdateOperationOptions()); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r KubernetesFleetUpdateStrategyResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Containers.FleetUpdateStrategiesClient
			schema := KubernetesFleetUpdateStrategyResourceSchema{}

			id, err := fleetupdatestrategies.ParseUpdateStrategyID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if model := resp.Model; model != nil {
				schema.Name = id.UpdateStrategyName
				schema.KubernetesFleetManagerId = fleetupdatestrategies.NewFleetID(id.SubscriptionId, id.ResourceGroupName, id.FleetName).ID()
				if model.Properties != nil {
					schema.Stage = flattenKubernetesFleetUpdateStrategyStage(model.Properties.Strategy.Stages)
				}
			}

			return metadata.Encode(&schema)
		},
	}
}
func (r KubernetesFleetUpdateStrategyResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Containers.FleetUpdateStrategiesClient

			id, err := fleetupdatestrategies.ParseUpdateStrategyID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id, fleetupdatestrategies.DefaultDeleteOperationOptions()); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func expandKubernetesFleetUpdateStrategyStage(input []KubernetesFleetUpdateStrategyResourceUpdateStageSchema) []fleetupdatestrategies.UpdateStage {
	output := make([]fleetupdatestrategies.UpdateStage, 0)
	for _, stage := range input {
		output = append(output, fleetupdatestrategies.UpdateStage{
			Name:                    stage.Name,
			AfterStageWaitInSeconds: pointer.FromInt64(stage.AfterStageWaitInSeconds),
			Groups:                  expandKubernetesFleetUpdateStrategyGroup(stage.Group),
		})
	}
	return output
}

func expandKubernetesFleetUpdateStrategyGroup(input []KubernetesFleetUpdateStrategyResourceUpdateGroupSchema) *[]fleetupdatestrategies.UpdateGroup {
	output := make([]fleetupdatestrategies.UpdateGroup, 0)
	for _, group := range input {
		output = append(output, fleetupdatestrategies.UpdateGroup{
			Name: group.Name,
		})
	}
	return &output
}

func flattenKubernetesFleetUpdateStrategyStage(input []fleetupdatestrategies.UpdateStage) []KubernetesFleetUpdateStrategyResourceUpdateStageSchema {
	output := make([]KubernetesFleetUpdateStrategyResourceUpdateStageSchema, 0)
	for _, stage := range input {
		output = append(output, KubernetesFleetUpdateStrategyResourceUpdateStageSchema{
			Name:                    stage.Name,
			AfterStageWaitInSeconds: pointer.ToInt64(stage.AfterStageWaitInSeconds),
			Group:                   flattenKubernetesFleetUpdateStrategyGroup(stage.Groups),
		})
	}
	return output

}

func flattenKubernetesFleetUpdateStrategyGroup(input *[]fleetupdatestrategies.UpdateGroup) []KubernetesFleetUpdateStrategyResourceUpdateGroupSchema {
	output := make([]KubernetesFleetUpdateStrategyResourceUpdateGroupSchema, 0)
	for _, group := range *input {
		output = append(output, KubernetesFleetUpdateStrategyResourceUpdateGroupSchema{
			Name: group.Name,
		})
	}
	return output
}
