package endpoints

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ Credentials = AzureKeyVaultSmbCredentials{}

type AzureKeyVaultSmbCredentials struct {
	PasswordUri *string `json:"passwordUri,omitempty"`
	UsernameUri *string `json:"usernameUri,omitempty"`

	// Fields inherited from Credentials

	Type CredentialType `json:"type"`
}

func (s AzureKeyVaultSmbCredentials) Credentials() BaseCredentialsImpl {
	return BaseCredentialsImpl{
		Type: s.Type,
	}
}

var _ json.Marshaler = AzureKeyVaultSmbCredentials{}

func (s AzureKeyVaultSmbCredentials) MarshalJSON() ([]byte, error) {
	type wrapper AzureKeyVaultSmbCredentials
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AzureKeyVaultSmbCredentials: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AzureKeyVaultSmbCredentials: %+v", err)
	}

	decoded["type"] = "AzureKeyVaultSmb"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AzureKeyVaultSmbCredentials: %+v", err)
	}

	return encoded, nil
}
