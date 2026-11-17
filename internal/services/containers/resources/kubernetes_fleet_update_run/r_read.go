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

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Containers.FleetUpdateRunsClient
			schema := KubernetesFleetUpdateRunResourceSchema{}

			id, err := updateruns.ParseUpdateRunID(metadata.ResourceData.Id())
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
				schema.Name = id.UpdateRunName
				schema.KubernetesFleetManagerId = commonids.NewKubernetesFleetID(id.SubscriptionId, id.ResourceGroupName, id.FleetName).ID()

				if props := model.Properties; props != nil {
					schema.ManagedClusterUpdate = flattenKubernetesFleetUpdateRunManagedClusterUpdate(props.ManagedClusterUpdate)
					schema.FleetUpdateStrategyId = pointer.From(props.UpdateStrategyId)

					if props.Strategy != nil {
						schema.Stage = flattenKubernetesFleetUpdateRunStage(props.Strategy.Stages)
					}
				}
			}

			return metadata.Encode(&schema)
		},
	}
}

type KubernetesFleetUpdateRunResourceSchema struct {
	KubernetesFleetManagerId string                                                       `tfschema:"kubernetes_fleet_manager_id"`
	Name                     string                                                       `tfschema:"name"`
	ManagedClusterUpdate     []KubernetesFleetUpdateRunResourceManagedClusterUpdateSchema `tfschema:"managed_cluster_update"`
	FleetUpdateStrategyId    string                                                       `tfschema:"fleet_update_strategy_id"`
	Stage                    []KubernetesFleetUpdateRunResourceUpdateStageSchema          `tfschema:"stage"`
}

type KubernetesFleetUpdateRunResourceManagedClusterUpdateSchema struct {
	NodeImageSelection []KubernetesFleetUpdateRunResourceManagedClusterUpdateNodeImageSelectionSchema `tfschema:"node_image_selection"`
	Upgrade            []KubernetesFleetUpdateRunResourceManagedClusterUpdateUpgradeSchema            `tfschema:"upgrade"`
}

type KubernetesFleetUpdateRunResourceManagedClusterUpdateNodeImageSelectionSchema struct {
	Type string `tfschema:"type"`
}

type KubernetesFleetUpdateRunResourceManagedClusterUpdateUpgradeSchema struct {
	Type              string `tfschema:"type"`
	KubernetesVersion string `tfschema:"kubernetes_version"`
}

type KubernetesFleetUpdateRunResourceUpdateGroupSchema struct {
	Name string `tfschema:"name"`
}

type KubernetesFleetUpdateRunResourceUpdateStageSchema struct {
	AfterStageWaitInSeconds int64                                               `tfschema:"after_stage_wait_in_seconds"`
	Group                   []KubernetesFleetUpdateRunResourceUpdateGroupSchema `tfschema:"group"`
	Name                    string                                              `tfschema:"name"`
}
