package taskruns

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ RunRequest = DockerBuildRequest{}

type DockerBuildRequest struct {
	AgentConfiguration *AgentProperties   `json:"agentConfiguration,omitempty"`
	Arguments          *[]Argument        `json:"arguments,omitempty"`
	Credentials        *Credentials       `json:"credentials,omitempty"`
	DockerFilePath     string             `json:"dockerFilePath"`
	ImageNames         *[]string          `json:"imageNames,omitempty"`
	IsPushEnabled      *bool              `json:"isPushEnabled,omitempty"`
	NoCache            *bool              `json:"noCache,omitempty"`
	Platform           PlatformProperties `json:"platform"`
	SourceLocation     *string            `json:"sourceLocation,omitempty"`
	Target             *string            `json:"target,omitempty"`
	Timeout            *int64             `json:"timeout,omitempty"`

	// Fields inherited from RunRequest

	AgentPoolName    *string `json:"agentPoolName,omitempty"`
	IsArchiveEnabled *bool   `json:"isArchiveEnabled,omitempty"`
	LogTemplate      *string `json:"logTemplate,omitempty"`
	Type             string  `json:"type"`
}

func (s DockerBuildRequest) RunRequest() BaseRunRequestImpl {
	return BaseRunRequestImpl{
		AgentPoolName:    s.AgentPoolName,
		IsArchiveEnabled: s.IsArchiveEnabled,
		LogTemplate:      s.LogTemplate,
		Type:             s.Type,
	}
}

var _ json.Marshaler = DockerBuildRequest{}

func (s DockerBuildRequest) MarshalJSON() ([]byte, error) {
	type wrapper DockerBuildRequest
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling DockerBuildRequest: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling DockerBuildRequest: %+v", err)
	}

	decoded["type"] = "DockerBuildRequest"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling DockerBuildRequest: %+v", err)
	}

	return encoded, nil
}
