package linkedservices

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WarehouseLinkedServiceTypeProperties struct {
	ArtifactId                     interface{}  `json:"artifactId"`
	EncryptedCredential            *string      `json:"encryptedCredential,omitempty"`
	Endpoint                       interface{}  `json:"endpoint"`
	ServicePrincipalCredential     SecretBase   `json:"servicePrincipalCredential"`
	ServicePrincipalCredentialType *interface{} `json:"servicePrincipalCredentialType,omitempty"`
	ServicePrincipalId             *interface{} `json:"servicePrincipalId,omitempty"`
	ServicePrincipalKey            SecretBase   `json:"servicePrincipalKey"`
	Tenant                         *interface{} `json:"tenant,omitempty"`
	WorkspaceId                    *interface{} `json:"workspaceId,omitempty"`
}

var _ json.Unmarshaler = &WarehouseLinkedServiceTypeProperties{}

func (s *WarehouseLinkedServiceTypeProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		ArtifactId                     interface{}  `json:"artifactId"`
		EncryptedCredential            *string      `json:"encryptedCredential,omitempty"`
		Endpoint                       interface{}  `json:"endpoint"`
		ServicePrincipalCredentialType *interface{} `json:"servicePrincipalCredentialType,omitempty"`
		ServicePrincipalId             *interface{} `json:"servicePrincipalId,omitempty"`
		Tenant                         *interface{} `json:"tenant,omitempty"`
		WorkspaceId                    *interface{} `json:"workspaceId,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.ArtifactId = decoded.ArtifactId
	s.EncryptedCredential = decoded.EncryptedCredential
	s.Endpoint = decoded.Endpoint
	s.ServicePrincipalCredentialType = decoded.ServicePrincipalCredentialType
	s.ServicePrincipalId = decoded.ServicePrincipalId
	s.Tenant = decoded.Tenant
	s.WorkspaceId = decoded.WorkspaceId

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling WarehouseLinkedServiceTypeProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["servicePrincipalCredential"]; ok {
		impl, err := UnmarshalSecretBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'ServicePrincipalCredential' for 'WarehouseLinkedServiceTypeProperties': %+v", err)
		}
		s.ServicePrincipalCredential = impl
	}

	if v, ok := temp["servicePrincipalKey"]; ok {
		impl, err := UnmarshalSecretBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'ServicePrincipalKey' for 'WarehouseLinkedServiceTypeProperties': %+v", err)
		}
		s.ServicePrincipalKey = impl
	}

	return nil
}
