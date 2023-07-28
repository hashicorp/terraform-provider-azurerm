package schedule

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ForecastHorizon = AutoForecastHorizon{}

type AutoForecastHorizon struct {

	// Fields inherited from ForecastHorizon
}

var _ json.Marshaler = AutoForecastHorizon{}

func (s AutoForecastHorizon) MarshalJSON() ([]byte, error) {
	type wrapper AutoForecastHorizon
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AutoForecastHorizon: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AutoForecastHorizon: %+v", err)
	}
	decoded["mode"] = "Auto"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AutoForecastHorizon: %+v", err)
	}

	return encoded, nil
}
