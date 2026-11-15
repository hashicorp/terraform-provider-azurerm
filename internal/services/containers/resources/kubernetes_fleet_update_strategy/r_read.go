// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package kubernetes_fleet_update_strategy

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2025-03-01/fleetupdatestrategies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
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

func flattenKubernetesFleetUpdateStrategyStage(input []fleetupdatestrategies.UpdateStage) []KubernetesFleetUpdateStrategyResourceUpdateStageSchema {
	output := make([]KubernetesFleetUpdateStrategyResourceUpdateStageSchema, 0)
	for _, stage := range input {
		output = append(output, KubernetesFleetUpdateStrategyResourceUpdateStageSchema{
			Name:                    stage.Name,
			AfterStageWaitInSeconds: pointer.From(stage.AfterStageWaitInSeconds),
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
