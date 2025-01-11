package streamingjobs

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ FunctionBinding = AzureMachineLearningStudioFunctionBinding{}

type AzureMachineLearningStudioFunctionBinding struct {
	Properties *AzureMachineLearningStudioFunctionBindingProperties `json:"properties,omitempty"`

	// Fields inherited from FunctionBinding

	Type string `json:"type"`
}

func (s AzureMachineLearningStudioFunctionBinding) FunctionBinding() BaseFunctionBindingImpl {
	return BaseFunctionBindingImpl{
		Type: s.Type,
	}
}

var _ json.Marshaler = AzureMachineLearningStudioFunctionBinding{}

func (s AzureMachineLearningStudioFunctionBinding) MarshalJSON() ([]byte, error) {
	type wrapper AzureMachineLearningStudioFunctionBinding
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AzureMachineLearningStudioFunctionBinding: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AzureMachineLearningStudioFunctionBinding: %+v", err)
	}

	decoded["type"] = "Microsoft.MachineLearning/WebService"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AzureMachineLearningStudioFunctionBinding: %+v", err)
	}

	return encoded, nil
}
