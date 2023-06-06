package metricalerts

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ MetricAlertCriteria = MetricAlertMultipleResourceMultipleMetricCriteria{}

type MetricAlertMultipleResourceMultipleMetricCriteria struct {
	AllOf *[]MultiMetricCriteria `json:"allOf,omitempty"`

	// Fields inherited from MetricAlertCriteria
}

var _ json.Marshaler = MetricAlertMultipleResourceMultipleMetricCriteria{}

func (s MetricAlertMultipleResourceMultipleMetricCriteria) MarshalJSON() ([]byte, error) {
	type wrapper MetricAlertMultipleResourceMultipleMetricCriteria
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling MetricAlertMultipleResourceMultipleMetricCriteria: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling MetricAlertMultipleResourceMultipleMetricCriteria: %+v", err)
	}
	decoded["odata.type"] = "Microsoft.Azure.Monitor.MultipleResourceMultipleMetricCriteria"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling MetricAlertMultipleResourceMultipleMetricCriteria: %+v", err)
	}

	return encoded, nil
}

var _ json.Unmarshaler = &MetricAlertMultipleResourceMultipleMetricCriteria{}

func (s *MetricAlertMultipleResourceMultipleMetricCriteria) UnmarshalJSON(bytes []byte) error {

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling MetricAlertMultipleResourceMultipleMetricCriteria into map[string]json.RawMessage: %+v", err)
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
				return fmt.Errorf("unmarshaling index %d field 'AllOf' for 'MetricAlertMultipleResourceMultipleMetricCriteria': %+v", i, err)
			}
			output = append(output, impl)
		}
		s.AllOf = &output
	}
	return nil
}
