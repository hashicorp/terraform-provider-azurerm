// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package service_plan

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := commonids.ParseAppServicePlanID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			client := metadata.Client.AppService.ServicePlanClient

			var state ServicePlanModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			model := *existing.Model

			if metadata.ResourceData.HasChange("per_site_scaling_enabled") {
				model.Properties.PerSiteScaling = pointer.To(state.PerSiteScaling)
			}

			if metadata.ResourceData.HasChange("sku_name") {
				model.Sku.Name = pointer.To(state.Sku)
			}

			if metadata.ResourceData.HasChange("tags") {
				model.Tags = pointer.To(state.Tags)
			}

			if metadata.ResourceData.HasChange("worker_count") {
				model.Sku.Capacity = pointer.To(state.WorkerCount)
			}

			if metadata.ResourceData.HasChange("premium_plan_auto_scale_enabled") {
				model.Properties.ElasticScaleEnabled = pointer.To(state.PremiumPlanAutoScaleEnabled)
			}

			if metadata.ResourceData.HasChange("maximum_elastic_worker_count") {
				model.Properties.MaximumElasticWorkerCount = pointer.To(state.MaximumElasticWorkerCount)
			}

			if metadata.ResourceData.HasChange("zone_balancing_enabled") {
				model.Properties.ZoneRedundant = pointer.To(state.ZoneBalancing)
			}

			if err = client.CreateOrUpdateThenPoll(ctx, *id, model); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}
