package metricalerts

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MultiMetricCriteria interface {
}

func unmarshalMultiMetricCriteriaImplementation(input []byte) (MultiMetricCriteria, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling MultiMetricCriteria into map[string]interface: %+v", err)
	}

	value, ok := temp["criterionType"].(string)
	if !ok {
		return nil, nil
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

	type RawMultiMetricCriteriaImpl struct {
		Type   string                 `json:"-"`
		Values map[string]interface{} `json:"-"`
	}
	out := RawMultiMetricCriteriaImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
