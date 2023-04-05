package taskruns

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ RunRequest = FileTaskRunRequest{}

type FileTaskRunRequest struct {
	AgentConfiguration *AgentProperties   `json:"agentConfiguration,omitempty"`
	Credentials        *Credentials       `json:"credentials,omitempty"`
	Platform           PlatformProperties `json:"platform"`
	SourceLocation     *string            `json:"sourceLocation,omitempty"`
	TaskFilePath       string             `json:"taskFilePath"`
	Timeout            *int64             `json:"timeout,omitempty"`
	Values             *[]SetValue        `json:"values,omitempty"`
	ValuesFilePath     *string            `json:"valuesFilePath,omitempty"`

	// Fields inherited from RunRequest
	AgentPoolName    *string `json:"agentPoolName,omitempty"`
	IsArchiveEnabled *bool   `json:"isArchiveEnabled,omitempty"`
	LogTemplate      *string `json:"logTemplate,omitempty"`
}

var _ json.Marshaler = FileTaskRunRequest{}

func (s FileTaskRunRequest) MarshalJSON() ([]byte, error) {
	type wrapper FileTaskRunRequest
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling FileTaskRunRequest: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling FileTaskRunRequest: %+v", err)
	}
	decoded["type"] = "FileTaskRunRequest"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling FileTaskRunRequest: %+v", err)
	}

	return encoded, nil
}
