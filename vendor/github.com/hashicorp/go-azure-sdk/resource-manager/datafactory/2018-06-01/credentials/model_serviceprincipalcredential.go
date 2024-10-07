package credentials

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ Credential = ServicePrincipalCredential{}

type ServicePrincipalCredential struct {
	TypeProperties ServicePrincipalCredentialTypeProperties `json:"typeProperties"`

	// Fields inherited from Credential

	Annotations *[]interface{} `json:"annotations,omitempty"`
	Description *string        `json:"description,omitempty"`
	Type        string         `json:"type"`
}

func (s ServicePrincipalCredential) Credential() BaseCredentialImpl {
	return BaseCredentialImpl{
		Annotations: s.Annotations,
		Description: s.Description,
		Type:        s.Type,
	}
}

var _ json.Marshaler = ServicePrincipalCredential{}

func (s ServicePrincipalCredential) MarshalJSON() ([]byte, error) {
	type wrapper ServicePrincipalCredential
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ServicePrincipalCredential: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ServicePrincipalCredential: %+v", err)
	}

	decoded["type"] = "ServicePrincipal"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ServicePrincipalCredential: %+v", err)
	}

	return encoded, nil
}
