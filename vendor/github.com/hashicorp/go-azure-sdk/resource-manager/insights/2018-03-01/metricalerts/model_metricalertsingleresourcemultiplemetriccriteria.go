package metricalerts

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ MetricAlertCriteria = MetricAlertSingleResourceMultipleMetricCriteria{}

type MetricAlertSingleResourceMultipleMetricCriteria struct {
	AllOf *[]MultiMetricCriteria `json:"allOf,omitempty"`

	// Fields inherited from MetricAlertCriteria
}

var _ json.Marshaler = MetricAlertSingleResourceMultipleMetricCriteria{}

func (s MetricAlertSingleResourceMultipleMetricCriteria) MarshalJSON() ([]byte, error) {
	type wrapper MetricAlertSingleResourceMultipleMetricCriteria
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling MetricAlertSingleResourceMultipleMetricCriteria: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling MetricAlertSingleResourceMultipleMetricCriteria: %+v", err)
	}
	decoded["odata.type"] = "Microsoft.Azure.Monitor.SingleResourceMultipleMetricCriteria"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling MetricAlertSingleResourceMultipleMetricCriteria: %+v", err)
	}

	return encoded, nil
}

var _ json.Unmarshaler = &MetricAlertSingleResourceMultipleMetricCriteria{}

func (s *MetricAlertSingleResourceMultipleMetricCriteria) UnmarshalJSON(bytes []byte) error {

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling MetricAlertSingleResourceMultipleMetricCriteria into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["allOf"]; ok {
		var listTemp []json.RawMessage
		if err := json.Unmarshal(v, &listTemp); err != nil {
			return fmt.Errorf("unmarshaling AllOf into list []json.RawMessage: %+v", err)
		}

		output := make([]MultiMetricCriteria, 0)
		for i, val := range listTemp {
			impl, err := unmarshalMultiMetricCriteriaImplementation(val)
			if err != nil {
				return fmt.Errorf("unmarshaling index %d field 'AllOf' for 'MetricAlertSingleResourceMultipleMetricCriteria': %+v", i, err)
			}
			output = append(output, impl)
		}
		s.AllOf = &output
	}
	return nil
}
