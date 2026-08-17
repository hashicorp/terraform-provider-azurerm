package connectors

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ StorageConnectorAuthPropertiesUpdate = ManagedIdentityAuthPropertiesUpdate{}

type ManagedIdentityAuthPropertiesUpdate struct {
	IdentityResourceId *string `json:"identityResourceId,omitempty"`

	// Fields inherited from StorageConnectorAuthPropertiesUpdate

	Type StorageConnectorAuthType `json:"type"`
}

func (s ManagedIdentityAuthPropertiesUpdate) StorageConnectorAuthPropertiesUpdate() BaseStorageConnectorAuthPropertiesUpdateImpl {
	return BaseStorageConnectorAuthPropertiesUpdateImpl{
		Type: s.Type,
	}
}

var _ json.Marshaler = ManagedIdentityAuthPropertiesUpdate{}

func (s ManagedIdentityAuthPropertiesUpdate) MarshalJSON() ([]byte, error) {
	type wrapper ManagedIdentityAuthPropertiesUpdate
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ManagedIdentityAuthPropertiesUpdate: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ManagedIdentityAuthPropertiesUpdate: %+v", err)
	}

	decoded["type"] = "ManagedIdentity"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ManagedIdentityAuthPropertiesUpdate: %+v", err)
	}

	return encoded, nil
}
