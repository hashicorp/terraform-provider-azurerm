package streamingjobs

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ FunctionBinding = JavaScriptFunctionBinding{}

type JavaScriptFunctionBinding struct {
	Properties *JavaScriptFunctionBindingProperties `json:"properties,omitempty"`

	// Fields inherited from FunctionBinding

	Type string `json:"type"`
}

func (s JavaScriptFunctionBinding) FunctionBinding() BaseFunctionBindingImpl {
	return BaseFunctionBindingImpl{
		Type: s.Type,
	}
}

var _ json.Marshaler = JavaScriptFunctionBinding{}

func (s JavaScriptFunctionBinding) MarshalJSON() ([]byte, error) {
	type wrapper JavaScriptFunctionBinding
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling JavaScriptFunctionBinding: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling JavaScriptFunctionBinding: %+v", err)
	}

	decoded["type"] = "Microsoft.StreamAnalytics/JavascriptUdf"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling JavaScriptFunctionBinding: %+v", err)
	}

	return encoded, nil
}
