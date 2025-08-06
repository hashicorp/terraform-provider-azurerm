package servicelinker

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ AuthInfoBase = ServicePrincipalSecretAuthInfo{}

type ServicePrincipalSecretAuthInfo struct {
	ClientId               string                  `json:"clientId"`
	DeleteOrUpdateBehavior *DeleteOrUpdateBehavior `json:"deleteOrUpdateBehavior,omitempty"`
	PrincipalId            string                  `json:"principalId"`
	Roles                  *[]string               `json:"roles,omitempty"`
	Secret                 string                  `json:"secret"`
	UserName               *string                 `json:"userName,omitempty"`

	// Fields inherited from AuthInfoBase

	AuthMode *AuthMode `json:"authMode,omitempty"`
	AuthType AuthType  `json:"authType"`
}

func (s ServicePrincipalSecretAuthInfo) AuthInfoBase() BaseAuthInfoBaseImpl {
	return BaseAuthInfoBaseImpl{
		AuthMode: s.AuthMode,
		AuthType: s.AuthType,
	}
}

var _ json.Marshaler = ServicePrincipalSecretAuthInfo{}

func (s ServicePrincipalSecretAuthInfo) MarshalJSON() ([]byte, error) {
	type wrapper ServicePrincipalSecretAuthInfo
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ServicePrincipalSecretAuthInfo: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ServicePrincipalSecretAuthInfo: %+v", err)
	}

	decoded["authType"] = "servicePrincipalSecret"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ServicePrincipalSecretAuthInfo: %+v", err)
	}

	return encoded, nil
}
