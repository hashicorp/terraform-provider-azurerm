// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package resource_group_cost_management_view

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/costmanagement/2025-03-01/viewoperationgroup"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/costmanagement/helpers"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.CostManagement.ViewOperationGroupClient

			id, err := viewoperationgroup.ParseScopedViewID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var config ResourceGroupCostManagementViewModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			// Update operation requires latest eTag to be set in the request.
			existing, err := client.ViewsGetByScope(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}
			model := existing.Model

			if model != nil {
				if model.ETag == nil {
					return fmt.Errorf("add %s: etag was nil", *id)
				}
			}

			if model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", *id)
			}

			if metadata.ResourceData.HasChange("display_name") {
				model.Properties.DisplayName = pointer.To(config.DisplayName)
			}

			if metadata.ResourceData.HasChange("chart_type") {
				model.Properties.Chart = pointer.ToEnum[viewoperationgroup.ChartType](config.ChartType)
			}

			if metadata.ResourceData.HasChange("dataset") {
				model.Properties.Query.DataSet = helpers.ExpandDatasetFromModel(config.Dataset)
			}

			if metadata.ResourceData.HasChange("timeframe") {
				model.Properties.Query.Timeframe = viewoperationgroup.ReportTimeframeType(config.Timeframe)
			}

			if metadata.ResourceData.HasChange("kpi") {
				model.Properties.Kpis = helpers.ExpandKpisFromModel(config.Kpi)
			}

			if metadata.ResourceData.HasChange("pivot") {
				model.Properties.Pivots = helpers.ExpandPivotsFromModel(config.Pivot)
			}

			if _, err = client.ViewsCreateOrUpdateByScope(ctx, *id, *model); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}
