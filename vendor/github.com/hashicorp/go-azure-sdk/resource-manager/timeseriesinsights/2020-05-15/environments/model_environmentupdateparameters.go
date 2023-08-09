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

// RawModeOfTransitImpl is returned when the Discriminated Value
// doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawEnvironmentUpdateParametersImpl struct {
	Type   string
	Values map[string]interface{}
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

	out := RawEnvironmentUpdateParametersImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
