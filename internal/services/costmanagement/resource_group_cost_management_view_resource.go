// Copyright IBM Corp. 2023, 2026
// SPDX-License-Identifier: MPL-2.0

package costmanagement

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/costmanagement/2023-08-01/views"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/costmanagement/validate"
	resourceValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type ResourceGroupCostManagementViewResource struct {
	base costManagementViewBaseResource
}

type ResourceGroupCostManagementViewModel struct {
	Name            string                           `tfschema:"name"`
	ResourceGroupId string                           `tfschema:"resource_group_id"`
	DisplayName     string                           `tfschema:"display_name"`
	ChartType       string                           `tfschema:"chart_type"`
	Accumulated     bool                             `tfschema:"accumulated"`
	ReportType      string                           `tfschema:"report_type"`
	Timeframe       string                           `tfschema:"timeframe"`
	Dataset         []CostManagementViewDatasetModel `tfschema:"dataset"`
	Kpi             []CostManagementViewKpiModel     `tfschema:"kpi"`
	Pivot           []CostManagementViewPivotModel   `tfschema:"pivot"`
}

var _ sdk.Resource = ResourceGroupCostManagementViewResource{}

func (r ResourceGroupCostManagementViewResource) Arguments() map[string]*pluginsdk.Schema {
	schema := map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotWhiteSpace,
		},
		"resource_group_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: resourceValidate.ResourceGroupID,
		},
	}
	return r.base.arguments(schema)
}

func (r ResourceGroupCostManagementViewResource) Attributes() map[string]*pluginsdk.Schema {
	return r.base.attributes()
}

func (r ResourceGroupCostManagementViewResource) ModelObject() interface{} {
	return &ResourceGroupCostManagementViewModel{}
}

func (r ResourceGroupCostManagementViewResource) ResourceType() string {
	return "azurerm_resource_group_cost_management_view"
}

func (r ResourceGroupCostManagementViewResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.ResourceGroupCostManagementViewID
}

func (r ResourceGroupCostManagementViewResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.CostManagement.ViewsClient

			var config ResourceGroupCostManagementViewModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := views.NewScopedViewID(config.ResourceGroupId, config.Name)

			existing, err := client.GetByScope(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError(r.ResourceType(), id.ID())
			}

			accumulated := views.AccumulatedTypeFalse
			if config.Accumulated {
				accumulated = views.AccumulatedTypeTrue
			}

			props := views.View{
				Properties: &views.ViewProperties{
					Accumulated: pointer.To(accumulated),
					DisplayName: pointer.To(config.DisplayName),
					Chart:       pointer.To(views.ChartType(config.ChartType)),
					Query: &views.ReportConfigDefinition{
						DataSet:   expandDatasetFromModel(config.Dataset),
						Timeframe: views.ReportTimeframeType(config.Timeframe),
						Type:      views.ReportTypeUsage,
					},
					Kpis:   expandKpisFromModel(config.Kpi),
					Pivots: expandPivotsFromModel(config.Pivot),
				},
			}

			if _, err = client.CreateOrUpdateByScope(ctx, id, props); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r ResourceGroupCostManagementViewResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.CostManagement.ViewsClient

			id, err := views.ParseScopedViewID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.GetByScope(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := ResourceGroupCostManagementViewModel{
				Name:            id.ViewName,
				ResourceGroupId: id.Scope,
			}

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					state.ChartType = pointer.FromEnum(props.Chart)
					state.DisplayName = pointer.From(props.DisplayName)

					state.Accumulated = views.AccumulatedTypeTrue == pointer.From(props.Accumulated)

					state.Kpi = flattenKpisToModel(props.Kpis)
					state.Pivot = flattenPivotsToModel(props.Pivots)

					if query := props.Query; query != nil {
						state.Timeframe = string(query.Timeframe)
						state.ReportType = string(query.Type)
						if query.DataSet != nil {
							state.Dataset = flattenDatasetToModel(query.DataSet)
						}
					}
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r ResourceGroupCostManagementViewResource) Delete() sdk.ResourceFunc {
	return r.base.deleteFunc()
}

func (r ResourceGroupCostManagementViewResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.CostManagement.ViewsClient

			id, err := views.ParseScopedViewID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var config ResourceGroupCostManagementViewModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			// Update operation requires latest eTag to be set in the request.
			existing, err := client.GetByScope(ctx, *id)
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
				model.Properties.Chart = pointer.To(views.ChartType(config.ChartType))
			}

			if metadata.ResourceData.HasChange("dataset") {
				model.Properties.Query.DataSet = expandDatasetFromModel(config.Dataset)
			}

			if metadata.ResourceData.HasChange("timeframe") {
				model.Properties.Query.Timeframe = views.ReportTimeframeType(config.Timeframe)
			}

			if metadata.ResourceData.HasChange("kpi") {
				model.Properties.Kpis = expandKpisFromModel(config.Kpi)
			}

			if metadata.ResourceData.HasChange("pivot") {
				model.Properties.Pivots = expandPivotsFromModel(config.Pivot)
			}

			if _, err = client.CreateOrUpdateByScope(ctx, *id, *model); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}
