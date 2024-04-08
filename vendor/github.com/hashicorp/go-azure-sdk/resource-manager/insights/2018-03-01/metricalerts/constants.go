package metricalerts

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AggregationTypeEnum string

const (
	AggregationTypeEnumAverage AggregationTypeEnum = "Average"
	AggregationTypeEnumCount   AggregationTypeEnum = "Count"
	AggregationTypeEnumMaximum AggregationTypeEnum = "Maximum"
	AggregationTypeEnumMinimum AggregationTypeEnum = "Minimum"
	AggregationTypeEnumTotal   AggregationTypeEnum = "Total"
)

func PossibleValuesForAggregationTypeEnum() []string {
	return []string{
		string(AggregationTypeEnumAverage),
		string(AggregationTypeEnumCount),
		string(AggregationTypeEnumMaximum),
		string(AggregationTypeEnumMinimum),
		string(AggregationTypeEnumTotal),
	}
}

func (s *AggregationTypeEnum) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAggregationTypeEnum(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAggregationTypeEnum(input string) (*AggregationTypeEnum, error) {
	vals := map[string]AggregationTypeEnum{
		"average": AggregationTypeEnumAverage,
		"count":   AggregationTypeEnumCount,
		"maximum": AggregationTypeEnumMaximum,
		"minimum": AggregationTypeEnumMinimum,
		"total":   AggregationTypeEnumTotal,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AggregationTypeEnum(input)
	return &out, nil
}

type CriterionType string

const (
	CriterionTypeDynamicThresholdCriterion CriterionType = "DynamicThresholdCriterion"
	CriterionTypeStaticThresholdCriterion  CriterionType = "StaticThresholdCriterion"
)

func PossibleValuesForCriterionType() []string {
	return []string{
		string(CriterionTypeDynamicThresholdCriterion),
		string(CriterionTypeStaticThresholdCriterion),
	}
}

func (s *CriterionType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCriterionType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseCriterionType(input string) (*CriterionType, error) {
	vals := map[string]CriterionType{
		"dynamicthresholdcriterion": CriterionTypeDynamicThresholdCriterion,
		"staticthresholdcriterion":  CriterionTypeStaticThresholdCriterion,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CriterionType(input)
	return &out, nil
}

type DynamicThresholdOperator string

const (
	DynamicThresholdOperatorGreaterOrLessThan DynamicThresholdOperator = "GreaterOrLessThan"
	DynamicThresholdOperatorGreaterThan       DynamicThresholdOperator = "GreaterThan"
	DynamicThresholdOperatorLessThan          DynamicThresholdOperator = "LessThan"
)

func PossibleValuesForDynamicThresholdOperator() []string {
	return []string{
		string(DynamicThresholdOperatorGreaterOrLessThan),
		string(DynamicThresholdOperatorGreaterThan),
		string(DynamicThresholdOperatorLessThan),
	}
}

func (s *DynamicThresholdOperator) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDynamicThresholdOperator(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDynamicThresholdOperator(input string) (*DynamicThresholdOperator, error) {
	vals := map[string]DynamicThresholdOperator{
		"greaterorlessthan": DynamicThresholdOperatorGreaterOrLessThan,
		"greaterthan":       DynamicThresholdOperatorGreaterThan,
		"lessthan":          DynamicThresholdOperatorLessThan,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DynamicThresholdOperator(input)
	return &out, nil
}

type DynamicThresholdSensitivity string

const (
	DynamicThresholdSensitivityHigh   DynamicThresholdSensitivity = "High"
	DynamicThresholdSensitivityLow    DynamicThresholdSensitivity = "Low"
	DynamicThresholdSensitivityMedium DynamicThresholdSensitivity = "Medium"
)

func PossibleValuesForDynamicThresholdSensitivity() []string {
	return []string{
		string(DynamicThresholdSensitivityHigh),
		string(DynamicThresholdSensitivityLow),
		string(DynamicThresholdSensitivityMedium),
	}
}

func (s *DynamicThresholdSensitivity) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDynamicThresholdSensitivity(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDynamicThresholdSensitivity(input string) (*DynamicThresholdSensitivity, error) {
	vals := map[string]DynamicThresholdSensitivity{
		"high":   DynamicThresholdSensitivityHigh,
		"low":    DynamicThresholdSensitivityLow,
		"medium": DynamicThresholdSensitivityMedium,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DynamicThresholdSensitivity(input)
	return &out, nil
}

type Odatatype string

const (
	OdatatypeMicrosoftPointAzurePointMonitorPointMultipleResourceMultipleMetricCriteria Odatatype = "Microsoft.Azure.Monitor.MultipleResourceMultipleMetricCriteria"
	OdatatypeMicrosoftPointAzurePointMonitorPointSingleResourceMultipleMetricCriteria   Odatatype = "Microsoft.Azure.Monitor.SingleResourceMultipleMetricCriteria"
	OdatatypeMicrosoftPointAzurePointMonitorPointWebtestLocationAvailabilityCriteria    Odatatype = "Microsoft.Azure.Monitor.WebtestLocationAvailabilityCriteria"
)

func PossibleValuesForOdatatype() []string {
	return []string{
		string(OdatatypeMicrosoftPointAzurePointMonitorPointMultipleResourceMultipleMetricCriteria),
		string(OdatatypeMicrosoftPointAzurePointMonitorPointSingleResourceMultipleMetricCriteria),
		string(OdatatypeMicrosoftPointAzurePointMonitorPointWebtestLocationAvailabilityCriteria),
	}
}

func (s *Odatatype) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseOdatatype(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseOdatatype(input string) (*Odatatype, error) {
	vals := map[string]Odatatype{
		"microsoft.azure.monitor.multipleresourcemultiplemetriccriteria": OdatatypeMicrosoftPointAzurePointMonitorPointMultipleResourceMultipleMetricCriteria,
		"microsoft.azure.monitor.singleresourcemultiplemetriccriteria":   OdatatypeMicrosoftPointAzurePointMonitorPointSingleResourceMultipleMetricCriteria,
		"microsoft.azure.monitor.webtestlocationavailabilitycriteria":    OdatatypeMicrosoftPointAzurePointMonitorPointWebtestLocationAvailabilityCriteria,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Odatatype(input)
	return &out, nil
}

type Operator string

const (
	OperatorEquals             Operator = "Equals"
	OperatorGreaterThan        Operator = "GreaterThan"
	OperatorGreaterThanOrEqual Operator = "GreaterThanOrEqual"
	OperatorLessThan           Operator = "LessThan"
	OperatorLessThanOrEqual    Operator = "LessThanOrEqual"
)

func PossibleValuesForOperator() []string {
	return []string{
		string(OperatorEquals),
		string(OperatorGreaterThan),
		string(OperatorGreaterThanOrEqual),
		string(OperatorLessThan),
		string(OperatorLessThanOrEqual),
	}
}

func (s *Operator) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseOperator(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseOperator(input string) (*Operator, error) {
	vals := map[string]Operator{
		"equals":             OperatorEquals,
		"greaterthan":        OperatorGreaterThan,
		"greaterthanorequal": OperatorGreaterThanOrEqual,
		"lessthan":           OperatorLessThan,
		"lessthanorequal":    OperatorLessThanOrEqual,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Operator(input)
	return &out, nil
}
