package usagedetails

import "strings"

type GenerateDetailedCostReportMetricType string

const (
	GenerateDetailedCostReportMetricTypeActualCost    GenerateDetailedCostReportMetricType = "ActualCost"
	GenerateDetailedCostReportMetricTypeAmortizedCost GenerateDetailedCostReportMetricType = "AmortizedCost"
)

func PossibleValuesForGenerateDetailedCostReportMetricType() []string {
	return []string{
		string(GenerateDetailedCostReportMetricTypeActualCost),
		string(GenerateDetailedCostReportMetricTypeAmortizedCost),
	}
}

func parseGenerateDetailedCostReportMetricType(input string) (*GenerateDetailedCostReportMetricType, error) {
	vals := map[string]GenerateDetailedCostReportMetricType{
		"actualcost":    GenerateDetailedCostReportMetricTypeActualCost,
		"amortizedcost": GenerateDetailedCostReportMetricTypeAmortizedCost,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := GenerateDetailedCostReportMetricType(input)
	return &out, nil
}
