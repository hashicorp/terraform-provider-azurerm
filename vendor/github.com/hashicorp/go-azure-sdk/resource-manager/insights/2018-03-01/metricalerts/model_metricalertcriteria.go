package metricalerts

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MetricAlertCriteria interface {
}

func unmarshalMetricAlertCriteriaImplementation(input []byte) (MetricAlertCriteria, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling MetricAlertCriteria into map[string]interface: %+v", err)
	}

	value, ok := temp["odata.type"].(string)
	if !ok {
		return nil, nil
	}

	if strings.EqualFold(value, "Microsoft.Azure.Monitor.MultipleResourceMultipleMetricCriteria") {
		var out MetricAlertMultipleResourceMultipleMetricCriteria
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MetricAlertMultipleResourceMultipleMetricCriteria: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Microsoft.Azure.Monitor.SingleResourceMultipleMetricCriteria") {
		var out MetricAlertSingleResourceMultipleMetricCriteria
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MetricAlertSingleResourceMultipleMetricCriteria: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Microsoft.Azure.Monitor.WebtestLocationAvailabilityCriteria") {
		var out WebtestLocationAvailabilityCriteria
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into WebtestLocationAvailabilityCriteria: %+v", err)
		}
		return out, nil
	}

	type RawMetricAlertCriteriaImpl struct {
		Type   string                 `json:"-"`
		Values map[string]interface{} `json:"-"`
	}
	out := RawMetricAlertCriteriaImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
