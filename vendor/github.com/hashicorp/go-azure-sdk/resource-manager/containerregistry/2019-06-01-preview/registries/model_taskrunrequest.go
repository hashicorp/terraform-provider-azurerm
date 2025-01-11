package registries

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ RunRequest = TaskRunRequest{}

type TaskRunRequest struct {
	OverrideTaskStepProperties *OverrideTaskStepProperties `json:"overrideTaskStepProperties,omitempty"`
	TaskId                     string                      `json:"taskId"`

	// Fields inherited from RunRequest

	AgentPoolName    *string `json:"agentPoolName,omitempty"`
	IsArchiveEnabled *bool   `json:"isArchiveEnabled,omitempty"`
	LogTemplate      *string `json:"logTemplate,omitempty"`
	Type             string  `json:"type"`
}

func (s TaskRunRequest) RunRequest() BaseRunRequestImpl {
	return BaseRunRequestImpl{
		AgentPoolName:    s.AgentPoolName,
		IsArchiveEnabled: s.IsArchiveEnabled,
		LogTemplate:      s.LogTemplate,
		Type:             s.Type,
	}
}

var _ json.Marshaler = TaskRunRequest{}

func (s TaskRunRequest) MarshalJSON() ([]byte, error) {
	type wrapper TaskRunRequest
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling TaskRunRequest: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling TaskRunRequest: %+v", err)
	}

	decoded["type"] = "TaskRunRequest"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling TaskRunRequest: %+v", err)
	}

	return encoded, nil
}
