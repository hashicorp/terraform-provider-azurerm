package functions

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FunctionProperties interface {
}

func unmarshalFunctionPropertiesImplementation(input []byte) (FunctionProperties, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling FunctionProperties into map[string]interface: %+v", err)
	}

	value, ok := temp["type"].(string)
	if !ok {
		return nil, nil
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

	type RawFunctionPropertiesImpl struct {
		Type   string                 `json:"-"`
		Values map[string]interface{} `json:"-"`
	}
	out := RawFunctionPropertiesImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
