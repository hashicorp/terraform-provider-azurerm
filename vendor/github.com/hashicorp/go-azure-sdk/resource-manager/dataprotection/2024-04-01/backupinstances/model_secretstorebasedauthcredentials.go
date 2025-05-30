package backupinstances

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ AuthCredentials = SecretStoreBasedAuthCredentials{}

type SecretStoreBasedAuthCredentials struct {
	SecretStoreResource *SecretStoreResource `json:"secretStoreResource,omitempty"`

	// Fields inherited from AuthCredentials

	ObjectType string `json:"objectType"`
}

func (s SecretStoreBasedAuthCredentials) AuthCredentials() BaseAuthCredentialsImpl {
	return BaseAuthCredentialsImpl{
		ObjectType: s.ObjectType,
	}
}

var _ json.Marshaler = SecretStoreBasedAuthCredentials{}

func (s SecretStoreBasedAuthCredentials) MarshalJSON() ([]byte, error) {
	type wrapper SecretStoreBasedAuthCredentials
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling SecretStoreBasedAuthCredentials: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling SecretStoreBasedAuthCredentials: %+v", err)
	}

	decoded["objectType"] = "SecretStoreBasedAuthCredentials"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling SecretStoreBasedAuthCredentials: %+v", err)
	}

	return encoded, nil
}
