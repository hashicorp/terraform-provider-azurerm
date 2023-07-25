package functions

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ FunctionRetrieveDefaultDefinitionParameters = JavaScriptFunctionRetrieveDefaultDefinitionParameters{}

type JavaScriptFunctionRetrieveDefaultDefinitionParameters struct {
	BindingRetrievalProperties *JavaScriptFunctionBindingRetrievalProperties `json:"bindingRetrievalProperties,omitempty"`

	// Fields inherited from FunctionRetrieveDefaultDefinitionParameters
}

var _ json.Marshaler = JavaScriptFunctionRetrieveDefaultDefinitionParameters{}

func (s JavaScriptFunctionRetrieveDefaultDefinitionParameters) MarshalJSON() ([]byte, error) {
	type wrapper JavaScriptFunctionRetrieveDefaultDefinitionParameters
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling JavaScriptFunctionRetrieveDefaultDefinitionParameters: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling JavaScriptFunctionRetrieveDefaultDefinitionParameters: %+v", err)
	}
	decoded["bindingType"] = "Microsoft.StreamAnalytics/JavascriptUdf"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling JavaScriptFunctionRetrieveDefaultDefinitionParameters: %+v", err)
	}

	return encoded, nil
}
