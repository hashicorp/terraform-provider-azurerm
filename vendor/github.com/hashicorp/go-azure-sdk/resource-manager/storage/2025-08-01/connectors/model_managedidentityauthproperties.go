package connectors

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ StorageConnectorAuthProperties = ManagedIdentityAuthProperties{}

type ManagedIdentityAuthProperties struct {
	IdentityResourceId *string `json:"identityResourceId,omitempty"`

	// Fields inherited from StorageConnectorAuthProperties

	Type StorageConnectorAuthType `json:"type"`
}

func (s ManagedIdentityAuthProperties) StorageConnectorAuthProperties() BaseStorageConnectorAuthPropertiesImpl {
	return BaseStorageConnectorAuthPropertiesImpl{
		Type: s.Type,
	}
}

var _ json.Marshaler = ManagedIdentityAuthProperties{}

func (s ManagedIdentityAuthProperties) MarshalJSON() ([]byte, error) {
	type wrapper ManagedIdentityAuthProperties
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ManagedIdentityAuthProperties: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ManagedIdentityAuthProperties: %+v", err)
	}

	decoded["type"] = "ManagedIdentity"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ManagedIdentityAuthProperties: %+v", err)
	}

	return encoded, nil
}
