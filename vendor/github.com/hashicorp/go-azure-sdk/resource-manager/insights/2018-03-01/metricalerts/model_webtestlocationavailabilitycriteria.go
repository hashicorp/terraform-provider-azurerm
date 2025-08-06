package metricalerts

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ MetricAlertCriteria = WebtestLocationAvailabilityCriteria{}

type WebtestLocationAvailabilityCriteria struct {
	ComponentId         string  `json:"componentId"`
	FailedLocationCount float64 `json:"failedLocationCount"`
	WebTestId           string  `json:"webTestId"`

	// Fields inherited from MetricAlertCriteria

	OdataType Odatatype `json:"odata.type"`
}

func (s WebtestLocationAvailabilityCriteria) MetricAlertCriteria() BaseMetricAlertCriteriaImpl {
	return BaseMetricAlertCriteriaImpl{
		OdataType: s.OdataType,
	}
}

var _ json.Marshaler = WebtestLocationAvailabilityCriteria{}

func (s WebtestLocationAvailabilityCriteria) MarshalJSON() ([]byte, error) {
	type wrapper WebtestLocationAvailabilityCriteria
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling WebtestLocationAvailabilityCriteria: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling WebtestLocationAvailabilityCriteria: %+v", err)
	}

	decoded["odata.type"] = "Microsoft.Azure.Monitor.WebtestLocationAvailabilityCriteria"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling WebtestLocationAvailabilityCriteria: %+v", err)
	}

	return encoded, nil
}
