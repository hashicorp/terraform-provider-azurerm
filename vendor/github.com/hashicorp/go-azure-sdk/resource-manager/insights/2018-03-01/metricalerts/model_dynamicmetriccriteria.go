package metricalerts

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ MultiMetricCriteria = DynamicMetricCriteria{}

type DynamicMetricCriteria struct {
	AlertSensitivity DynamicThresholdSensitivity    `json:"alertSensitivity"`
	FailingPeriods   DynamicThresholdFailingPeriods `json:"failingPeriods"`
	IgnoreDataBefore *string                        `json:"ignoreDataBefore,omitempty"`
	Operator         DynamicThresholdOperator       `json:"operator"`

	// Fields inherited from MultiMetricCriteria

	CriterionType        CriterionType       `json:"criterionType"`
	Dimensions           *[]MetricDimension  `json:"dimensions,omitempty"`
	MetricName           string              `json:"metricName"`
	MetricNamespace      *string             `json:"metricNamespace,omitempty"`
	Name                 string              `json:"name"`
	SkipMetricValidation *bool               `json:"skipMetricValidation,omitempty"`
	TimeAggregation      AggregationTypeEnum `json:"timeAggregation"`
}

func (s DynamicMetricCriteria) MultiMetricCriteria() BaseMultiMetricCriteriaImpl {
	return BaseMultiMetricCriteriaImpl{
		CriterionType:        s.CriterionType,
		Dimensions:           s.Dimensions,
		MetricName:           s.MetricName,
		MetricNamespace:      s.MetricNamespace,
		Name:                 s.Name,
		SkipMetricValidation: s.SkipMetricValidation,
		TimeAggregation:      s.TimeAggregation,
	}
}

var _ json.Marshaler = DynamicMetricCriteria{}

func (s DynamicMetricCriteria) MarshalJSON() ([]byte, error) {
	type wrapper DynamicMetricCriteria
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling DynamicMetricCriteria: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling DynamicMetricCriteria: %+v", err)
	}

	decoded["criterionType"] = "DynamicThresholdCriterion"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling DynamicMetricCriteria: %+v", err)
	}

	return encoded, nil
}
