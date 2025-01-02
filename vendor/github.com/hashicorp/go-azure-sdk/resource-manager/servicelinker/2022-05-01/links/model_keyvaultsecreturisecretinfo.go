package links

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ SecretInfoBase = KeyVaultSecretUriSecretInfo{}

type KeyVaultSecretUriSecretInfo struct {
	Value *string `json:"value,omitempty"`

	// Fields inherited from SecretInfoBase

	SecretType SecretType `json:"secretType"`
}

func (s KeyVaultSecretUriSecretInfo) SecretInfoBase() BaseSecretInfoBaseImpl {
	return BaseSecretInfoBaseImpl{
		SecretType: s.SecretType,
	}
}

var _ json.Marshaler = KeyVaultSecretUriSecretInfo{}

func (s KeyVaultSecretUriSecretInfo) MarshalJSON() ([]byte, error) {
	type wrapper KeyVaultSecretUriSecretInfo
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling KeyVaultSecretUriSecretInfo: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling KeyVaultSecretUriSecretInfo: %+v", err)
	}

	decoded["secretType"] = "keyVaultSecretUri"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling KeyVaultSecretUriSecretInfo: %+v", err)
	}

	return encoded, nil
}
