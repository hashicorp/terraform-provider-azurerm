package views

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AccumulatedType string

const (
	AccumulatedTypeFalse AccumulatedType = "false"
	AccumulatedTypeTrue  AccumulatedType = "true"
)

func PossibleValuesForAccumulatedType() []string {
	return []string{
		string(AccumulatedTypeFalse),
		string(AccumulatedTypeTrue),
	}
}

func parseAccumulatedType(input string) (*AccumulatedType, error) {
	vals := map[string]AccumulatedType{
		"false": AccumulatedTypeFalse,
		"true":  AccumulatedTypeTrue,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AccumulatedType(input)
	return &out, nil
}

type ChartType string

const (
	ChartTypeArea          ChartType = "Area"
	ChartTypeGroupedColumn ChartType = "GroupedColumn"
	ChartTypeLine          ChartType = "Line"
	ChartTypeStackedColumn ChartType = "StackedColumn"
	ChartTypeTable         ChartType = "Table"
)

func PossibleValuesForChartType() []string {
	return []string{
		string(ChartTypeArea),
		string(ChartTypeGroupedColumn),
		string(ChartTypeLine),
		string(ChartTypeStackedColumn),
		string(ChartTypeTable),
	}
}

func parseChartType(input string) (*ChartType, error) {
	vals := map[string]ChartType{
		"area":          ChartTypeArea,
		"groupedcolumn": ChartTypeGroupedColumn,
		"line":          ChartTypeLine,
		"stackedcolumn": ChartTypeStackedColumn,
		"table":         ChartTypeTable,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ChartType(input)
	return &out, nil
}

type FunctionType string

const (
	FunctionTypeSum FunctionType = "Sum"
)

func PossibleValuesForFunctionType() []string {
	return []string{
		string(FunctionTypeSum),
	}
}

func parseFunctionType(input string) (*FunctionType, error) {
	vals := map[string]FunctionType{
		"sum": FunctionTypeSum,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := FunctionType(input)
	return &out, nil
}

type KpiTypeType string

const (
	KpiTypeTypeBudget   KpiTypeType = "Budget"
	KpiTypeTypeForecast KpiTypeType = "Forecast"
)

func PossibleValuesForKpiTypeType() []string {
	return []string{
		string(KpiTypeTypeBudget),
		string(KpiTypeTypeForecast),
	}
}

func parseKpiTypeType(input string) (*KpiTypeType, error) {
	vals := map[string]KpiTypeType{
		"budget":   KpiTypeTypeBudget,
		"forecast": KpiTypeTypeForecast,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := KpiTypeType(input)
	return &out, nil
}

type MetricType string

const (
	MetricTypeAHUB          MetricType = "AHUB"
	MetricTypeActualCost    MetricType = "ActualCost"
	MetricTypeAmortizedCost MetricType = "AmortizedCost"
)

func PossibleValuesForMetricType() []string {
	return []string{
		string(MetricTypeAHUB),
		string(MetricTypeActualCost),
		string(MetricTypeAmortizedCost),
	}
}

func parseMetricType(input string) (*MetricType, error) {
	vals := map[string]MetricType{
		"ahub":          MetricTypeAHUB,
		"actualcost":    MetricTypeActualCost,
		"amortizedcost": MetricTypeAmortizedCost,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := MetricType(input)
	return &out, nil
}

type OperatorType string

const (
	OperatorTypeContains OperatorType = "Contains"
	OperatorTypeIn       OperatorType = "In"
)

func PossibleValuesForOperatorType() []string {
	return []string{
		string(OperatorTypeContains),
		string(OperatorTypeIn),
	}
}

func parseOperatorType(input string) (*OperatorType, error) {
	vals := map[string]OperatorType{
		"contains": OperatorTypeContains,
		"in":       OperatorTypeIn,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := OperatorType(input)
	return &out, nil
}

type PivotTypeType string

const (
	PivotTypeTypeDimension PivotTypeType = "Dimension"
	PivotTypeTypeTagKey    PivotTypeType = "TagKey"
)

func PossibleValuesForPivotTypeType() []string {
	return []string{
		string(PivotTypeTypeDimension),
		string(PivotTypeTypeTagKey),
	}
}

func parsePivotTypeType(input string) (*PivotTypeType, error) {
	vals := map[string]PivotTypeType{
		"dimension": PivotTypeTypeDimension,
		"tagkey":    PivotTypeTypeTagKey,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PivotTypeType(input)
	return &out, nil
}

type QueryColumnType string

const (
	QueryColumnTypeDimension QueryColumnType = "Dimension"
	QueryColumnTypeTagKey    QueryColumnType = "TagKey"
)

func PossibleValuesForQueryColumnType() []string {
	return []string{
		string(QueryColumnTypeDimension),
		string(QueryColumnTypeTagKey),
	}
}

func parseQueryColumnType(input string) (*QueryColumnType, error) {
	vals := map[string]QueryColumnType{
		"dimension": QueryColumnTypeDimension,
		"tagkey":    QueryColumnTypeTagKey,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := QueryColumnType(input)
	return &out, nil
}

type ReportConfigSortingType string

const (
	ReportConfigSortingTypeAscending  ReportConfigSortingType = "Ascending"
	ReportConfigSortingTypeDescending ReportConfigSortingType = "Descending"
)

func PossibleValuesForReportConfigSortingType() []string {
	return []string{
		string(ReportConfigSortingTypeAscending),
		string(ReportConfigSortingTypeDescending),
	}
}

func parseReportConfigSortingType(input string) (*ReportConfigSortingType, error) {
	vals := map[string]ReportConfigSortingType{
		"ascending":  ReportConfigSortingTypeAscending,
		"descending": ReportConfigSortingTypeDescending,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ReportConfigSortingType(input)
	return &out, nil
}

type ReportGranularityType string

const (
	ReportGranularityTypeDaily   ReportGranularityType = "Daily"
	ReportGranularityTypeMonthly ReportGranularityType = "Monthly"
)

func PossibleValuesForReportGranularityType() []string {
	return []string{
		string(ReportGranularityTypeDaily),
		string(ReportGranularityTypeMonthly),
	}
}

func parseReportGranularityType(input string) (*ReportGranularityType, error) {
	vals := map[string]ReportGranularityType{
		"daily":   ReportGranularityTypeDaily,
		"monthly": ReportGranularityTypeMonthly,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ReportGranularityType(input)
	return &out, nil
}

type ReportTimeframeType string

const (
	ReportTimeframeTypeCustom      ReportTimeframeType = "Custom"
	ReportTimeframeTypeMonthToDate ReportTimeframeType = "MonthToDate"
	ReportTimeframeTypeWeekToDate  ReportTimeframeType = "WeekToDate"
	ReportTimeframeTypeYearToDate  ReportTimeframeType = "YearToDate"
)

func PossibleValuesForReportTimeframeType() []string {
	return []string{
		string(ReportTimeframeTypeCustom),
		string(ReportTimeframeTypeMonthToDate),
		string(ReportTimeframeTypeWeekToDate),
		string(ReportTimeframeTypeYearToDate),
	}
}

func parseReportTimeframeType(input string) (*ReportTimeframeType, error) {
	vals := map[string]ReportTimeframeType{
		"custom":      ReportTimeframeTypeCustom,
		"monthtodate": ReportTimeframeTypeMonthToDate,
		"weektodate":  ReportTimeframeTypeWeekToDate,
		"yeartodate":  ReportTimeframeTypeYearToDate,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ReportTimeframeType(input)
	return &out, nil
}

type ReportType string

const (
	ReportTypeUsage ReportType = "Usage"
)

func PossibleValuesForReportType() []string {
	return []string{
		string(ReportTypeUsage),
	}
}

func parseReportType(input string) (*ReportType, error) {
	vals := map[string]ReportType{
		"usage": ReportTypeUsage,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ReportType(input)
	return &out, nil
}
