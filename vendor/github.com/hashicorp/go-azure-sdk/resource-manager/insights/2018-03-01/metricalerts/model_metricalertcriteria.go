package metricalerts

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MetricAlertCriteria interface {
	MetricAlertCriteria() BaseMetricAlertCriteriaImpl
}

var _ MetricAlertCriteria = BaseMetricAlertCriteriaImpl{}

type BaseMetricAlertCriteriaImpl struct {
	OdataType Odatatype `json:"odata.type"`
}

func (s BaseMetricAlertCriteriaImpl) MetricAlertCriteria() BaseMetricAlertCriteriaImpl {
	return s
}

var _ MetricAlertCriteria = RawMetricAlertCriteriaImpl{}

// RawMetricAlertCriteriaImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawMetricAlertCriteriaImpl struct {
	metricAlertCriteria BaseMetricAlertCriteriaImpl
	Type                string
	Values              map[string]interface{}
}

func (s RawMetricAlertCriteriaImpl) MetricAlertCriteria() BaseMetricAlertCriteriaImpl {
	return s.metricAlertCriteria
}

func UnmarshalMetricAlertCriteriaImplementation(input []byte) (MetricAlertCriteria, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling MetricAlertCriteria into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["odata.type"]; ok {
		value = fmt.Sprintf("%v", v)
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

	var parent BaseMetricAlertCriteriaImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseMetricAlertCriteriaImpl: %+v", err)
	}

	return RawMetricAlertCriteriaImpl{
		metricAlertCriteria: parent,
		Type:                value,
		Values:              temp,
	}, nil

}
