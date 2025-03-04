package linkedservices

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ SecretBase = SecureString{}

type SecureString struct {
	Value string `json:"value"`

	// Fields inherited from SecretBase

	Type string `json:"type"`
}

func (s SecureString) SecretBase() BaseSecretBaseImpl {
	return BaseSecretBaseImpl{
		Type: s.Type,
	}
}

var _ json.Marshaler = SecureString{}

func (s SecureString) MarshalJSON() ([]byte, error) {
	type wrapper SecureString
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling SecureString: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling SecureString: %+v", err)
	}

	decoded["type"] = "SecureString"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling SecureString: %+v", err)
	}

	return encoded, nil
}
