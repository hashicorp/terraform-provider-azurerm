package streamingjobs

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ FunctionProperties = ScalarFunctionProperties{}

type ScalarFunctionProperties struct {

	// Fields inherited from FunctionProperties

	Etag       *string                `json:"etag,omitempty"`
	Properties *FunctionConfiguration `json:"properties,omitempty"`
	Type       string                 `json:"type"`
}

func (s ScalarFunctionProperties) FunctionProperties() BaseFunctionPropertiesImpl {
	return BaseFunctionPropertiesImpl{
		Etag:       s.Etag,
		Properties: s.Properties,
		Type:       s.Type,
	}
}

var _ json.Marshaler = ScalarFunctionProperties{}

func (s ScalarFunctionProperties) MarshalJSON() ([]byte, error) {
	type wrapper ScalarFunctionProperties
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ScalarFunctionProperties: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ScalarFunctionProperties: %+v", err)
	}

	decoded["type"] = "Scalar"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ScalarFunctionProperties: %+v", err)
	}

	return encoded, nil
}
