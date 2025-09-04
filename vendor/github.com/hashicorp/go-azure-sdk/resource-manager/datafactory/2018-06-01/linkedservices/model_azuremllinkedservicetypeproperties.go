package linkedservices

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureMLLinkedServiceTypeProperties struct {
	ApiKey                 SecretBase   `json:"apiKey"`
	Authentication         *interface{} `json:"authentication,omitempty"`
	EncryptedCredential    *string      `json:"encryptedCredential,omitempty"`
	MlEndpoint             interface{}  `json:"mlEndpoint"`
	ServicePrincipalId     *interface{} `json:"servicePrincipalId,omitempty"`
	ServicePrincipalKey    SecretBase   `json:"servicePrincipalKey"`
	Tenant                 *interface{} `json:"tenant,omitempty"`
	UpdateResourceEndpoint *interface{} `json:"updateResourceEndpoint,omitempty"`
}

var _ json.Unmarshaler = &AzureMLLinkedServiceTypeProperties{}

func (s *AzureMLLinkedServiceTypeProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		Authentication         *interface{} `json:"authentication,omitempty"`
		EncryptedCredential    *string      `json:"encryptedCredential,omitempty"`
		MlEndpoint             interface{}  `json:"mlEndpoint"`
		ServicePrincipalId     *interface{} `json:"servicePrincipalId,omitempty"`
		Tenant                 *interface{} `json:"tenant,omitempty"`
		UpdateResourceEndpoint *interface{} `json:"updateResourceEndpoint,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.Authentication = decoded.Authentication
	s.EncryptedCredential = decoded.EncryptedCredential
	s.MlEndpoint = decoded.MlEndpoint
	s.ServicePrincipalId = decoded.ServicePrincipalId
	s.Tenant = decoded.Tenant
	s.UpdateResourceEndpoint = decoded.UpdateResourceEndpoint

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling AzureMLLinkedServiceTypeProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["apiKey"]; ok {
		impl, err := UnmarshalSecretBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'ApiKey' for 'AzureMLLinkedServiceTypeProperties': %+v", err)
		}
		s.ApiKey = impl
	}

	if v, ok := temp["servicePrincipalKey"]; ok {
		impl, err := UnmarshalSecretBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'ServicePrincipalKey' for 'AzureMLLinkedServiceTypeProperties': %+v", err)
		}
		s.ServicePrincipalKey = impl
	}

	return nil
}
