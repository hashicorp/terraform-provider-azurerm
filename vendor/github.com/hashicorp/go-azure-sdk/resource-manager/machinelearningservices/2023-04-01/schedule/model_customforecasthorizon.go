package schedule

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ForecastHorizon = CustomForecastHorizon{}

type CustomForecastHorizon struct {
	Value int64 `json:"value"`

	// Fields inherited from ForecastHorizon
}

var _ json.Marshaler = CustomForecastHorizon{}

func (s CustomForecastHorizon) MarshalJSON() ([]byte, error) {
	type wrapper CustomForecastHorizon
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling CustomForecastHorizon: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling CustomForecastHorizon: %+v", err)
	}
	decoded["mode"] = "Custom"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling CustomForecastHorizon: %+v", err)
	}

	return encoded, nil
}
