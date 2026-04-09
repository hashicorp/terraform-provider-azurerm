package servicelinker

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ SecretInfoBase = KeyVaultSecretReferenceSecretInfo{}

type KeyVaultSecretReferenceSecretInfo struct {
	Name    *string `json:"name,omitempty"`
	Version *string `json:"version,omitempty"`

	// Fields inherited from SecretInfoBase

	SecretType SecretType `json:"secretType"`
}

func (s KeyVaultSecretReferenceSecretInfo) SecretInfoBase() BaseSecretInfoBaseImpl {
	return BaseSecretInfoBaseImpl{
		SecretType: s.SecretType,
	}
}

var _ json.Marshaler = KeyVaultSecretReferenceSecretInfo{}

func (s KeyVaultSecretReferenceSecretInfo) MarshalJSON() ([]byte, error) {
	type wrapper KeyVaultSecretReferenceSecretInfo
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling KeyVaultSecretReferenceSecretInfo: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling KeyVaultSecretReferenceSecretInfo: %+v", err)
	}

	decoded["secretType"] = "keyVaultSecretReference"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling KeyVaultSecretReferenceSecretInfo: %+v", err)
	}

	return encoded, nil
}
