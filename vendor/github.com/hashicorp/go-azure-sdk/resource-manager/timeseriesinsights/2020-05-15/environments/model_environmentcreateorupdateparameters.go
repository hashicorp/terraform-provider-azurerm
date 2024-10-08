package environments

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EnvironmentCreateOrUpdateParameters interface {
	EnvironmentCreateOrUpdateParameters() BaseEnvironmentCreateOrUpdateParametersImpl
}

var _ EnvironmentCreateOrUpdateParameters = BaseEnvironmentCreateOrUpdateParametersImpl{}

type BaseEnvironmentCreateOrUpdateParametersImpl struct {
	Kind     EnvironmentKind    `json:"kind"`
	Location string             `json:"location"`
	Sku      Sku                `json:"sku"`
	Tags     *map[string]string `json:"tags,omitempty"`
}

func (s BaseEnvironmentCreateOrUpdateParametersImpl) EnvironmentCreateOrUpdateParameters() BaseEnvironmentCreateOrUpdateParametersImpl {
	return s
}

var _ EnvironmentCreateOrUpdateParameters = RawEnvironmentCreateOrUpdateParametersImpl{}

// RawEnvironmentCreateOrUpdateParametersImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawEnvironmentCreateOrUpdateParametersImpl struct {
	environmentCreateOrUpdateParameters BaseEnvironmentCreateOrUpdateParametersImpl
	Type                                string
	Values                              map[string]interface{}
}

func (s RawEnvironmentCreateOrUpdateParametersImpl) EnvironmentCreateOrUpdateParameters() BaseEnvironmentCreateOrUpdateParametersImpl {
	return s.environmentCreateOrUpdateParameters
}

func UnmarshalEnvironmentCreateOrUpdateParametersImplementation(input []byte) (EnvironmentCreateOrUpdateParameters, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling EnvironmentCreateOrUpdateParameters into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["kind"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "Gen1") {
		var out Gen1EnvironmentCreateOrUpdateParameters
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into Gen1EnvironmentCreateOrUpdateParameters: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Gen2") {
		var out Gen2EnvironmentCreateOrUpdateParameters
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into Gen2EnvironmentCreateOrUpdateParameters: %+v", err)
		}
		return out, nil
	}

	var parent BaseEnvironmentCreateOrUpdateParametersImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseEnvironmentCreateOrUpdateParametersImpl: %+v", err)
	}

	return RawEnvironmentCreateOrUpdateParametersImpl{
		environmentCreateOrUpdateParameters: parent,
		Type:                                value,
		Values:                              temp,
	}, nil

}
