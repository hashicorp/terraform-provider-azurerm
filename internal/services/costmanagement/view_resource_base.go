// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package costmanagement

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/costmanagement/2022-10-01/views"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type costManagementViewBaseResource struct{}

func (br costManagementViewBaseResource) arguments(fields map[string]*pluginsdk.Schema) map[string]*pluginsdk.Schema {
	output := map[string]*pluginsdk.Schema{
		"display_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"chart_type": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(views.PossibleValuesForChartType(), false),
		},

		"accumulated": {
			Type:     pluginsdk.TypeBool,
			Required: true,
			ForceNew: true,
		},

		"report_type": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(views.ReportTypeUsage),
			}, false),
		},

		"timeframe": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(views.PossibleValuesForReportTimeframeType(), false),
		},

		"dataset": {
			Type:     pluginsdk.TypeList,
			MaxItems: 1,
			Required: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"granularity": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(views.PossibleValuesForReportGranularityType(), false),
					},

					"aggregation": {
						Type:     pluginsdk.TypeSet,
						Required: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Type:     pluginsdk.TypeString,
									Required: true,
									ForceNew: true,
								},
								"column_name": {
									Type:     pluginsdk.TypeString,
									Required: true,
									ForceNew: true,
								},
							},
						},
					},

					"sorting": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"direction": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringInSlice(views.PossibleValuesForReportConfigSortingType(), false),
								},
								"name": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},
							},
						},
					},

					"grouping": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"type": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringInSlice(views.PossibleValuesForQueryColumnType(), false),
								},
								"name": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},
							},
						},
					},
				},
			},
		},

		"kpi": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"type": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(views.PossibleValuesForKpiTypeType(), false),
					},
				},
			},
		},

		"pivot": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"type": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(views.PossibleValuesForPivotTypeType(), false),
					},
				},
			},
		},
	}

	for k, v := range fields {
		output[k] = v
	}

	return output
}

func (br costManagementViewBaseResource) attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (br costManagementViewBaseResource) createFunc(resourceName, scopeFieldName string) sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.CostManagement.ViewsClient
			id := views.NewScopedViewID(metadata.ResourceData.Get(scopeFieldName).(string), metadata.ResourceData.Get("name").(string))

			existing, err := client.GetByScope(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError(resourceName, id.ID())
			}

			accumulated := views.AccumulatedTypeFalse
			if accumulatedRaw := metadata.ResourceData.Get("accumulated").(bool); accumulatedRaw {
				accumulated = views.AccumulatedTypeTrue
			}

			props := views.View{
				Properties: &views.ViewProperties{
					Accumulated: pointer.To(accumulated),
					DisplayName: utils.String(metadata.ResourceData.Get("display_name").(string)),
					Chart:       pointer.To(views.ChartType(metadata.ResourceData.Get("chart_type").(string))),
					Query: &views.ReportConfigDefinition{
						DataSet:   expandDataset(metadata.ResourceData.Get("dataset").([]interface{})),
						Timeframe: views.ReportTimeframeType(metadata.ResourceData.Get("timeframe").(string)),
						Type:      views.ReportTypeUsage,
					},
					Kpis:   expandKpis(metadata.ResourceData.Get("kpi").([]interface{})),
					Pivots: expandPivots(metadata.ResourceData.Get("pivot").([]interface{})),
				},
			}

			_, err = client.CreateOrUpdateByScope(ctx, id, props)
			if err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (br costManagementViewBaseResource) readFunc(scopeFieldName string) sdk.ResourceFunc {
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
				return fmt.Errorf("reading %s: %+v", *id, err)
			}

			metadata.ResourceData.Set("name", id.ViewName)
			// lintignore:R001
			metadata.ResourceData.Set(scopeFieldName, id.Scope)

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					metadata.ResourceData.Set("chart_type", string(pointer.From(props.Chart)))

					accumulated := false
					if props.Accumulated != nil {
						accumulated = views.AccumulatedTypeTrue == *props.Accumulated
					}
					metadata.ResourceData.Set("accumulated", accumulated)

					metadata.ResourceData.Set("display_name", props.DisplayName)
					metadata.ResourceData.Set("kpi", flattenKpis(props.Kpis))
					metadata.ResourceData.Set("pivot", flattenPivots(props.Pivots))

					if query := props.Query; query != nil {
						metadata.ResourceData.Set("timeframe", string(query.Timeframe))
						metadata.ResourceData.Set("report_type", string(query.Type))
						if query.DataSet != nil {
							metadata.ResourceData.Set("dataset", flattenDataset(query.DataSet))
						}
					}
				}
			}

			return nil
		},
	}
}

func (br costManagementViewBaseResource) deleteFunc() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.CostManagement.ViewsClient

			id, err := views.ParseScopedViewID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err = client.DeleteByScope(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (br costManagementViewBaseResource) updateFunc() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.CostManagement.ViewsClient

			id, err := views.ParseScopedViewID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			// Update operation requires latest eTag to be set in the request.
			existing, err := client.GetByScope(ctx, *id)
			if err != nil {
				return fmt.Errorf("reading %s: %+v", *id, err)
			}
			model := existing.Model

			if model != nil {
				if model.ETag == nil {
					return fmt.Errorf("add %s: etag was nil", *id)
				}
			}

			if model.Properties == nil {
				return fmt.Errorf("retreiving properties for %s for update: %+v", *id, err)
			}

			if metadata.ResourceData.HasChange("display_name") {
				model.Properties.DisplayName = utils.String(metadata.ResourceData.Get("display_name").(string))
			}

			if metadata.ResourceData.HasChange("chart_type") {
				model.Properties.Chart = pointer.To(views.ChartType(metadata.ResourceData.Get("chart_type").(string)))
			}

			if metadata.ResourceData.HasChange("dataset") {
				model.Properties.Query.DataSet = expandDataset(metadata.ResourceData.Get("dataset").([]interface{}))
			}
			if metadata.ResourceData.HasChange("timeframe") {
				model.Properties.Query.Timeframe = views.ReportTimeframeType(metadata.ResourceData.Get("timeframe").(string))
			}

			if metadata.ResourceData.HasChange("kpi") {
				model.Properties.Kpis = expandKpis(metadata.ResourceData.Get("kpi").([]interface{}))
			}

			if metadata.ResourceData.HasChange("pivot") {
				model.Properties.Pivots = expandPivots(metadata.ResourceData.Get("pivot").([]interface{}))
			}

			_, err = client.CreateOrUpdateByScope(ctx, *id, *model)
			if err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func expandDataset(input []interface{}) *views.ReportConfigDataset {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	attrs := input[0].(map[string]interface{})
	dataset := &views.ReportConfigDataset{
		Granularity: pointer.To(views.ReportGranularityType(attrs["granularity"].(string))),
	}

	if aggregation := attrs["aggregation"].(*pluginsdk.Set).List(); len(aggregation) > 0 {
		dataset.Aggregation = expandAggregation(aggregation)
	}

	if sorting := attrs["sorting"].([]interface{}); len(sorting) > 0 {
		dataset.Sorting = expandSorting(sorting)
	}

	if grouping := attrs["grouping"].([]interface{}); len(grouping) > 0 {
		dataset.Grouping = expandGrouping(grouping)
	}

	return dataset
}

func expandAggregation(input []interface{}) *map[string]views.ReportConfigAggregation {
	outputSorting := map[string]views.ReportConfigAggregation{}
	if len(input) == 0 || input[0] == nil {
		return &outputSorting
	}

	for _, item := range input {
		v := item.(map[string]interface{})
		name := v["name"].(string)
		outputSorting[name] = views.ReportConfigAggregation{
			Name:     v["column_name"].(string),
			Function: views.FunctionTypeSum,
		}
	}

	return &outputSorting
}

func expandGrouping(input []interface{}) *[]views.ReportConfigGrouping {
	outputGrouping := []views.ReportConfigGrouping{}
	if len(input) == 0 || input[0] == nil {
		return &outputGrouping
	}

	for _, item := range input {
		v := item.(map[string]interface{})
		outputGrouping = append(outputGrouping, views.ReportConfigGrouping{
			Type: views.QueryColumnType(v["type"].(string)),
			Name: v["name"].(string),
		})
	}

	return &outputGrouping
}

func expandSorting(input []interface{}) *[]views.ReportConfigSorting {
	outputSorting := []views.ReportConfigSorting{}
	if len(input) == 0 || input[0] == nil {
		return &outputSorting
	}

	for _, item := range input {
		v := item.(map[string]interface{})
		outputSorting = append(outputSorting, views.ReportConfigSorting{
			Direction: pointer.To(views.ReportConfigSortingType(v["direction"].(string))),
			Name:      v["name"].(string),
		})
	}

	return &outputSorting
}

func expandKpis(input []interface{}) *[]views.KpiProperties {
	outputKpis := []views.KpiProperties{}
	if len(input) == 0 || input[0] == nil {
		return &outputKpis
	}

	for _, item := range input {
		v := item.(map[string]interface{})
		outputKpis = append(outputKpis, views.KpiProperties{
			Type:    pointer.To(views.KpiTypeType(v["type"].(string))),
			Enabled: utils.Bool(true),
		})
	}

	return &outputKpis
}

func expandPivots(input []interface{}) *[]views.PivotProperties {
	outputPivots := []views.PivotProperties{}
	if len(input) == 0 || input[0] == nil {
		return &outputPivots
	}

	for _, item := range input {
		v := item.(map[string]interface{})
		outputPivots = append(outputPivots, views.PivotProperties{
			Type: pointer.To(views.PivotTypeType(v["type"].(string))),
			Name: utils.String((v["name"].(string))),
		})
	}

	return &outputPivots
}

func flattenKpis(input *[]views.KpiProperties) []interface{} {
	outputKpis := []interface{}{}
	if input == nil || len(*input) == 0 {
		return outputKpis
	}

	for _, item := range *input {
		kpiType := ""
		if v := item.Type; v != nil && item.Enabled != nil && *item.Enabled {
			kpiType = string(*v)
		}

		outputKpis = append(outputKpis, map[string]interface{}{
			"type": kpiType,
		})
	}

	return outputKpis
}

func flattenPivots(input *[]views.PivotProperties) []interface{} {
	outputPivots := []interface{}{}
	if input == nil || len(*input) == 0 {
		return outputPivots
	}

	for _, item := range *input {
		pivotType := ""
		if v := item.Type; v != nil {
			pivotType = string(*v)
		}

		name := ""
		if p := item.Name; p != nil {
			name = *p
		}

		outputPivots = append(outputPivots, map[string]interface{}{
			"name": name,
			"type": pivotType,
		})
	}

	return outputPivots
}

func flattenDataset(input *views.ReportConfigDataset) []interface{} {
	outputDataset := map[string]interface{}{
		"aggregation": flattenAggregation(input.Aggregation),
		"sorting":     flattenSorting(input.Sorting),
		"grouping":    flattenGrouping(input.Grouping),
	}

	if input.Granularity != nil {
		outputDataset["granularity"] = string(*input.Granularity)
	}

	return []interface{}{outputDataset}
}

func flattenAggregation(input *map[string]views.ReportConfigAggregation) []interface{} {
	outputAggregations := []interface{}{}
	if input == nil || len(*input) == 0 {
		return outputAggregations
	}

	for name, item := range *input {
		outputAggregations = append(outputAggregations, map[string]interface{}{
			"name":        name,
			"column_name": item.Name,
		})
	}

	return outputAggregations
}

func flattenGrouping(input *[]views.ReportConfigGrouping) []interface{} {
	outputGroupings := []interface{}{}
	if input == nil || len(*input) == 0 {
		return outputGroupings
	}

	for _, item := range *input {
		outputGroupings = append(outputGroupings, map[string]interface{}{
			"name": item.Name,
			"type": string(item.Type),
		})
	}

	return outputGroupings
}

func flattenSorting(input *[]views.ReportConfigSorting) []interface{} {
	outputSortings := []interface{}{}
	if input == nil || len(*input) == 0 {
		return outputSortings
	}

	for _, item := range *input {
		if item.Direction == nil {
			continue
		}
		outputSortings = append(outputSortings, map[string]interface{}{
			"name":      item.Name,
			"direction": string(*item.Direction),
		})
	}

	return outputSortings
}
