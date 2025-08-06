package tasks

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ TaskStepProperties = DockerBuildStep{}

type DockerBuildStep struct {
	Arguments      *[]Argument `json:"arguments,omitempty"`
	DockerFilePath string      `json:"dockerFilePath"`
	ImageNames     *[]string   `json:"imageNames,omitempty"`
	IsPushEnabled  *bool       `json:"isPushEnabled,omitempty"`
	NoCache        *bool       `json:"noCache,omitempty"`
	Target         *string     `json:"target,omitempty"`

	// Fields inherited from TaskStepProperties

	BaseImageDependencies *[]BaseImageDependency `json:"baseImageDependencies,omitempty"`
	ContextAccessToken    *string                `json:"contextAccessToken,omitempty"`
	ContextPath           *string                `json:"contextPath,omitempty"`
	Type                  StepType               `json:"type"`
}

func (s DockerBuildStep) TaskStepProperties() BaseTaskStepPropertiesImpl {
	return BaseTaskStepPropertiesImpl{
		BaseImageDependencies: s.BaseImageDependencies,
		ContextAccessToken:    s.ContextAccessToken,
		ContextPath:           s.ContextPath,
		Type:                  s.Type,
	}
}

var _ json.Marshaler = DockerBuildStep{}

func (s DockerBuildStep) MarshalJSON() ([]byte, error) {
	type wrapper DockerBuildStep
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling DockerBuildStep: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling DockerBuildStep: %+v", err)
	}

	decoded["type"] = "Docker"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling DockerBuildStep: %+v", err)
	}

	return encoded, nil
}
