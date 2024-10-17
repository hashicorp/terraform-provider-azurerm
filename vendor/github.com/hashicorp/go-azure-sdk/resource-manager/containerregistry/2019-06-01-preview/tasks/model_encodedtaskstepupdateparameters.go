package tasks

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ TaskStepUpdateParameters = EncodedTaskStepUpdateParameters{}

type EncodedTaskStepUpdateParameters struct {
	EncodedTaskContent   *string     `json:"encodedTaskContent,omitempty"`
	EncodedValuesContent *string     `json:"encodedValuesContent,omitempty"`
	Values               *[]SetValue `json:"values,omitempty"`

	// Fields inherited from TaskStepUpdateParameters

	ContextAccessToken *string  `json:"contextAccessToken,omitempty"`
	ContextPath        *string  `json:"contextPath,omitempty"`
	Type               StepType `json:"type"`
}

func (s EncodedTaskStepUpdateParameters) TaskStepUpdateParameters() BaseTaskStepUpdateParametersImpl {
	return BaseTaskStepUpdateParametersImpl{
		ContextAccessToken: s.ContextAccessToken,
		ContextPath:        s.ContextPath,
		Type:               s.Type,
	}
}

var _ json.Marshaler = EncodedTaskStepUpdateParameters{}

func (s EncodedTaskStepUpdateParameters) MarshalJSON() ([]byte, error) {
	type wrapper EncodedTaskStepUpdateParameters
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling EncodedTaskStepUpdateParameters: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling EncodedTaskStepUpdateParameters: %+v", err)
	}

	decoded["type"] = "EncodedTask"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling EncodedTaskStepUpdateParameters: %+v", err)
	}

	return encoded, nil
}
