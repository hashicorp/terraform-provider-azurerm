package schedule

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ JobOutput = TritonModelJobOutput{}

type TritonModelJobOutput struct {
	Mode *OutputDeliveryMode `json:"mode,omitempty"`
	Uri  *string             `json:"uri,omitempty"`

	// Fields inherited from JobOutput
	Description *string `json:"description,omitempty"`
}

var _ json.Marshaler = TritonModelJobOutput{}

func (s TritonModelJobOutput) MarshalJSON() ([]byte, error) {
	type wrapper TritonModelJobOutput
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling TritonModelJobOutput: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling TritonModelJobOutput: %+v", err)
	}
	decoded["jobOutputType"] = "triton_model"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling TritonModelJobOutput: %+v", err)
	}

	return encoded, nil
}
