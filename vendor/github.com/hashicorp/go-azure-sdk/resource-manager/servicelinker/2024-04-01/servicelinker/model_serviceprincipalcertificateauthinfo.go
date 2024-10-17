package servicelinker

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ AuthInfoBase = ServicePrincipalCertificateAuthInfo{}

type ServicePrincipalCertificateAuthInfo struct {
	Certificate            string                  `json:"certificate"`
	ClientId               string                  `json:"clientId"`
	DeleteOrUpdateBehavior *DeleteOrUpdateBehavior `json:"deleteOrUpdateBehavior,omitempty"`
	PrincipalId            string                  `json:"principalId"`
	Roles                  *[]string               `json:"roles,omitempty"`

	// Fields inherited from AuthInfoBase

	AuthMode *AuthMode `json:"authMode,omitempty"`
	AuthType AuthType  `json:"authType"`
}

func (s ServicePrincipalCertificateAuthInfo) AuthInfoBase() BaseAuthInfoBaseImpl {
	return BaseAuthInfoBaseImpl{
		AuthMode: s.AuthMode,
		AuthType: s.AuthType,
	}
}

var _ json.Marshaler = ServicePrincipalCertificateAuthInfo{}

func (s ServicePrincipalCertificateAuthInfo) MarshalJSON() ([]byte, error) {
	type wrapper ServicePrincipalCertificateAuthInfo
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ServicePrincipalCertificateAuthInfo: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ServicePrincipalCertificateAuthInfo: %+v", err)
	}

	decoded["authType"] = "servicePrincipalCertificate"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ServicePrincipalCertificateAuthInfo: %+v", err)
	}

	return encoded, nil
}
