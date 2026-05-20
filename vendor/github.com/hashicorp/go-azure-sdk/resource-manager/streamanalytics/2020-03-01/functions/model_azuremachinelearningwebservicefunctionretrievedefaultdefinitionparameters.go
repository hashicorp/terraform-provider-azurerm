package functions

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ FunctionRetrieveDefaultDefinitionParameters = AzureMachineLearningWebServiceFunctionRetrieveDefaultDefinitionParameters{}

type AzureMachineLearningWebServiceFunctionRetrieveDefaultDefinitionParameters struct {
	BindingRetrievalProperties *AzureMachineLearningWebServiceFunctionBindingRetrievalProperties `json:"bindingRetrievalProperties,omitempty"`

	// Fields inherited from FunctionRetrieveDefaultDefinitionParameters

	BindingType string `json:"bindingType"`
}

func (s AzureMachineLearningWebServiceFunctionRetrieveDefaultDefinitionParameters) FunctionRetrieveDefaultDefinitionParameters() BaseFunctionRetrieveDefaultDefinitionParametersImpl {
	return BaseFunctionRetrieveDefaultDefinitionParametersImpl{
		BindingType: s.BindingType,
	}
}

var _ json.Marshaler = AzureMachineLearningWebServiceFunctionRetrieveDefaultDefinitionParameters{}

func (s AzureMachineLearningWebServiceFunctionRetrieveDefaultDefinitionParameters) MarshalJSON() ([]byte, error) {
	type wrapper AzureMachineLearningWebServiceFunctionRetrieveDefaultDefinitionParameters
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AzureMachineLearningWebServiceFunctionRetrieveDefaultDefinitionParameters: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AzureMachineLearningWebServiceFunctionRetrieveDefaultDefinitionParameters: %+v", err)
	}

	decoded["bindingType"] = "Microsoft.MachineLearning/WebService"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AzureMachineLearningWebServiceFunctionRetrieveDefaultDefinitionParameters: %+v", err)
	}

	return encoded, nil
}
