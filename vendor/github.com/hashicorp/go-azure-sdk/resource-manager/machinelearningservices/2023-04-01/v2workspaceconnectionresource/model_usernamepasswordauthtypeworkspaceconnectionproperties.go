package v2workspaceconnectionresource

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ WorkspaceConnectionPropertiesV2 = UsernamePasswordAuthTypeWorkspaceConnectionProperties{}

type UsernamePasswordAuthTypeWorkspaceConnectionProperties struct {
	Credentials *WorkspaceConnectionUsernamePassword `json:"credentials,omitempty"`

	// Fields inherited from WorkspaceConnectionPropertiesV2
	Category    *ConnectionCategory `json:"category,omitempty"`
	Target      *string             `json:"target,omitempty"`
	Value       *string             `json:"value,omitempty"`
	ValueFormat *ValueFormat        `json:"valueFormat,omitempty"`
}

var _ json.Marshaler = UsernamePasswordAuthTypeWorkspaceConnectionProperties{}

func (s UsernamePasswordAuthTypeWorkspaceConnectionProperties) MarshalJSON() ([]byte, error) {
	type wrapper UsernamePasswordAuthTypeWorkspaceConnectionProperties
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling UsernamePasswordAuthTypeWorkspaceConnectionProperties: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling UsernamePasswordAuthTypeWorkspaceConnectionProperties: %+v", err)
	}
	decoded["authType"] = "UsernamePassword"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling UsernamePasswordAuthTypeWorkspaceConnectionProperties: %+v", err)
	}

	return encoded, nil
}
