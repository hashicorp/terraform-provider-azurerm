package v2workspaceconnectionresource

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ WorkspaceConnectionPropertiesV2 = ManagedIdentityAuthTypeWorkspaceConnectionProperties{}

type ManagedIdentityAuthTypeWorkspaceConnectionProperties struct {
	Credentials *WorkspaceConnectionManagedIdentity `json:"credentials,omitempty"`

	// Fields inherited from WorkspaceConnectionPropertiesV2
	Category    *ConnectionCategory `json:"category,omitempty"`
	Target      *string             `json:"target,omitempty"`
	Value       *string             `json:"value,omitempty"`
	ValueFormat *ValueFormat        `json:"valueFormat,omitempty"`
}

var _ json.Marshaler = ManagedIdentityAuthTypeWorkspaceConnectionProperties{}

func (s ManagedIdentityAuthTypeWorkspaceConnectionProperties) MarshalJSON() ([]byte, error) {
	type wrapper ManagedIdentityAuthTypeWorkspaceConnectionProperties
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ManagedIdentityAuthTypeWorkspaceConnectionProperties: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ManagedIdentityAuthTypeWorkspaceConnectionProperties: %+v", err)
	}
	decoded["authType"] = "ManagedIdentity"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ManagedIdentityAuthTypeWorkspaceConnectionProperties: %+v", err)
	}

	return encoded, nil
}
