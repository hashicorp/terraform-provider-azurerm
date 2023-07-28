package schedule

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ JobInput = MLFlowModelJobInput{}

type MLFlowModelJobInput struct {
	Mode *InputDeliveryMode `json:"mode,omitempty"`
	Uri  string             `json:"uri"`

	// Fields inherited from JobInput
	Description *string `json:"description,omitempty"`
}

var _ json.Marshaler = MLFlowModelJobInput{}

func (s MLFlowModelJobInput) MarshalJSON() ([]byte, error) {
	type wrapper MLFlowModelJobInput
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling MLFlowModelJobInput: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling MLFlowModelJobInput: %+v", err)
	}
	decoded["jobInputType"] = "mlflow_model"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling MLFlowModelJobInput: %+v", err)
	}

	return encoded, nil
}
