// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package kubernetes_fleet_update_strategy

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2025-03-01/fleetupdatestrategies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Create() sdk.ResourceFunc {
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

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil {
					if !response.WasNotFound(existing.HttpResponse) {
						return fmt.Errorf("checking for the presence of an existing %s: %+v", id, err)
					}
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			payload := fleetupdatestrategies.FleetUpdateStrategy{
				Properties: &fleetupdatestrategies.FleetUpdateStrategyProperties{
					Strategy: fleetupdatestrategies.UpdateRunStrategy{
						Stages: expandKubernetesFleetUpdateStrategyStage(config.Stage),
					},
				},
			}

			if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, payload, fleetupdatestrategies.DefaultCreateOrUpdateOperationOptions(), metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func expandKubernetesFleetUpdateStrategyStage(input []KubernetesFleetUpdateStrategyResourceUpdateStageSchema) []fleetupdatestrategies.UpdateStage {
	output := make([]fleetupdatestrategies.UpdateStage, 0)
	for _, stage := range input {
		output = append(output, fleetupdatestrategies.UpdateStage{
			Name:                    stage.Name,
			AfterStageWaitInSeconds: pointer.To(stage.AfterStageWaitInSeconds),
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
