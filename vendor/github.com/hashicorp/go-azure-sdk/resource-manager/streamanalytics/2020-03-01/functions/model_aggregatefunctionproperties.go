package functions

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ FunctionProperties = AggregateFunctionProperties{}

type AggregateFunctionProperties struct {

	// Fields inherited from FunctionProperties

	Etag       *string                `json:"etag,omitempty"`
	Properties *FunctionConfiguration `json:"properties,omitempty"`
	Type       string                 `json:"type"`
}

func (s AggregateFunctionProperties) FunctionProperties() BaseFunctionPropertiesImpl {
	return BaseFunctionPropertiesImpl{
		Etag:       s.Etag,
		Properties: s.Properties,
		Type:       s.Type,
	}
}

var _ json.Marshaler = AggregateFunctionProperties{}

func (s AggregateFunctionProperties) MarshalJSON() ([]byte, error) {
	type wrapper AggregateFunctionProperties
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AggregateFunctionProperties: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AggregateFunctionProperties: %+v", err)
	}

	decoded["type"] = "Aggregate"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AggregateFunctionProperties: %+v", err)
	}

	return encoded, nil
}
