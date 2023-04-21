package scheduledqueryrules

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AlertSeverity int64

const (
	AlertSeverityFour  AlertSeverity = 4
	AlertSeverityOne   AlertSeverity = 1
	AlertSeverityThree AlertSeverity = 3
	AlertSeverityTwo   AlertSeverity = 2
	AlertSeverityZero  AlertSeverity = 0
)

func PossibleValuesForAlertSeverity() []int64 {
	return []int64{
		int64(AlertSeverityFour),
		int64(AlertSeverityOne),
		int64(AlertSeverityThree),
		int64(AlertSeverityTwo),
		int64(AlertSeverityZero),
	}
}

type ConditionOperator string

const (
	ConditionOperatorEquals             ConditionOperator = "Equals"
	ConditionOperatorGreaterThan        ConditionOperator = "GreaterThan"
	ConditionOperatorGreaterThanOrEqual ConditionOperator = "GreaterThanOrEqual"
	ConditionOperatorLessThan           ConditionOperator = "LessThan"
	ConditionOperatorLessThanOrEqual    ConditionOperator = "LessThanOrEqual"
)

func PossibleValuesForConditionOperator() []string {
	return []string{
		string(ConditionOperatorEquals),
		string(ConditionOperatorGreaterThan),
		string(ConditionOperatorGreaterThanOrEqual),
		string(ConditionOperatorLessThan),
		string(ConditionOperatorLessThanOrEqual),
	}
}

func parseConditionOperator(input string) (*ConditionOperator, error) {
	vals := map[string]ConditionOperator{
		"equals":             ConditionOperatorEquals,
		"greaterthan":        ConditionOperatorGreaterThan,
		"greaterthanorequal": ConditionOperatorGreaterThanOrEqual,
		"lessthan":           ConditionOperatorLessThan,
		"lessthanorequal":    ConditionOperatorLessThanOrEqual,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ConditionOperator(input)
	return &out, nil
}

type DimensionOperator string

const (
	DimensionOperatorExclude DimensionOperator = "Exclude"
	DimensionOperatorInclude DimensionOperator = "Include"
)

func PossibleValuesForDimensionOperator() []string {
	return []string{
		string(DimensionOperatorExclude),
		string(DimensionOperatorInclude),
	}
}

func parseDimensionOperator(input string) (*DimensionOperator, error) {
	vals := map[string]DimensionOperator{
		"exclude": DimensionOperatorExclude,
		"include": DimensionOperatorInclude,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DimensionOperator(input)
	return &out, nil
}

type Kind string

const (
	KindLogAlert    Kind = "LogAlert"
	KindLogToMetric Kind = "LogToMetric"
)

func PossibleValuesForKind() []string {
	return []string{
		string(KindLogAlert),
		string(KindLogToMetric),
	}
}

func parseKind(input string) (*Kind, error) {
	vals := map[string]Kind{
		"logalert":    KindLogAlert,
		"logtometric": KindLogToMetric,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Kind(input)
	return &out, nil
}

type TimeAggregation string

const (
	TimeAggregationAverage TimeAggregation = "Average"
	TimeAggregationCount   TimeAggregation = "Count"
	TimeAggregationMaximum TimeAggregation = "Maximum"
	TimeAggregationMinimum TimeAggregation = "Minimum"
	TimeAggregationTotal   TimeAggregation = "Total"
)

func PossibleValuesForTimeAggregation() []string {
	return []string{
		string(TimeAggregationAverage),
		string(TimeAggregationCount),
		string(TimeAggregationMaximum),
		string(TimeAggregationMinimum),
		string(TimeAggregationTotal),
	}
}

func parseTimeAggregation(input string) (*TimeAggregation, error) {
	vals := map[string]TimeAggregation{
		"average": TimeAggregationAverage,
		"count":   TimeAggregationCount,
		"maximum": TimeAggregationMaximum,
		"minimum": TimeAggregationMinimum,
		"total":   TimeAggregationTotal,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := TimeAggregation(input)
	return &out, nil
}
