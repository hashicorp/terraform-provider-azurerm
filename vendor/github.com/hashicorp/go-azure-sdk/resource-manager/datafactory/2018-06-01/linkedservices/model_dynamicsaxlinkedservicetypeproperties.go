package linkedservices

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DynamicsAXLinkedServiceTypeProperties struct {
	AadResourceId       interface{} `json:"aadResourceId"`
	EncryptedCredential *string     `json:"encryptedCredential,omitempty"`
	ServicePrincipalId  interface{} `json:"servicePrincipalId"`
	ServicePrincipalKey SecretBase  `json:"servicePrincipalKey"`
	Tenant              interface{} `json:"tenant"`
	Url                 interface{} `json:"url"`
}

var _ json.Unmarshaler = &DynamicsAXLinkedServiceTypeProperties{}

func (s *DynamicsAXLinkedServiceTypeProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		AadResourceId       interface{} `json:"aadResourceId"`
		EncryptedCredential *string     `json:"encryptedCredential,omitempty"`
		ServicePrincipalId  interface{} `json:"servicePrincipalId"`
		Tenant              interface{} `json:"tenant"`
		Url                 interface{} `json:"url"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.AadResourceId = decoded.AadResourceId
	s.EncryptedCredential = decoded.EncryptedCredential
	s.ServicePrincipalId = decoded.ServicePrincipalId
	s.Tenant = decoded.Tenant
	s.Url = decoded.Url

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling DynamicsAXLinkedServiceTypeProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["servicePrincipalKey"]; ok {
		impl, err := UnmarshalSecretBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'ServicePrincipalKey' for 'DynamicsAXLinkedServiceTypeProperties': %+v", err)
		}
		s.ServicePrincipalKey = impl
	}

	return nil
}
