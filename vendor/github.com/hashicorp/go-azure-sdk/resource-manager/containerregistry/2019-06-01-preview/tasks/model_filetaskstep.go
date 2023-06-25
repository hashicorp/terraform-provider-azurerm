package tasks

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ TaskStepProperties = FileTaskStep{}

type FileTaskStep struct {
	TaskFilePath   string      `json:"taskFilePath"`
	Values         *[]SetValue `json:"values,omitempty"`
	ValuesFilePath *string     `json:"valuesFilePath,omitempty"`

	// Fields inherited from TaskStepProperties
	BaseImageDependencies *[]BaseImageDependency `json:"baseImageDependencies,omitempty"`
	ContextAccessToken    *string                `json:"contextAccessToken,omitempty"`
	ContextPath           *string                `json:"contextPath,omitempty"`
}

var _ json.Marshaler = FileTaskStep{}

func (s FileTaskStep) MarshalJSON() ([]byte, error) {
	type wrapper FileTaskStep
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling FileTaskStep: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling FileTaskStep: %+v", err)
	}
	decoded["type"] = "FileTask"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling FileTaskStep: %+v", err)
	}

	return encoded, nil
}
