package linkedservices

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServiceNowLinkedServiceTypeProperties struct {
	AuthenticationType    ServiceNowAuthenticationType `json:"authenticationType"`
	ClientId              *interface{}                 `json:"clientId,omitempty"`
	ClientSecret          SecretBase                   `json:"clientSecret"`
	EncryptedCredential   *string                      `json:"encryptedCredential,omitempty"`
	Endpoint              interface{}                  `json:"endpoint"`
	Password              SecretBase                   `json:"password"`
	UseEncryptedEndpoints *bool                        `json:"useEncryptedEndpoints,omitempty"`
	UseHostVerification   *bool                        `json:"useHostVerification,omitempty"`
	UsePeerVerification   *bool                        `json:"usePeerVerification,omitempty"`
	Username              *interface{}                 `json:"username,omitempty"`
}

var _ json.Unmarshaler = &ServiceNowLinkedServiceTypeProperties{}

func (s *ServiceNowLinkedServiceTypeProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		AuthenticationType    ServiceNowAuthenticationType `json:"authenticationType"`
		ClientId              *interface{}                 `json:"clientId,omitempty"`
		EncryptedCredential   *string                      `json:"encryptedCredential,omitempty"`
		Endpoint              interface{}                  `json:"endpoint"`
		UseEncryptedEndpoints *bool                        `json:"useEncryptedEndpoints,omitempty"`
		UseHostVerification   *bool                        `json:"useHostVerification,omitempty"`
		UsePeerVerification   *bool                        `json:"usePeerVerification,omitempty"`
		Username              *interface{}                 `json:"username,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.AuthenticationType = decoded.AuthenticationType
	s.ClientId = decoded.ClientId
	s.EncryptedCredential = decoded.EncryptedCredential
	s.Endpoint = decoded.Endpoint
	s.UseEncryptedEndpoints = decoded.UseEncryptedEndpoints
	s.UseHostVerification = decoded.UseHostVerification
	s.UsePeerVerification = decoded.UsePeerVerification
	s.Username = decoded.Username

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling ServiceNowLinkedServiceTypeProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["clientSecret"]; ok {
		impl, err := UnmarshalSecretBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'ClientSecret' for 'ServiceNowLinkedServiceTypeProperties': %+v", err)
		}
		s.ClientSecret = impl
	}

	if v, ok := temp["password"]; ok {
		impl, err := UnmarshalSecretBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Password' for 'ServiceNowLinkedServiceTypeProperties': %+v", err)
		}
		s.Password = impl
	}

	return nil
}
