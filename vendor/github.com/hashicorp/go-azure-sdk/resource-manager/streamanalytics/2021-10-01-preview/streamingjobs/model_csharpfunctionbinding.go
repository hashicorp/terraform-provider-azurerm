package streamingjobs

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ FunctionBinding = CSharpFunctionBinding{}

type CSharpFunctionBinding struct {
	Properties *CSharpFunctionBindingProperties `json:"properties,omitempty"`

	// Fields inherited from FunctionBinding

	Type string `json:"type"`
}

func (s CSharpFunctionBinding) FunctionBinding() BaseFunctionBindingImpl {
	return BaseFunctionBindingImpl{
		Type: s.Type,
	}
}

var _ json.Marshaler = CSharpFunctionBinding{}

func (s CSharpFunctionBinding) MarshalJSON() ([]byte, error) {
	type wrapper CSharpFunctionBinding
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling CSharpFunctionBinding: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling CSharpFunctionBinding: %+v", err)
	}

	decoded["type"] = "Microsoft.StreamAnalytics/CLRUdf"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling CSharpFunctionBinding: %+v", err)
	}

	return encoded, nil
}
