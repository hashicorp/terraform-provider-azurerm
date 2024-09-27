package environments

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EnvironmentUpdateParameters interface {
	EnvironmentUpdateParameters() BaseEnvironmentUpdateParametersImpl
}

var _ EnvironmentUpdateParameters = BaseEnvironmentUpdateParametersImpl{}

type BaseEnvironmentUpdateParametersImpl struct {
	Kind EnvironmentKind    `json:"kind"`
	Tags *map[string]string `json:"tags,omitempty"`
}

func (s BaseEnvironmentUpdateParametersImpl) EnvironmentUpdateParameters() BaseEnvironmentUpdateParametersImpl {
	return s
}

var _ EnvironmentUpdateParameters = RawEnvironmentUpdateParametersImpl{}

// RawEnvironmentUpdateParametersImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawEnvironmentUpdateParametersImpl struct {
	environmentUpdateParameters BaseEnvironmentUpdateParametersImpl
	Type                        string
	Values                      map[string]interface{}
}

func (s RawEnvironmentUpdateParametersImpl) EnvironmentUpdateParameters() BaseEnvironmentUpdateParametersImpl {
	return s.environmentUpdateParameters
}

func UnmarshalEnvironmentUpdateParametersImplementation(input []byte) (EnvironmentUpdateParameters, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling EnvironmentUpdateParameters into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["kind"]; ok {
		value = fmt.Sprintf("%v", v)
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

	var parent BaseEnvironmentUpdateParametersImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseEnvironmentUpdateParametersImpl: %+v", err)
	}

	return RawEnvironmentUpdateParametersImpl{
		environmentUpdateParameters: parent,
		Type:                        value,
		Values:                      temp,
	}, nil

}
