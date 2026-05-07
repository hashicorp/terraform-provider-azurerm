// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package costmanagement

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/costmanagement/2023-08-01/views"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

// Shared nested model structs for cost management view resources

type CostManagementViewAggregationModel struct {
	Name       string `tfschema:"name"`
	ColumnName string `tfschema:"column_name"`
}

type CostManagementViewSortingModel struct {
	Direction string `tfschema:"direction"`
	Name      string `tfschema:"name"`
}

type CostManagementViewGroupingModel struct {
	Type string `tfschema:"type"`
	Name string `tfschema:"name"`
}

type CostManagementViewDatasetModel struct {
	Granularity string                               `tfschema:"granularity"`
	Aggregation []CostManagementViewAggregationModel `tfschema:"aggregation"`
	Sorting     []CostManagementViewSortingModel     `tfschema:"sorting"`
	Grouping    []CostManagementViewGroupingModel    `tfschema:"grouping"`
}

type CostManagementViewKpiModel struct {
	Type string `tfschema:"type"`
}

type CostManagementViewPivotModel struct {
	Name string `tfschema:"name"`
	Type string `tfschema:"type"`
}

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

// Typed model expand/flatten helpers

func expandDatasetFromModel(input []CostManagementViewDatasetModel) *views.ReportConfigDataset {
	if len(input) == 0 {
		return nil
	}

	ds := input[0]
	dataset := &views.ReportConfigDataset{
		Granularity: pointer.To(views.ReportGranularityType(ds.Granularity)),
	}

	if len(ds.Aggregation) > 0 {
		aggregation := map[string]views.ReportConfigAggregation{}
		for _, a := range ds.Aggregation {
			aggregation[a.Name] = views.ReportConfigAggregation{
				Name:     a.ColumnName,
				Function: views.FunctionTypeSum,
			}
		}
		dataset.Aggregation = &aggregation
	}

	if len(ds.Sorting) > 0 {
		sorting := make([]views.ReportConfigSorting, 0)
		for _, s := range ds.Sorting {
			sorting = append(sorting, views.ReportConfigSorting{
				Direction: pointer.To(views.ReportConfigSortingType(s.Direction)),
				Name:      s.Name,
			})
		}
		dataset.Sorting = &sorting
	}

	if len(ds.Grouping) > 0 {
		grouping := make([]views.ReportConfigGrouping, 0)
		for _, g := range ds.Grouping {
			grouping = append(grouping, views.ReportConfigGrouping{
				Type: views.QueryColumnType(g.Type),
				Name: g.Name,
			})
		}
		dataset.Grouping = &grouping
	}

	return dataset
}

func flattenDatasetToModel(input *views.ReportConfigDataset) []CostManagementViewDatasetModel {
	if input == nil {
		return []CostManagementViewDatasetModel{}
	}

	ds := CostManagementViewDatasetModel{}

	if input.Granularity != nil {
		ds.Granularity = string(*input.Granularity)
	}

	if input.Aggregation != nil {
		aggregation := make([]CostManagementViewAggregationModel, 0)
		for name, item := range *input.Aggregation {
			aggregation = append(aggregation, CostManagementViewAggregationModel{
				Name:       name,
				ColumnName: item.Name,
			})
		}
		ds.Aggregation = aggregation
	}

	if input.Sorting != nil {
		sorting := make([]CostManagementViewSortingModel, 0)
		for _, item := range *input.Sorting {
			if item.Direction == nil {
				continue
			}
			sorting = append(sorting, CostManagementViewSortingModel{
				Name:      item.Name,
				Direction: string(*item.Direction),
			})
		}
		ds.Sorting = sorting
	}

	if input.Grouping != nil {
		grouping := make([]CostManagementViewGroupingModel, 0)
		for _, item := range *input.Grouping {
			grouping = append(grouping, CostManagementViewGroupingModel{
				Name: item.Name,
				Type: string(item.Type),
			})
		}
		ds.Grouping = grouping
	}

	return []CostManagementViewDatasetModel{ds}
}

func expandKpisFromModel(input []CostManagementViewKpiModel) *[]views.KpiProperties {
	kpis := make([]views.KpiProperties, 0)
	for _, k := range input {
		kpis = append(kpis, views.KpiProperties{
			Type:    pointer.To(views.KpiTypeType(k.Type)),
			Enabled: pointer.To(true),
		})
	}
	return &kpis
}

func flattenKpisToModel(input *[]views.KpiProperties) []CostManagementViewKpiModel {
	if input == nil || len(*input) == 0 {
		return []CostManagementViewKpiModel{}
	}

	result := make([]CostManagementViewKpiModel, 0)
	for _, item := range *input {
		kpiType := ""
		if v := item.Type; v != nil && item.Enabled != nil && *item.Enabled {
			kpiType = string(*v)
		}
		result = append(result, CostManagementViewKpiModel{
			Type: kpiType,
		})
	}
	return result
}

func expandPivotsFromModel(input []CostManagementViewPivotModel) *[]views.PivotProperties {
	pivots := make([]views.PivotProperties, 0)
	for _, p := range input {
		pivots = append(pivots, views.PivotProperties{
			Type: pointer.To(views.PivotTypeType(p.Type)),
			Name: pointer.To(p.Name),
		})
	}
	return &pivots
}

func flattenPivotsToModel(input *[]views.PivotProperties) []CostManagementViewPivotModel {
	if input == nil || len(*input) == 0 {
		return []CostManagementViewPivotModel{}
	}

	result := make([]CostManagementViewPivotModel, 0)
	for _, item := range *input {
		pivotType := ""
		if v := item.Type; v != nil {
			pivotType = string(*v)
		}
		name := ""
		if p := item.Name; p != nil {
			name = *p
		}
		result = append(result, CostManagementViewPivotModel{
			Name: name,
			Type: pivotType,
		})
	}
	return result
}
