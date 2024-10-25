package streamingjobs

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ FunctionBinding = AzureMachineLearningServiceFunctionBinding{}

type AzureMachineLearningServiceFunctionBinding struct {
	Properties *AzureMachineLearningServiceFunctionBindingProperties `json:"properties,omitempty"`

	// Fields inherited from FunctionBinding

	Type string `json:"type"`
}

func (s AzureMachineLearningServiceFunctionBinding) FunctionBinding() BaseFunctionBindingImpl {
	return BaseFunctionBindingImpl{
		Type: s.Type,
	}
}

var _ json.Marshaler = AzureMachineLearningServiceFunctionBinding{}

func (s AzureMachineLearningServiceFunctionBinding) MarshalJSON() ([]byte, error) {
	type wrapper AzureMachineLearningServiceFunctionBinding
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AzureMachineLearningServiceFunctionBinding: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AzureMachineLearningServiceFunctionBinding: %+v", err)
	}

	decoded["type"] = "Microsoft.MachineLearningServices"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AzureMachineLearningServiceFunctionBinding: %+v", err)
	}

	return encoded, nil
}
