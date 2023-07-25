package environments

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EnvironmentUpdateParameters interface {
}

func unmarshalEnvironmentUpdateParametersImplementation(input []byte) (EnvironmentUpdateParameters, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling EnvironmentUpdateParameters into map[string]interface: %+v", err)
	}

	value, ok := temp["kind"].(string)
	if !ok {
		return nil, nil
	}

	if strings.EqualFold(value, "Gen1") {
		var out Gen1EnvironmentUpdateParameters
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into Gen1EnvironmentUpdateParameters: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Gen2") {
		var out Gen2EnvironmentUpdateParameters
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into Gen2EnvironmentUpdateParameters: %+v", err)
		}
		return out, nil
	}

	type RawEnvironmentUpdateParametersImpl struct {
		Type   string                 `json:"-"`
		Values map[string]interface{} `json:"-"`
	}
	out := RawEnvironmentUpdateParametersImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
