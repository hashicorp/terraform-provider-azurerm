package schedule

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ForecastHorizon interface {
}

func unmarshalForecastHorizonImplementation(input []byte) (ForecastHorizon, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling ForecastHorizon into map[string]interface: %+v", err)
	}

	value, ok := temp["mode"].(string)
	if !ok {
		return nil, nil
	}

	if strings.EqualFold(value, "Auto") {
		var out AutoForecastHorizon
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AutoForecastHorizon: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Custom") {
		var out CustomForecastHorizon
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into CustomForecastHorizon: %+v", err)
		}
		return out, nil
	}

	type RawForecastHorizonImpl struct {
		Type   string                 `json:"-"`
		Values map[string]interface{} `json:"-"`
	}
	out := RawForecastHorizonImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
