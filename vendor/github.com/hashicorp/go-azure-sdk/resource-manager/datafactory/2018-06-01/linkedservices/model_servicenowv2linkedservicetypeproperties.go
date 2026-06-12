package linkedservices

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServiceNowV2LinkedServiceTypeProperties struct {
	AuthenticationType  ServiceNowV2AuthenticationType `json:"authenticationType"`
	ClientId            *interface{}                   `json:"clientId,omitempty"`
	ClientSecret        SecretBase                     `json:"clientSecret"`
	EncryptedCredential *string                        `json:"encryptedCredential,omitempty"`
	Endpoint            interface{}                    `json:"endpoint"`
	GrantType           *interface{}                   `json:"grantType,omitempty"`
	Password            SecretBase                     `json:"password"`
	Username            *interface{}                   `json:"username,omitempty"`
}

var _ json.Unmarshaler = &ServiceNowV2LinkedServiceTypeProperties{}

func (s *ServiceNowV2LinkedServiceTypeProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		AuthenticationType  ServiceNowV2AuthenticationType `json:"authenticationType"`
		ClientId            *interface{}                   `json:"clientId,omitempty"`
		EncryptedCredential *string                        `json:"encryptedCredential,omitempty"`
		Endpoint            interface{}                    `json:"endpoint"`
		GrantType           *interface{}                   `json:"grantType,omitempty"`
		Username            *interface{}                   `json:"username,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.AuthenticationType = decoded.AuthenticationType
	s.ClientId = decoded.ClientId
	s.EncryptedCredential = decoded.EncryptedCredential
	s.Endpoint = decoded.Endpoint
	s.GrantType = decoded.GrantType
	s.Username = decoded.Username

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling ServiceNowV2LinkedServiceTypeProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["clientSecret"]; ok {
		impl, err := UnmarshalSecretBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'ClientSecret' for 'ServiceNowV2LinkedServiceTypeProperties': %+v", err)
		}
		s.ClientSecret = impl
	}

	if v, ok := temp["password"]; ok {
		impl, err := UnmarshalSecretBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Password' for 'ServiceNowV2LinkedServiceTypeProperties': %+v", err)
		}
		s.Password = impl
	}

	return nil
}
