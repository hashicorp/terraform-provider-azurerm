package managedidentities

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type UserAssignedIdentityProperties struct {
	ClientId    *string `json:"clientId,omitempty"`
	PrincipalId *string `json:"principalId,omitempty"`
	TenantId    *string `json:"tenantId,omitempty"`
}

var _ json.Marshaler = UserAssignedIdentityProperties{}

func (s UserAssignedIdentityProperties) MarshalJSON() ([]byte, error) {
	type wrapper UserAssignedIdentityProperties
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling UserAssignedIdentityProperties: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling UserAssignedIdentityProperties: %+v", err)
	}

	delete(decoded, "clientId")
	delete(decoded, "principalId")
	delete(decoded, "tenantId")

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling UserAssignedIdentityProperties: %+v", err)
	}

	return encoded, nil
}
