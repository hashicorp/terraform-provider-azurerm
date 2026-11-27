// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"context"
	"fmt"
	"maps"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/costmanagement/2025-03-01/viewoperationgroup"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	costmanagementClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/costmanagement/client"
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

type CostManagementViewBaseResource struct{}

func (br CostManagementViewBaseResource) Arguments(fields map[string]*pluginsdk.Schema) map[string]*pluginsdk.Schema {
	output := map[string]*pluginsdk.Schema{
		"display_name": {
			Description: "User visible input name of the Cost Management View.",
			Type:        pluginsdk.TypeString,
			Required:    true,
		},

		"chart_type": {
			Description:  "Chart type of the main view in Cost Analysis.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(viewoperationgroup.PossibleValuesForChartType(), false),
		},

		"accumulated": {
			Description: "Whether the costs data in the Cost Management View are accumulated over time.",
			Type:        pluginsdk.TypeBool,
			Required:    true,
			ForceNew:    true,
		},

		"report_type": {
			Description:  "The type of the report.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(viewoperationgroup.PossibleValuesForReportType(), false),
		},

		"timeframe": {
			Description:  "The time frame for pulling data for the report.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(viewoperationgroup.PossibleValuesForReportTimeframeType(), false),
		},

		"dataset": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"granularity": {
						Description:  "The granularity of the rows in the report.",
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
									Description: "The name which should be used for this aggregation.",
									Type:        pluginsdk.TypeString,
									Required:    true,
									ForceNew:    true,
								},
								"column_name": {
									Description: "The name of the column to aggregate.",
									Type:        pluginsdk.TypeString,
									Required:    true,
									ForceNew:    true,
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
									Description:  "Direction of sort.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringInSlice(viewoperationgroup.PossibleValuesForReportConfigSortingType(), false),
								},
								"name": {
									Description: "The name of the column to sort.",
									Type:        pluginsdk.TypeString,
									Required:    true,
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
									Description:  "The type of column.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringInSlice(viewoperationgroup.PossibleValuesForQueryColumnType(), false),
								},
								"name": {
									Description: "The name of the column to group.",
									Type:        pluginsdk.TypeString,
									Required:    true,
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
						Description:  "KPI type.",
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
						Description: "The name of the column which should be used for this sub-view in the Cost Analysis UI.",
						Type:        pluginsdk.TypeString,
						Required:    true,
					},
					"type": {
						Description:  "The data type to show in this sub-view.",
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

func (br CostManagementViewBaseResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (br CostManagementViewBaseResource) DeleteFunc() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := costmanagementClient.From(metadata.Client).ViewOperationGroupClient

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

func ExpandDatasetFromModel(input []CostManagementViewDatasetModel) *viewoperationgroup.ReportConfigDataset {
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

func FlattenDatasetToModel(input *viewoperationgroup.ReportConfigDataset) []CostManagementViewDatasetModel {
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

func ExpandKpisFromModel(input []CostManagementViewKpiModel) *[]viewoperationgroup.KpiProperties {
	kpis := make([]viewoperationgroup.KpiProperties, 0)
	for _, k := range input {
		kpis = append(kpis, viewoperationgroup.KpiProperties{
			Type:    pointer.ToEnum[viewoperationgroup.KpiTypeType](k.Type),
			Enabled: pointer.To(true),
		})
	}
	return &kpis
}

func FlattenKpisToModel(input *[]viewoperationgroup.KpiProperties) []CostManagementViewKpiModel {
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

func ExpandPivotsFromModel(input []CostManagementViewPivotModel) *[]viewoperationgroup.PivotProperties {
	pivots := make([]viewoperationgroup.PivotProperties, 0)
	for _, p := range input {
		pivots = append(pivots, viewoperationgroup.PivotProperties{
			Type: pointer.ToEnum[viewoperationgroup.PivotTypeType](p.Type),
			Name: pointer.To(p.Name),
		})
	}
	return &pivots
}

func FlattenPivotsToModel(input *[]viewoperationgroup.PivotProperties) []CostManagementViewPivotModel {
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
