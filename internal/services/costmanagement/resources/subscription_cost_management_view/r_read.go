// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package subscription_cost_management_view

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/costmanagement/2025-03-01/viewoperationgroup"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/costmanagement/helpers"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.CostManagement.ViewOperationGroupClient

			id, err := viewoperationgroup.ParseScopedViewID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.ViewsGetByScope(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := SubscriptionCostManagementViewModel{
				Name:  id.ViewName,
				SubId: id.Scope,
			}

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					state.ChartType = pointer.FromEnum(props.Chart)
					state.DisplayName = pointer.From(props.DisplayName)

					state.Accumulated = viewoperationgroup.AccumulatedTypeTrue == pointer.From(props.Accumulated)

					state.Kpi = helpers.FlattenKpisToModel(props.Kpis)
					state.Pivot = helpers.FlattenPivotsToModel(props.Pivots)

					if query := props.Query; query != nil {
						state.Timeframe = string(query.Timeframe)
						state.ReportType = string(query.Type)
						if query.DataSet != nil {
							state.Dataset = helpers.FlattenDatasetToModel(query.DataSet)
						}
					}
				}
			}

			return metadata.Encode(&state)
		},
	}
}
