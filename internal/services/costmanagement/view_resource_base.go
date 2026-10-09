// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package costmanagement

import (
	"context"
	"fmt"
	"maps"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/costmanagement/2025-03-01/viewoperationgroup"
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
			ValidateFunc: validation.StringInSlice(viewoperationgroup.PossibleValuesForChartType(), false),
		},

		"accumulated": {
			Type:     pluginsdk.TypeBool,
			Required: true,
			ForceNew: true,
		},

		"report_type": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(viewoperationgroup.PossibleValuesForReportType(), false),
		},

		"timeframe": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(viewoperationgroup.PossibleValuesForReportTimeframeType(), false),
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
						ValidateFunc: validation.StringInSlice(viewoperationgroup.PossibleValuesForReportGranularityType(), false),
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
									ValidateFunc: validation.StringInSlice(viewoperationgroup.PossibleValuesForReportConfigSortingType(), false),
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
									ValidateFunc: validation.StringInSlice(viewoperationgroup.PossibleValuesForQueryColumnType(), false),
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
						ValidateFunc: validation.StringInSlice(viewoperationgroup.PossibleValuesForKpiTypeType(), false),
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
						ValidateFunc: validation.StringInSlice(viewoperationgroup.PossibleValuesForPivotTypeType(), false),
					},
				},
			},
		},
	}

	maps.Copy(output, fields)

	return output
}

func (br costManagementViewBaseResource) attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (br costManagementViewBaseResource) deleteFunc() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.CostManagement.ViewOperationGroupClient

			id, err := viewoperationgroup.ParseScopedViewID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err = client.ViewsDeleteByScope(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

// Typed model expand/flatten helpers

func expandDatasetFromModel(input []CostManagementViewDatasetModel) *viewoperationgroup.ReportConfigDataset {
	if len(input) == 0 {
		return nil
	}

	ds := input[0]
	dataset := &viewoperationgroup.ReportConfigDataset{
		Granularity: pointer.ToEnum[viewoperationgroup.ReportGranularityType](ds.Granularity),
	}

	if len(ds.Aggregation) > 0 {
		aggregation := map[string]viewoperationgroup.ReportConfigAggregation{}
		for _, a := range ds.Aggregation {
			aggregation[a.Name] = viewoperationgroup.ReportConfigAggregation{
				Name:     a.ColumnName,
				Function: viewoperationgroup.FunctionTypeSum,
			}
		}
		dataset.Aggregation = &aggregation
	}

	if len(ds.Sorting) > 0 {
		sorting := make([]viewoperationgroup.ReportConfigSorting, 0)
		for _, s := range ds.Sorting {
			sorting = append(sorting, viewoperationgroup.ReportConfigSorting{
				Direction: pointer.ToEnum[viewoperationgroup.ReportConfigSortingType](s.Direction),
				Name:      s.Name,
			})
		}
		dataset.Sorting = &sorting
	}

	if len(ds.Grouping) > 0 {
		grouping := make([]viewoperationgroup.ReportConfigGrouping, 0)
		for _, g := range ds.Grouping {
			grouping = append(grouping, viewoperationgroup.ReportConfigGrouping{
				Type: viewoperationgroup.QueryColumnType(g.Type),
				Name: g.Name,
			})
		}
		dataset.Grouping = &grouping
	}

	return dataset
}

func flattenDatasetToModel(input *viewoperationgroup.ReportConfigDataset) []CostManagementViewDatasetModel {
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

func expandKpisFromModel(input []CostManagementViewKpiModel) *[]viewoperationgroup.KpiProperties {
	kpis := make([]viewoperationgroup.KpiProperties, 0)
	for _, k := range input {
		kpis = append(kpis, viewoperationgroup.KpiProperties{
			Type:    pointer.ToEnum[viewoperationgroup.KpiTypeType](k.Type),
			Enabled: pointer.To(true),
		})
	}
	return &kpis
}

func flattenKpisToModel(input *[]viewoperationgroup.KpiProperties) []CostManagementViewKpiModel {
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

func expandPivotsFromModel(input []CostManagementViewPivotModel) *[]viewoperationgroup.PivotProperties {
	pivots := make([]viewoperationgroup.PivotProperties, 0)
	for _, p := range input {
		pivots = append(pivots, viewoperationgroup.PivotProperties{
			Type: pointer.ToEnum[viewoperationgroup.PivotTypeType](p.Type),
			Name: pointer.To(p.Name),
		})
	}
	return &pivots
}

func flattenPivotsToModel(input *[]viewoperationgroup.PivotProperties) []CostManagementViewPivotModel {
	if input == nil || len(*input) == 0 {
		return []CostManagementViewPivotModel{}
	}

	result := make([]CostManagementViewPivotModel, 0)
	for _, item := range *input {
		pivotType := ""
		if v := item.Type; v != nil {
			pivotType = string(*v)
		}
		name := pointer.From(item.Name)
		result = append(result, CostManagementViewPivotModel{
			Name: name,
			Type: pivotType,
		})
	}
	return result
}
