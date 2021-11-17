package forecast

import "strings"

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

type ForecastTimeframeType string

const (
	ForecastTimeframeTypeBillingMonthToDate  ForecastTimeframeType = "BillingMonthToDate"
	ForecastTimeframeTypeCustom              ForecastTimeframeType = "Custom"
	ForecastTimeframeTypeMonthToDate         ForecastTimeframeType = "MonthToDate"
	ForecastTimeframeTypeTheLastBillingMonth ForecastTimeframeType = "TheLastBillingMonth"
	ForecastTimeframeTypeTheLastMonth        ForecastTimeframeType = "TheLastMonth"
	ForecastTimeframeTypeWeekToDate          ForecastTimeframeType = "WeekToDate"
)

func PossibleValuesForForecastTimeframeType() []string {
	return []string{
		string(ForecastTimeframeTypeBillingMonthToDate),
		string(ForecastTimeframeTypeCustom),
		string(ForecastTimeframeTypeMonthToDate),
		string(ForecastTimeframeTypeTheLastBillingMonth),
		string(ForecastTimeframeTypeTheLastMonth),
		string(ForecastTimeframeTypeWeekToDate),
	}
}

func parseForecastTimeframeType(input string) (*ForecastTimeframeType, error) {
	vals := map[string]ForecastTimeframeType{
		"billingmonthtodate":  ForecastTimeframeTypeBillingMonthToDate,
		"custom":              ForecastTimeframeTypeCustom,
		"monthtodate":         ForecastTimeframeTypeMonthToDate,
		"thelastbillingmonth": ForecastTimeframeTypeTheLastBillingMonth,
		"thelastmonth":        ForecastTimeframeTypeTheLastMonth,
		"weektodate":          ForecastTimeframeTypeWeekToDate,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ForecastTimeframeType(input)
	return &out, nil
}

type ForecastType string

const (
	ForecastTypeActualCost    ForecastType = "ActualCost"
	ForecastTypeAmortizedCost ForecastType = "AmortizedCost"
	ForecastTypeUsage         ForecastType = "Usage"
)

func PossibleValuesForForecastType() []string {
	return []string{
		string(ForecastTypeActualCost),
		string(ForecastTypeAmortizedCost),
		string(ForecastTypeUsage),
	}
}

func parseForecastType(input string) (*ForecastType, error) {
	vals := map[string]ForecastType{
		"actualcost":    ForecastTypeActualCost,
		"amortizedcost": ForecastTypeAmortizedCost,
		"usage":         ForecastTypeUsage,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ForecastType(input)
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
