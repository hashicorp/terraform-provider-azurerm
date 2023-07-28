package job

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ JobOutput = CustomModelJobOutput{}

type CustomModelJobOutput struct {
	Mode *OutputDeliveryMode `json:"mode,omitempty"`
	Uri  *string             `json:"uri,omitempty"`

	// Fields inherited from JobOutput
	Description *string `json:"description,omitempty"`
}

var _ json.Marshaler = CustomModelJobOutput{}

func (s CustomModelJobOutput) MarshalJSON() ([]byte, error) {
	type wrapper CustomModelJobOutput
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling CustomModelJobOutput: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling CustomModelJobOutput: %+v", err)
	}
	decoded["jobOutputType"] = "custom_model"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling CustomModelJobOutput: %+v", err)
	}

	return encoded, nil
}
