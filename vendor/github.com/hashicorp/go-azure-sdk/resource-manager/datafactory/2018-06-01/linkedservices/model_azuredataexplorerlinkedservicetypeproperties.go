package linkedservices

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureDataExplorerLinkedServiceTypeProperties struct {
	Credential          *CredentialReference `json:"credential,omitempty"`
	Database            interface{}          `json:"database"`
	Endpoint            interface{}          `json:"endpoint"`
	ServicePrincipalId  *interface{}         `json:"servicePrincipalId,omitempty"`
	ServicePrincipalKey SecretBase           `json:"servicePrincipalKey"`
	Tenant              *interface{}         `json:"tenant,omitempty"`
}

var _ json.Unmarshaler = &AzureDataExplorerLinkedServiceTypeProperties{}

func (s *AzureDataExplorerLinkedServiceTypeProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		Credential         *CredentialReference `json:"credential,omitempty"`
		Database           interface{}          `json:"database"`
		Endpoint           interface{}          `json:"endpoint"`
		ServicePrincipalId *interface{}         `json:"servicePrincipalId,omitempty"`
		Tenant             *interface{}         `json:"tenant,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.Credential = decoded.Credential
	s.Database = decoded.Database
	s.Endpoint = decoded.Endpoint
	s.ServicePrincipalId = decoded.ServicePrincipalId
	s.Tenant = decoded.Tenant

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling AzureDataExplorerLinkedServiceTypeProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["servicePrincipalKey"]; ok {
		impl, err := UnmarshalSecretBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'ServicePrincipalKey' for 'AzureDataExplorerLinkedServiceTypeProperties': %+v", err)
		}
		s.ServicePrincipalKey = impl
	}

	return nil
}
