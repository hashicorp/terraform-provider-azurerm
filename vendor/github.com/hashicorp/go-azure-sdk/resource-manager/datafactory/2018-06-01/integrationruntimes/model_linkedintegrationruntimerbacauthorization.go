package integrationruntimes

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ LinkedIntegrationRuntimeType = LinkedIntegrationRuntimeRbacAuthorization{}

type LinkedIntegrationRuntimeRbacAuthorization struct {
	Credential *CredentialReference `json:"credential,omitempty"`
	ResourceId string               `json:"resourceId"`

	// Fields inherited from LinkedIntegrationRuntimeType

	AuthorizationType string `json:"authorizationType"`
}

func (s LinkedIntegrationRuntimeRbacAuthorization) LinkedIntegrationRuntimeType() BaseLinkedIntegrationRuntimeTypeImpl {
	return BaseLinkedIntegrationRuntimeTypeImpl{
		AuthorizationType: s.AuthorizationType,
	}
}

var _ json.Marshaler = LinkedIntegrationRuntimeRbacAuthorization{}

func (s LinkedIntegrationRuntimeRbacAuthorization) MarshalJSON() ([]byte, error) {
	type wrapper LinkedIntegrationRuntimeRbacAuthorization
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling LinkedIntegrationRuntimeRbacAuthorization: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling LinkedIntegrationRuntimeRbacAuthorization: %+v", err)
	}

	decoded["authorizationType"] = "RBAC"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling LinkedIntegrationRuntimeRbacAuthorization: %+v", err)
	}

	return encoded, nil
}
