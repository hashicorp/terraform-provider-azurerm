package streamingjobs

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FunctionProperties interface {
	FunctionProperties() BaseFunctionPropertiesImpl
}

var _ FunctionProperties = BaseFunctionPropertiesImpl{}

type BaseFunctionPropertiesImpl struct {
	Etag       *string                `json:"etag,omitempty"`
	Properties *FunctionConfiguration `json:"properties,omitempty"`
	Type       string                 `json:"type"`
}

func (s BaseFunctionPropertiesImpl) FunctionProperties() BaseFunctionPropertiesImpl {
	return s
}

var _ FunctionProperties = RawFunctionPropertiesImpl{}

// RawFunctionPropertiesImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawFunctionPropertiesImpl struct {
	functionProperties BaseFunctionPropertiesImpl
	Type               string
	Values             map[string]interface{}
}

func (s RawFunctionPropertiesImpl) FunctionProperties() BaseFunctionPropertiesImpl {
	return s.functionProperties
}

func UnmarshalFunctionPropertiesImplementation(input []byte) (FunctionProperties, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling FunctionProperties into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["type"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "Aggregate") {
		var out AggregateFunctionProperties
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AggregateFunctionProperties: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Scalar") {
		var out ScalarFunctionProperties
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ScalarFunctionProperties: %+v", err)
		}
		return out, nil
	}

	var parent BaseFunctionPropertiesImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseFunctionPropertiesImpl: %+v", err)
	}

	return RawFunctionPropertiesImpl{
		functionProperties: parent,
		Type:               value,
		Values:             temp,
	}, nil

}
