package linkedservices

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ SecretBase = AzureKeyVaultSecretReference{}

type AzureKeyVaultSecretReference struct {
	SecretName    interface{}            `json:"secretName"`
	SecretVersion *interface{}           `json:"secretVersion,omitempty"`
	Store         LinkedServiceReference `json:"store"`

	// Fields inherited from SecretBase

	Type string `json:"type"`
}

func (s AzureKeyVaultSecretReference) SecretBase() BaseSecretBaseImpl {
	return BaseSecretBaseImpl{
		Type: s.Type,
	}
}

var _ json.Marshaler = AzureKeyVaultSecretReference{}

func (s AzureKeyVaultSecretReference) MarshalJSON() ([]byte, error) {
	type wrapper AzureKeyVaultSecretReference
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AzureKeyVaultSecretReference: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AzureKeyVaultSecretReference: %+v", err)
	}

	decoded["type"] = "AzureKeyVaultSecret"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AzureKeyVaultSecretReference: %+v", err)
	}

	return encoded, nil
}
