package v2workspaceconnectionresource

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ WorkspaceConnectionPropertiesV2 = PATAuthTypeWorkspaceConnectionProperties{}

type PATAuthTypeWorkspaceConnectionProperties struct {
	Credentials *WorkspaceConnectionPersonalAccessToken `json:"credentials,omitempty"`

	// Fields inherited from WorkspaceConnectionPropertiesV2
	Category    *ConnectionCategory `json:"category,omitempty"`
	Target      *string             `json:"target,omitempty"`
	Value       *string             `json:"value,omitempty"`
	ValueFormat *ValueFormat        `json:"valueFormat,omitempty"`
}

var _ json.Marshaler = PATAuthTypeWorkspaceConnectionProperties{}

func (s PATAuthTypeWorkspaceConnectionProperties) MarshalJSON() ([]byte, error) {
	type wrapper PATAuthTypeWorkspaceConnectionProperties
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling PATAuthTypeWorkspaceConnectionProperties: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling PATAuthTypeWorkspaceConnectionProperties: %+v", err)
	}
	decoded["authType"] = "PAT"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling PATAuthTypeWorkspaceConnectionProperties: %+v", err)
	}

	return encoded, nil
}
