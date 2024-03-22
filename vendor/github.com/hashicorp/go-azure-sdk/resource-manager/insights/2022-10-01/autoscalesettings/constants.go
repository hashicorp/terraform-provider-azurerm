package autoscalesettings

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ComparisonOperationType string

const (
	ComparisonOperationTypeEquals             ComparisonOperationType = "Equals"
	ComparisonOperationTypeGreaterThan        ComparisonOperationType = "GreaterThan"
	ComparisonOperationTypeGreaterThanOrEqual ComparisonOperationType = "GreaterThanOrEqual"
	ComparisonOperationTypeLessThan           ComparisonOperationType = "LessThan"
	ComparisonOperationTypeLessThanOrEqual    ComparisonOperationType = "LessThanOrEqual"
	ComparisonOperationTypeNotEquals          ComparisonOperationType = "NotEquals"
)

func PossibleValuesForComparisonOperationType() []string {
	return []string{
		string(ComparisonOperationTypeEquals),
		string(ComparisonOperationTypeGreaterThan),
		string(ComparisonOperationTypeGreaterThanOrEqual),
		string(ComparisonOperationTypeLessThan),
		string(ComparisonOperationTypeLessThanOrEqual),
		string(ComparisonOperationTypeNotEquals),
	}
}

func (s *ComparisonOperationType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseComparisonOperationType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseComparisonOperationType(input string) (*ComparisonOperationType, error) {
	vals := map[string]ComparisonOperationType{
		"equals":             ComparisonOperationTypeEquals,
		"greaterthan":        ComparisonOperationTypeGreaterThan,
		"greaterthanorequal": ComparisonOperationTypeGreaterThanOrEqual,
		"lessthan":           ComparisonOperationTypeLessThan,
		"lessthanorequal":    ComparisonOperationTypeLessThanOrEqual,
		"notequals":          ComparisonOperationTypeNotEquals,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ComparisonOperationType(input)
	return &out, nil
}

type MetricStatisticType string

const (
	MetricStatisticTypeAverage MetricStatisticType = "Average"
	MetricStatisticTypeCount   MetricStatisticType = "Count"
	MetricStatisticTypeMax     MetricStatisticType = "Max"
	MetricStatisticTypeMin     MetricStatisticType = "Min"
	MetricStatisticTypeSum     MetricStatisticType = "Sum"
)

func PossibleValuesForMetricStatisticType() []string {
	return []string{
		string(MetricStatisticTypeAverage),
		string(MetricStatisticTypeCount),
		string(MetricStatisticTypeMax),
		string(MetricStatisticTypeMin),
		string(MetricStatisticTypeSum),
	}
}

func (s *MetricStatisticType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseMetricStatisticType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseMetricStatisticType(input string) (*MetricStatisticType, error) {
	vals := map[string]MetricStatisticType{
		"average": MetricStatisticTypeAverage,
		"count":   MetricStatisticTypeCount,
		"max":     MetricStatisticTypeMax,
		"min":     MetricStatisticTypeMin,
		"sum":     MetricStatisticTypeSum,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := MetricStatisticType(input)
	return &out, nil
}

type OperationType string

const (
	OperationTypeScale OperationType = "Scale"
)

func PossibleValuesForOperationType() []string {
	return []string{
		string(OperationTypeScale),
	}
}

func (s *OperationType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseOperationType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseOperationType(input string) (*OperationType, error) {
	vals := map[string]OperationType{
		"scale": OperationTypeScale,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := OperationType(input)
	return &out, nil
}

type PredictiveAutoscalePolicyScaleMode string

const (
	PredictiveAutoscalePolicyScaleModeDisabled     PredictiveAutoscalePolicyScaleMode = "Disabled"
	PredictiveAutoscalePolicyScaleModeEnabled      PredictiveAutoscalePolicyScaleMode = "Enabled"
	PredictiveAutoscalePolicyScaleModeForecastOnly PredictiveAutoscalePolicyScaleMode = "ForecastOnly"
)

func PossibleValuesForPredictiveAutoscalePolicyScaleMode() []string {
	return []string{
		string(PredictiveAutoscalePolicyScaleModeDisabled),
		string(PredictiveAutoscalePolicyScaleModeEnabled),
		string(PredictiveAutoscalePolicyScaleModeForecastOnly),
	}
}

func (s *PredictiveAutoscalePolicyScaleMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePredictiveAutoscalePolicyScaleMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePredictiveAutoscalePolicyScaleMode(input string) (*PredictiveAutoscalePolicyScaleMode, error) {
	vals := map[string]PredictiveAutoscalePolicyScaleMode{
		"disabled":     PredictiveAutoscalePolicyScaleModeDisabled,
		"enabled":      PredictiveAutoscalePolicyScaleModeEnabled,
		"forecastonly": PredictiveAutoscalePolicyScaleModeForecastOnly,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PredictiveAutoscalePolicyScaleMode(input)
	return &out, nil
}

type RecurrenceFrequency string

const (
	RecurrenceFrequencyDay    RecurrenceFrequency = "Day"
	RecurrenceFrequencyHour   RecurrenceFrequency = "Hour"
	RecurrenceFrequencyMinute RecurrenceFrequency = "Minute"
	RecurrenceFrequencyMonth  RecurrenceFrequency = "Month"
	RecurrenceFrequencyNone   RecurrenceFrequency = "None"
	RecurrenceFrequencySecond RecurrenceFrequency = "Second"
	RecurrenceFrequencyWeek   RecurrenceFrequency = "Week"
	RecurrenceFrequencyYear   RecurrenceFrequency = "Year"
)

func PossibleValuesForRecurrenceFrequency() []string {
	return []string{
		string(RecurrenceFrequencyDay),
		string(RecurrenceFrequencyHour),
		string(RecurrenceFrequencyMinute),
		string(RecurrenceFrequencyMonth),
		string(RecurrenceFrequencyNone),
		string(RecurrenceFrequencySecond),
		string(RecurrenceFrequencyWeek),
		string(RecurrenceFrequencyYear),
	}
}

func (s *RecurrenceFrequency) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseRecurrenceFrequency(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseRecurrenceFrequency(input string) (*RecurrenceFrequency, error) {
	vals := map[string]RecurrenceFrequency{
		"day":    RecurrenceFrequencyDay,
		"hour":   RecurrenceFrequencyHour,
		"minute": RecurrenceFrequencyMinute,
		"month":  RecurrenceFrequencyMonth,
		"none":   RecurrenceFrequencyNone,
		"second": RecurrenceFrequencySecond,
		"week":   RecurrenceFrequencyWeek,
		"year":   RecurrenceFrequencyYear,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RecurrenceFrequency(input)
	return &out, nil
}

type ScaleDirection string

const (
	ScaleDirectionDecrease ScaleDirection = "Decrease"
	ScaleDirectionIncrease ScaleDirection = "Increase"
	ScaleDirectionNone     ScaleDirection = "None"
)

func PossibleValuesForScaleDirection() []string {
	return []string{
		string(ScaleDirectionDecrease),
		string(ScaleDirectionIncrease),
		string(ScaleDirectionNone),
	}
}

func (s *ScaleDirection) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseScaleDirection(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseScaleDirection(input string) (*ScaleDirection, error) {
	vals := map[string]ScaleDirection{
		"decrease": ScaleDirectionDecrease,
		"increase": ScaleDirectionIncrease,
		"none":     ScaleDirectionNone,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ScaleDirection(input)
	return &out, nil
}

type ScaleRuleMetricDimensionOperationType string

const (
	ScaleRuleMetricDimensionOperationTypeEquals    ScaleRuleMetricDimensionOperationType = "Equals"
	ScaleRuleMetricDimensionOperationTypeNotEquals ScaleRuleMetricDimensionOperationType = "NotEquals"
)

func PossibleValuesForScaleRuleMetricDimensionOperationType() []string {
	return []string{
		string(ScaleRuleMetricDimensionOperationTypeEquals),
		string(ScaleRuleMetricDimensionOperationTypeNotEquals),
	}
}

func (s *ScaleRuleMetricDimensionOperationType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseScaleRuleMetricDimensionOperationType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseScaleRuleMetricDimensionOperationType(input string) (*ScaleRuleMetricDimensionOperationType, error) {
	vals := map[string]ScaleRuleMetricDimensionOperationType{
		"equals":    ScaleRuleMetricDimensionOperationTypeEquals,
		"notequals": ScaleRuleMetricDimensionOperationTypeNotEquals,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ScaleRuleMetricDimensionOperationType(input)
	return &out, nil
}

type ScaleType string

const (
	ScaleTypeChangeCount             ScaleType = "ChangeCount"
	ScaleTypeExactCount              ScaleType = "ExactCount"
	ScaleTypePercentChangeCount      ScaleType = "PercentChangeCount"
	ScaleTypeServiceAllowedNextValue ScaleType = "ServiceAllowedNextValue"
)

func PossibleValuesForScaleType() []string {
	return []string{
		string(ScaleTypeChangeCount),
		string(ScaleTypeExactCount),
		string(ScaleTypePercentChangeCount),
		string(ScaleTypeServiceAllowedNextValue),
	}
}

func (s *ScaleType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseScaleType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseScaleType(input string) (*ScaleType, error) {
	vals := map[string]ScaleType{
		"changecount":             ScaleTypeChangeCount,
		"exactcount":              ScaleTypeExactCount,
		"percentchangecount":      ScaleTypePercentChangeCount,
		"serviceallowednextvalue": ScaleTypeServiceAllowedNextValue,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ScaleType(input)
	return &out, nil
}

type TimeAggregationType string

const (
	TimeAggregationTypeAverage TimeAggregationType = "Average"
	TimeAggregationTypeCount   TimeAggregationType = "Count"
	TimeAggregationTypeLast    TimeAggregationType = "Last"
	TimeAggregationTypeMaximum TimeAggregationType = "Maximum"
	TimeAggregationTypeMinimum TimeAggregationType = "Minimum"
	TimeAggregationTypeTotal   TimeAggregationType = "Total"
)

func PossibleValuesForTimeAggregationType() []string {
	return []string{
		string(TimeAggregationTypeAverage),
		string(TimeAggregationTypeCount),
		string(TimeAggregationTypeLast),
		string(TimeAggregationTypeMaximum),
		string(TimeAggregationTypeMinimum),
		string(TimeAggregationTypeTotal),
	}
}

func (s *TimeAggregationType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseTimeAggregationType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseTimeAggregationType(input string) (*TimeAggregationType, error) {
	vals := map[string]TimeAggregationType{
		"average": TimeAggregationTypeAverage,
		"count":   TimeAggregationTypeCount,
		"last":    TimeAggregationTypeLast,
		"maximum": TimeAggregationTypeMaximum,
		"minimum": TimeAggregationTypeMinimum,
		"total":   TimeAggregationTypeTotal,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := TimeAggregationType(input)
	return &out, nil
}
