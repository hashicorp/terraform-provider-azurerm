package servicelinker

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ AuthInfoBase = SystemAssignedIdentityAuthInfo{}

type SystemAssignedIdentityAuthInfo struct {
	DeleteOrUpdateBehavior *DeleteOrUpdateBehavior `json:"deleteOrUpdateBehavior,omitempty"`
	Roles                  *[]string               `json:"roles,omitempty"`
	UserName               *string                 `json:"userName,omitempty"`

	// Fields inherited from AuthInfoBase

	AuthMode *AuthMode `json:"authMode,omitempty"`
	AuthType AuthType  `json:"authType"`
}

func (s SystemAssignedIdentityAuthInfo) AuthInfoBase() BaseAuthInfoBaseImpl {
	return BaseAuthInfoBaseImpl{
		AuthMode: s.AuthMode,
		AuthType: s.AuthType,
	}
}

var _ json.Marshaler = SystemAssignedIdentityAuthInfo{}

func (s SystemAssignedIdentityAuthInfo) MarshalJSON() ([]byte, error) {
	type wrapper SystemAssignedIdentityAuthInfo
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling SystemAssignedIdentityAuthInfo: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling SystemAssignedIdentityAuthInfo: %+v", err)
	}

	decoded["authType"] = "systemAssignedIdentity"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling SystemAssignedIdentityAuthInfo: %+v", err)
	}

	return encoded, nil
}
