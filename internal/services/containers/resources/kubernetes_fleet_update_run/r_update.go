// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package kubernetes_fleet_update_run

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2025-03-01/updateruns"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Containers.FleetUpdateRunsClient

			id, err := updateruns.ParseUpdateRunID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var config KubernetesFleetUpdateRunResourceSchema
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: `Model` was nil", *id)
			}

			if existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `Properties` was nil", *id)
			}
			props := existing.Model.Properties

			// The API errors if this field is included in the Update request
			// unexpected status 400 (400 Bad Request) with error: PropertyChangeNotAllowed:
			// Changing immutable/read-only properties is not allowed: Status
			props.Status = nil

			if metadata.ResourceData.HasChange("managed_cluster_update") {
				props.ManagedClusterUpdate = expandKubernetesFleetUpdateRunManagedClusterUpdate(config.ManagedClusterUpdate)
			}

			if metadata.ResourceData.HasChange("fleet_update_strategy_id") {
				props.UpdateStrategyId = pointer.To(config.FleetUpdateStrategyId)
			}

			if metadata.ResourceData.HasChange("stage") {
				if props.Strategy == nil {
					props.Strategy = &updateruns.UpdateRunStrategy{}
				}
				props.Strategy.Stages = expandKubernetesFleetUpdateRunStage(config.Stage)
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, *existing.Model, updateruns.DefaultCreateOrUpdateOperationOptions()); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}
