package linkedservices

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LakeHouseLinkedServiceTypeProperties struct {
	ArtifactId                     *interface{} `json:"artifactId,omitempty"`
	EncryptedCredential            *string      `json:"encryptedCredential,omitempty"`
	ServicePrincipalCredential     SecretBase   `json:"servicePrincipalCredential"`
	ServicePrincipalCredentialType *interface{} `json:"servicePrincipalCredentialType,omitempty"`
	ServicePrincipalId             *interface{} `json:"servicePrincipalId,omitempty"`
	ServicePrincipalKey            SecretBase   `json:"servicePrincipalKey"`
	Tenant                         *interface{} `json:"tenant,omitempty"`
	WorkspaceId                    *interface{} `json:"workspaceId,omitempty"`
}

var _ json.Unmarshaler = &LakeHouseLinkedServiceTypeProperties{}

func (s *LakeHouseLinkedServiceTypeProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		ArtifactId                     *interface{} `json:"artifactId,omitempty"`
		EncryptedCredential            *string      `json:"encryptedCredential,omitempty"`
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
	s.ServicePrincipalCredentialType = decoded.ServicePrincipalCredentialType
	s.ServicePrincipalId = decoded.ServicePrincipalId
	s.Tenant = decoded.Tenant
	s.WorkspaceId = decoded.WorkspaceId

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling LakeHouseLinkedServiceTypeProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["servicePrincipalCredential"]; ok {
		impl, err := UnmarshalSecretBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'ServicePrincipalCredential' for 'LakeHouseLinkedServiceTypeProperties': %+v", err)
		}
		s.ServicePrincipalCredential = impl
	}

	if v, ok := temp["servicePrincipalKey"]; ok {
		impl, err := UnmarshalSecretBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'ServicePrincipalKey' for 'LakeHouseLinkedServiceTypeProperties': %+v", err)
		}
		s.ServicePrincipalKey = impl
	}

	return nil
}
