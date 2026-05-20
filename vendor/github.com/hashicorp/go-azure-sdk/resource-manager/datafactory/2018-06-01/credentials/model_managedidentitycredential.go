package credentials

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ Credential = ManagedIdentityCredential{}

type ManagedIdentityCredential struct {
	TypeProperties *ManagedIdentityTypeProperties `json:"typeProperties,omitempty"`

	// Fields inherited from Credential

	Annotations *[]interface{} `json:"annotations,omitempty"`
	Description *string        `json:"description,omitempty"`
	Type        string         `json:"type"`
}

func (s ManagedIdentityCredential) Credential() BaseCredentialImpl {
	return BaseCredentialImpl{
		Annotations: s.Annotations,
		Description: s.Description,
		Type:        s.Type,
	}
}

var _ json.Marshaler = ManagedIdentityCredential{}

func (s ManagedIdentityCredential) MarshalJSON() ([]byte, error) {
	type wrapper ManagedIdentityCredential
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ManagedIdentityCredential: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ManagedIdentityCredential: %+v", err)
	}

	decoded["type"] = "ManagedIdentity"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ManagedIdentityCredential: %+v", err)
	}

	return encoded, nil
}
