package functions

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FunctionRetrieveDefaultDefinitionParameters interface {
	FunctionRetrieveDefaultDefinitionParameters() BaseFunctionRetrieveDefaultDefinitionParametersImpl
}

var _ FunctionRetrieveDefaultDefinitionParameters = BaseFunctionRetrieveDefaultDefinitionParametersImpl{}

type BaseFunctionRetrieveDefaultDefinitionParametersImpl struct {
	BindingType string `json:"bindingType"`
}

func (s BaseFunctionRetrieveDefaultDefinitionParametersImpl) FunctionRetrieveDefaultDefinitionParameters() BaseFunctionRetrieveDefaultDefinitionParametersImpl {
	return s
}

var _ FunctionRetrieveDefaultDefinitionParameters = RawFunctionRetrieveDefaultDefinitionParametersImpl{}

// RawFunctionRetrieveDefaultDefinitionParametersImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawFunctionRetrieveDefaultDefinitionParametersImpl struct {
	functionRetrieveDefaultDefinitionParameters BaseFunctionRetrieveDefaultDefinitionParametersImpl
	Type                                        string
	Values                                      map[string]interface{}
}

func (s RawFunctionRetrieveDefaultDefinitionParametersImpl) FunctionRetrieveDefaultDefinitionParameters() BaseFunctionRetrieveDefaultDefinitionParametersImpl {
	return s.functionRetrieveDefaultDefinitionParameters
}

func UnmarshalFunctionRetrieveDefaultDefinitionParametersImplementation(input []byte) (FunctionRetrieveDefaultDefinitionParameters, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling FunctionRetrieveDefaultDefinitionParameters into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["bindingType"]; ok {
		value = fmt.Sprintf("%v", v)
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

	var parent BaseFunctionRetrieveDefaultDefinitionParametersImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseFunctionRetrieveDefaultDefinitionParametersImpl: %+v", err)
	}

	return RawFunctionRetrieveDefaultDefinitionParametersImpl{
		functionRetrieveDefaultDefinitionParameters: parent,
		Type:   value,
		Values: temp,
	}, nil

}
