package streamingjobs

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ FunctionBinding = AzureMachineLearningWebServiceFunctionBinding{}

type AzureMachineLearningWebServiceFunctionBinding struct {
	Properties *AzureMachineLearningWebServiceFunctionBindingProperties `json:"properties,omitempty"`

	// Fields inherited from FunctionBinding

	Type string `json:"type"`
}

func (s AzureMachineLearningWebServiceFunctionBinding) FunctionBinding() BaseFunctionBindingImpl {
	return BaseFunctionBindingImpl{
		Type: s.Type,
	}
}

var _ json.Marshaler = AzureMachineLearningWebServiceFunctionBinding{}

func (s AzureMachineLearningWebServiceFunctionBinding) MarshalJSON() ([]byte, error) {
	type wrapper AzureMachineLearningWebServiceFunctionBinding
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AzureMachineLearningWebServiceFunctionBinding: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AzureMachineLearningWebServiceFunctionBinding: %+v", err)
	}

	decoded["type"] = "Microsoft.MachineLearning/WebService"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AzureMachineLearningWebServiceFunctionBinding: %+v", err)
	}

	return encoded, nil
}
