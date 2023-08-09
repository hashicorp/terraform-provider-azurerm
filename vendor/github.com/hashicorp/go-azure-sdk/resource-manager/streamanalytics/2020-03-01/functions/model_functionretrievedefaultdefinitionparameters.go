package functions

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FunctionRetrieveDefaultDefinitionParameters interface {
}

// RawModeOfTransitImpl is returned when the Discriminated Value
// doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawFunctionRetrieveDefaultDefinitionParametersImpl struct {
	Type   string
	Values map[string]interface{}
}

func unmarshalFunctionRetrieveDefaultDefinitionParametersImplementation(input []byte) (FunctionRetrieveDefaultDefinitionParameters, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling FunctionRetrieveDefaultDefinitionParameters into map[string]interface: %+v", err)
	}

	value, ok := temp["bindingType"].(string)
	if !ok {
		return nil, nil
	}

	if strings.EqualFold(value, "Microsoft.MachineLearning/WebService") {
		var out AzureMachineLearningWebServiceFunctionRetrieveDefaultDefinitionParameters
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureMachineLearningWebServiceFunctionRetrieveDefaultDefinitionParameters: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Microsoft.StreamAnalytics/JavascriptUdf") {
		var out JavaScriptFunctionRetrieveDefaultDefinitionParameters
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into JavaScriptFunctionRetrieveDefaultDefinitionParameters: %+v", err)
		}
		return out, nil
	}

	out := RawFunctionRetrieveDefaultDefinitionParametersImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
