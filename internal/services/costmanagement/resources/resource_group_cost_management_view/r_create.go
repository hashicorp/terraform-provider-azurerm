// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package resource_group_cost_management_view

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/costmanagement/2025-03-01/viewoperationgroup"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/costmanagement/helpers"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.CostManagement.ViewOperationGroupClient

			var config ResourceGroupCostManagementViewModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := viewoperationgroup.NewScopedViewID(config.ResourceGroupId, config.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.ViewsGetByScope(ctx, id)
				if err != nil {
					if !response.WasNotFound(existing.HttpResponse) {
						return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
					}
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return tf.ImportAsExistsError(r.ResourceType(), id.ID())
				}
			}

			accumulated := viewoperationgroup.AccumulatedTypeFalse
			if config.Accumulated {
				accumulated = viewoperationgroup.AccumulatedTypeTrue
			}

			props := viewoperationgroup.View{
				Properties: &viewoperationgroup.ViewProperties{
					Accumulated: pointer.To(accumulated),
					DisplayName: pointer.To(config.DisplayName),
					Chart:       pointer.ToEnum[viewoperationgroup.ChartType](config.ChartType),
					Query: &viewoperationgroup.ReportConfigDefinition{
						DataSet:   helpers.ExpandDatasetFromModel(config.Dataset),
						Timeframe: viewoperationgroup.ReportTimeframeType(config.Timeframe),
						Type:      viewoperationgroup.ReportTypeUsage,
					},
					Kpis:   helpers.ExpandKpisFromModel(config.Kpi),
					Pivots: helpers.ExpandPivotsFromModel(config.Pivot),
				},
			}

			if _, err := client.ViewsCreateOrUpdateByScope(ctx, id, props); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}
