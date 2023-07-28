package job

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Seasonality interface {
}

func unmarshalSeasonalityImplementation(input []byte) (Seasonality, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling Seasonality into map[string]interface: %+v", err)
	}

	value, ok := temp["mode"].(string)
	if !ok {
		return nil, nil
	}

	if strings.EqualFold(value, "Auto") {
		var out AutoSeasonality
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AutoSeasonality: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Custom") {
		var out CustomSeasonality
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into CustomSeasonality: %+v", err)
		}
		return out, nil
	}

	type RawSeasonalityImpl struct {
		Type   string                 `json:"-"`
		Values map[string]interface{} `json:"-"`
	}
	out := RawSeasonalityImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
