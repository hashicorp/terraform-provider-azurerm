package servicelinker

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ AzureResourcePropertiesBase = AzureKeyVaultProperties{}

type AzureKeyVaultProperties struct {
	ConnectAsKubernetesCsiDriver *bool `json:"connectAsKubernetesCsiDriver,omitempty"`

	// Fields inherited from AzureResourcePropertiesBase

	Type AzureResourceType `json:"type"`
}

func (s AzureKeyVaultProperties) AzureResourcePropertiesBase() BaseAzureResourcePropertiesBaseImpl {
	return BaseAzureResourcePropertiesBaseImpl{
		Type: s.Type,
	}
}

var _ json.Marshaler = AzureKeyVaultProperties{}

func (s AzureKeyVaultProperties) MarshalJSON() ([]byte, error) {
	type wrapper AzureKeyVaultProperties
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AzureKeyVaultProperties: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AzureKeyVaultProperties: %+v", err)
	}

	decoded["type"] = "KeyVault"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AzureKeyVaultProperties: %+v", err)
	}

	return encoded, nil
}
