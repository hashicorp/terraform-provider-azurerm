// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package kubernetes_fleet_update_run

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2025-03-01/updateruns"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Containers.FleetUpdateRunsClient

			var config KubernetesFleetUpdateRunResourceSchema
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			fleetId, err := commonids.ParseKubernetesFleetID(config.KubernetesFleetManagerId)
			if err != nil {
				return err
			}

			id := updateruns.NewUpdateRunID(fleetId.SubscriptionId, fleetId.ResourceGroupName, fleetId.FleetName, config.Name)

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

			payload := updateruns.UpdateRun{
				Properties: &updateruns.UpdateRunProperties{
					ManagedClusterUpdate: expandKubernetesFleetUpdateRunManagedClusterUpdate(config.ManagedClusterUpdate),
					Strategy:             expandKubernetesFleetUpdateRunStrategy(config.Stage),
					UpdateStrategyId:     pointer.To(config.FleetUpdateStrategyId),
				},
			}

			if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, payload, updateruns.DefaultCreateOrUpdateOperationOptions(), metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func expandKubernetesFleetUpdateRunManagedClusterUpdate(input []KubernetesFleetUpdateRunResourceManagedClusterUpdateSchema) updateruns.ManagedClusterUpdate {
	if len(input) == 0 {
		return updateruns.ManagedClusterUpdate{}
	}
	return updateruns.ManagedClusterUpdate{
		NodeImageSelection: expandKubernetesFleetUpdateRunNodeImageSelection(input[0].NodeImageSelection),
		Upgrade:            expandKubernetesFleetUpdateRunUpgrade(input[0].Upgrade),
	}
}

func expandKubernetesFleetUpdateRunUpgrade(input []KubernetesFleetUpdateRunResourceManagedClusterUpdateUpgradeSchema) updateruns.ManagedClusterUpgradeSpec {
	if len(input) == 0 {
		return updateruns.ManagedClusterUpgradeSpec{}
	}
	return updateruns.ManagedClusterUpgradeSpec{
		Type:              updateruns.ManagedClusterUpgradeType(input[0].Type),
		KubernetesVersion: pointer.To(input[0].KubernetesVersion),
	}
}

func expandKubernetesFleetUpdateRunNodeImageSelection(input []KubernetesFleetUpdateRunResourceManagedClusterUpdateNodeImageSelectionSchema) *updateruns.NodeImageSelection {
	if len(input) == 0 {
		return nil
	}
	return &updateruns.NodeImageSelection{
		Type: updateruns.NodeImageSelectionType(input[0].Type),
	}
}

func expandKubernetesFleetUpdateRunStrategy(input []KubernetesFleetUpdateRunResourceUpdateStageSchema) *updateruns.UpdateRunStrategy {
	if len(input) == 0 {
		return nil
	}

	return &updateruns.UpdateRunStrategy{
		Stages: expandKubernetesFleetUpdateRunStage(input),
	}
}

func expandKubernetesFleetUpdateRunStage(input []KubernetesFleetUpdateRunResourceUpdateStageSchema) []updateruns.UpdateStage {
	output := make([]updateruns.UpdateStage, 0)
	for _, stage := range input {
		output = append(output, updateruns.UpdateStage{
			Name:                    stage.Name,
			AfterStageWaitInSeconds: pointer.To(stage.AfterStageWaitInSeconds),
			Groups:                  expandKubernetesFleetUpdateRunGroup(stage.Group),
		})
	}
	return output
}

func expandKubernetesFleetUpdateRunGroup(input []KubernetesFleetUpdateRunResourceUpdateGroupSchema) *[]updateruns.UpdateGroup {
	output := make([]updateruns.UpdateGroup, 0)
	for _, group := range input {
		output = append(output, updateruns.UpdateGroup{
			Name: group.Name,
		})
	}
	return &output
}
