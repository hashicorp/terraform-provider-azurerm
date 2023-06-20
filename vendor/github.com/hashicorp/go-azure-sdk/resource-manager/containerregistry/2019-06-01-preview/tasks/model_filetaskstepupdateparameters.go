package tasks

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ TaskStepUpdateParameters = FileTaskStepUpdateParameters{}

type FileTaskStepUpdateParameters struct {
	TaskFilePath   *string     `json:"taskFilePath,omitempty"`
	Values         *[]SetValue `json:"values,omitempty"`
	ValuesFilePath *string     `json:"valuesFilePath,omitempty"`

	// Fields inherited from TaskStepUpdateParameters
	ContextAccessToken *string `json:"contextAccessToken,omitempty"`
	ContextPath        *string `json:"contextPath,omitempty"`
}

var _ json.Marshaler = FileTaskStepUpdateParameters{}

func (s FileTaskStepUpdateParameters) MarshalJSON() ([]byte, error) {
	type wrapper FileTaskStepUpdateParameters
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling FileTaskStepUpdateParameters: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling FileTaskStepUpdateParameters: %+v", err)
	}
	decoded["type"] = "FileTask"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling FileTaskStepUpdateParameters: %+v", err)
	}

	return encoded, nil
}
