package tasks

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ TaskStepUpdateParameters = DockerBuildStepUpdateParameters{}

type DockerBuildStepUpdateParameters struct {
	Arguments      *[]Argument `json:"arguments,omitempty"`
	DockerFilePath *string     `json:"dockerFilePath,omitempty"`
	ImageNames     *[]string   `json:"imageNames,omitempty"`
	IsPushEnabled  *bool       `json:"isPushEnabled,omitempty"`
	NoCache        *bool       `json:"noCache,omitempty"`
	Target         *string     `json:"target,omitempty"`

	// Fields inherited from TaskStepUpdateParameters
	ContextAccessToken *string `json:"contextAccessToken,omitempty"`
	ContextPath        *string `json:"contextPath,omitempty"`
}

var _ json.Marshaler = DockerBuildStepUpdateParameters{}

func (s DockerBuildStepUpdateParameters) MarshalJSON() ([]byte, error) {
	type wrapper DockerBuildStepUpdateParameters
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling DockerBuildStepUpdateParameters: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling DockerBuildStepUpdateParameters: %+v", err)
	}
	decoded["type"] = "Docker"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling DockerBuildStepUpdateParameters: %+v", err)
	}

	return encoded, nil
}
