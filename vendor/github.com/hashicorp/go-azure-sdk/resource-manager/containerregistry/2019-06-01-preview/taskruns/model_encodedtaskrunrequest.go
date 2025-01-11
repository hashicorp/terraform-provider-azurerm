package taskruns

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ RunRequest = EncodedTaskRunRequest{}

type EncodedTaskRunRequest struct {
	AgentConfiguration   *AgentProperties   `json:"agentConfiguration,omitempty"`
	Credentials          *Credentials       `json:"credentials,omitempty"`
	EncodedTaskContent   string             `json:"encodedTaskContent"`
	EncodedValuesContent *string            `json:"encodedValuesContent,omitempty"`
	Platform             PlatformProperties `json:"platform"`
	SourceLocation       *string            `json:"sourceLocation,omitempty"`
	Timeout              *int64             `json:"timeout,omitempty"`
	Values               *[]SetValue        `json:"values,omitempty"`

	// Fields inherited from RunRequest

	AgentPoolName    *string `json:"agentPoolName,omitempty"`
	IsArchiveEnabled *bool   `json:"isArchiveEnabled,omitempty"`
	LogTemplate      *string `json:"logTemplate,omitempty"`
	Type             string  `json:"type"`
}

func (s EncodedTaskRunRequest) RunRequest() BaseRunRequestImpl {
	return BaseRunRequestImpl{
		AgentPoolName:    s.AgentPoolName,
		IsArchiveEnabled: s.IsArchiveEnabled,
		LogTemplate:      s.LogTemplate,
		Type:             s.Type,
	}
}

var _ json.Marshaler = EncodedTaskRunRequest{}

func (s EncodedTaskRunRequest) MarshalJSON() ([]byte, error) {
	type wrapper EncodedTaskRunRequest
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling EncodedTaskRunRequest: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling EncodedTaskRunRequest: %+v", err)
	}

	decoded["type"] = "EncodedTaskRunRequest"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling EncodedTaskRunRequest: %+v", err)
	}

	return encoded, nil
}
