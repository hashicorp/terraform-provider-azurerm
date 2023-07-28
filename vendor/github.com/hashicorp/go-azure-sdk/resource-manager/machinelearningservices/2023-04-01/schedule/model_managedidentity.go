package schedule

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ IdentityConfiguration = ManagedIdentity{}

type ManagedIdentity struct {
	ClientId   *string `json:"clientId,omitempty"`
	ObjectId   *string `json:"objectId,omitempty"`
	ResourceId *string `json:"resourceId,omitempty"`

	// Fields inherited from IdentityConfiguration
}

var _ json.Marshaler = ManagedIdentity{}

func (s ManagedIdentity) MarshalJSON() ([]byte, error) {
	type wrapper ManagedIdentity
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ManagedIdentity: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ManagedIdentity: %+v", err)
	}
	decoded["identityType"] = "Managed"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ManagedIdentity: %+v", err)
	}

	return encoded, nil
}
