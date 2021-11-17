package query

import "strings"

type ExportType string

const (
	ExportTypeActualCost    ExportType = "ActualCost"
	ExportTypeAmortizedCost ExportType = "AmortizedCost"
	ExportTypeUsage         ExportType = "Usage"
)

func PossibleValuesForExportType() []string {
	return []string{
		string(ExportTypeActualCost),
		string(ExportTypeAmortizedCost),
		string(ExportTypeUsage),
	}
}

func parseExportType(input string) (*ExportType, error) {
	vals := map[string]ExportType{
		"actualcost":    ExportTypeActualCost,
		"amortizedcost": ExportTypeAmortizedCost,
		"usage":         ExportTypeUsage,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ExportType(input)
	return &out, nil
}

type ExternalCloudProviderType string

const (
	ExternalCloudProviderTypeExternalBillingAccounts ExternalCloudProviderType = "externalBillingAccounts"
	ExternalCloudProviderTypeExternalSubscriptions   ExternalCloudProviderType = "externalSubscriptions"
)

func PossibleValuesForExternalCloudProviderType() []string {
	return []string{
		string(ExternalCloudProviderTypeExternalBillingAccounts),
		string(ExternalCloudProviderTypeExternalSubscriptions),
	}
}

func parseExternalCloudProviderType(input string) (*ExternalCloudProviderType, error) {
	vals := map[string]ExternalCloudProviderType{
		"externalbillingaccounts": ExternalCloudProviderTypeExternalBillingAccounts,
		"externalsubscriptions":   ExternalCloudProviderTypeExternalSubscriptions,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ExternalCloudProviderType(input)
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

type GranularityType string

const (
	GranularityTypeDaily GranularityType = "Daily"
)

func PossibleValuesForGranularityType() []string {
	return []string{
		string(GranularityTypeDaily),
	}
}

func parseGranularityType(input string) (*GranularityType, error) {
	vals := map[string]GranularityType{
		"daily": GranularityTypeDaily,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := GranularityType(input)
	return &out, nil
}

type QueryColumnType string

const (
	QueryColumnTypeDimension QueryColumnType = "Dimension"
	QueryColumnTypeTag       QueryColumnType = "Tag"
)

func PossibleValuesForQueryColumnType() []string {
	return []string{
		string(QueryColumnTypeDimension),
		string(QueryColumnTypeTag),
	}
}

func parseQueryColumnType(input string) (*QueryColumnType, error) {
	vals := map[string]QueryColumnType{
		"dimension": QueryColumnTypeDimension,
		"tag":       QueryColumnTypeTag,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := QueryColumnType(input)
	return &out, nil
}

type QueryOperatorType string

const (
	QueryOperatorTypeIn QueryOperatorType = "In"
)

func PossibleValuesForQueryOperatorType() []string {
	return []string{
		string(QueryOperatorTypeIn),
	}
}

func parseQueryOperatorType(input string) (*QueryOperatorType, error) {
	vals := map[string]QueryOperatorType{
		"in": QueryOperatorTypeIn,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := QueryOperatorType(input)
	return &out, nil
}

type TimeframeType string

const (
	TimeframeTypeBillingMonthToDate  TimeframeType = "BillingMonthToDate"
	TimeframeTypeCustom              TimeframeType = "Custom"
	TimeframeTypeMonthToDate         TimeframeType = "MonthToDate"
	TimeframeTypeTheLastBillingMonth TimeframeType = "TheLastBillingMonth"
	TimeframeTypeTheLastMonth        TimeframeType = "TheLastMonth"
	TimeframeTypeWeekToDate          TimeframeType = "WeekToDate"
)

func PossibleValuesForTimeframeType() []string {
	return []string{
		string(TimeframeTypeBillingMonthToDate),
		string(TimeframeTypeCustom),
		string(TimeframeTypeMonthToDate),
		string(TimeframeTypeTheLastBillingMonth),
		string(TimeframeTypeTheLastMonth),
		string(TimeframeTypeWeekToDate),
	}
}

func parseTimeframeType(input string) (*TimeframeType, error) {
	vals := map[string]TimeframeType{
		"billingmonthtodate":  TimeframeTypeBillingMonthToDate,
		"custom":              TimeframeTypeCustom,
		"monthtodate":         TimeframeTypeMonthToDate,
		"thelastbillingmonth": TimeframeTypeTheLastBillingMonth,
		"thelastmonth":        TimeframeTypeTheLastMonth,
		"weektodate":          TimeframeTypeWeekToDate,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := TimeframeType(input)
	return &out, nil
}
