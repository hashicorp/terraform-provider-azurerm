package streamingjobs

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FunctionBinding interface {
	FunctionBinding() BaseFunctionBindingImpl
}

var _ FunctionBinding = BaseFunctionBindingImpl{}

type BaseFunctionBindingImpl struct {
	Type string `json:"type"`
}

func (s BaseFunctionBindingImpl) FunctionBinding() BaseFunctionBindingImpl {
	return s
}

var _ FunctionBinding = RawFunctionBindingImpl{}

// RawFunctionBindingImpl is returned when the Discriminated Value doesn't match any of the defined types.
// It can also be used as a Request Payload to provide a raw JSON payload, which is useful
// for preserving arbitrary/extensible JSON properties across a round-trip.
type RawFunctionBindingImpl struct {
	functionBinding BaseFunctionBindingImpl
	Type            string
	Values          map[string]interface{}
}

func (s RawFunctionBindingImpl) FunctionBinding() BaseFunctionBindingImpl {
	return s.functionBinding
}

func (s RawFunctionBindingImpl) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Values)
}

func UnmarshalFunctionBindingImplementation(input []byte) (FunctionBinding, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling FunctionBinding into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["type"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "Microsoft.MachineLearning/WebService") {
		var out AzureMachineLearningWebServiceFunctionBinding
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureMachineLearningWebServiceFunctionBinding: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Microsoft.StreamAnalytics/JavascriptUdf") {
		var out JavaScriptFunctionBinding
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into JavaScriptFunctionBinding: %+v", err)
		}
		return out, nil
	}

	var parent BaseFunctionBindingImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseFunctionBindingImpl: %+v", err)
	}

	return RawFunctionBindingImpl{
		functionBinding: parent,
		Type:            value,
		Values:          temp,
	}, nil

}
