package tasks

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ TaskStepProperties = EncodedTaskStep{}

type EncodedTaskStep struct {
	EncodedTaskContent   string      `json:"encodedTaskContent"`
	EncodedValuesContent *string     `json:"encodedValuesContent,omitempty"`
	Values               *[]SetValue `json:"values,omitempty"`

	// Fields inherited from TaskStepProperties
	BaseImageDependencies *[]BaseImageDependency `json:"baseImageDependencies,omitempty"`
	ContextAccessToken    *string                `json:"contextAccessToken,omitempty"`
	ContextPath           *string                `json:"contextPath,omitempty"`
}

var _ json.Marshaler = EncodedTaskStep{}

func (s EncodedTaskStep) MarshalJSON() ([]byte, error) {
	type wrapper EncodedTaskStep
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling EncodedTaskStep: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling EncodedTaskStep: %+v", err)
	}
	decoded["type"] = "EncodedTask"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling EncodedTaskStep: %+v", err)
	}

	return encoded, nil
}
