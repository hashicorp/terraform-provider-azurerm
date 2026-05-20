package metricalerts

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MultiMetricCriteria interface {
	MultiMetricCriteria() BaseMultiMetricCriteriaImpl
}

var _ MultiMetricCriteria = BaseMultiMetricCriteriaImpl{}

type BaseMultiMetricCriteriaImpl struct {
	CriterionType        CriterionType       `json:"criterionType"`
	Dimensions           *[]MetricDimension  `json:"dimensions,omitempty"`
	MetricName           string              `json:"metricName"`
	MetricNamespace      *string             `json:"metricNamespace,omitempty"`
	Name                 string              `json:"name"`
	SkipMetricValidation *bool               `json:"skipMetricValidation,omitempty"`
	TimeAggregation      AggregationTypeEnum `json:"timeAggregation"`
}

func (s BaseMultiMetricCriteriaImpl) MultiMetricCriteria() BaseMultiMetricCriteriaImpl {
	return s
}

var _ MultiMetricCriteria = RawMultiMetricCriteriaImpl{}

// RawMultiMetricCriteriaImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawMultiMetricCriteriaImpl struct {
	multiMetricCriteria BaseMultiMetricCriteriaImpl
	Type                string
	Values              map[string]interface{}
}

func (s RawMultiMetricCriteriaImpl) MultiMetricCriteria() BaseMultiMetricCriteriaImpl {
	return s.multiMetricCriteria
}

func UnmarshalMultiMetricCriteriaImplementation(input []byte) (MultiMetricCriteria, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling MultiMetricCriteria into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["criterionType"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "DynamicThresholdCriterion") {
		var out DynamicMetricCriteria
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into DynamicMetricCriteria: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "StaticThresholdCriterion") {
		var out MetricCriteria
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MetricCriteria: %+v", err)
		}
		return out, nil
	}

	var parent BaseMultiMetricCriteriaImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseMultiMetricCriteriaImpl: %+v", err)
	}

	return RawMultiMetricCriteriaImpl{
		multiMetricCriteria: parent,
		Type:                value,
		Values:              temp,
	}, nil

}
