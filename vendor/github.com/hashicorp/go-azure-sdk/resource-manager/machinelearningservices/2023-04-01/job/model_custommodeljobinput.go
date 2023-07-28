package job

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ JobInput = CustomModelJobInput{}

type CustomModelJobInput struct {
	Mode *InputDeliveryMode `json:"mode,omitempty"`
	Uri  string             `json:"uri"`

	// Fields inherited from JobInput
	Description *string `json:"description,omitempty"`
}

var _ json.Marshaler = CustomModelJobInput{}

func (s CustomModelJobInput) MarshalJSON() ([]byte, error) {
	type wrapper CustomModelJobInput
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling CustomModelJobInput: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling CustomModelJobInput: %+v", err)
	}
	decoded["jobInputType"] = "custom_model"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling CustomModelJobInput: %+v", err)
	}

	return encoded, nil
}
