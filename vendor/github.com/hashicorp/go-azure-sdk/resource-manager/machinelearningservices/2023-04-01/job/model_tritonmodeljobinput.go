package job

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ JobInput = TritonModelJobInput{}

type TritonModelJobInput struct {
	Mode *InputDeliveryMode `json:"mode,omitempty"`
	Uri  string             `json:"uri"`

	// Fields inherited from JobInput
	Description *string `json:"description,omitempty"`
}

var _ json.Marshaler = TritonModelJobInput{}

func (s TritonModelJobInput) MarshalJSON() ([]byte, error) {
	type wrapper TritonModelJobInput
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling TritonModelJobInput: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling TritonModelJobInput: %+v", err)
	}
	decoded["jobInputType"] = "triton_model"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling TritonModelJobInput: %+v", err)
	}

	return encoded, nil
}
