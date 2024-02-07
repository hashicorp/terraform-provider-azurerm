package metricalerts

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ MetricAlertCriteria = MetricAlertSingleResourceMultipleMetricCriteria{}

type MetricAlertSingleResourceMultipleMetricCriteria struct {
	AllOf *[]MetricCriteria `json:"allOf,omitempty"`

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
