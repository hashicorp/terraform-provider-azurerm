package linkedservices

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type HTTPLinkedServiceTypeProperties struct {
	AuthHeaders                       *map[string]string      `json:"authHeaders,omitempty"`
	AuthenticationType                *HTTPAuthenticationType `json:"authenticationType,omitempty"`
	CertThumbprint                    *interface{}            `json:"certThumbprint,omitempty"`
	EmbeddedCertData                  *interface{}            `json:"embeddedCertData,omitempty"`
	EnableServerCertificateValidation *bool                   `json:"enableServerCertificateValidation,omitempty"`
	EncryptedCredential               *string                 `json:"encryptedCredential,omitempty"`
	Password                          SecretBase              `json:"password"`
	Url                               interface{}             `json:"url"`
	UserName                          *interface{}            `json:"userName,omitempty"`
}

var _ json.Unmarshaler = &HTTPLinkedServiceTypeProperties{}

func (s *HTTPLinkedServiceTypeProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		AuthHeaders                       *map[string]string      `json:"authHeaders,omitempty"`
		AuthenticationType                *HTTPAuthenticationType `json:"authenticationType,omitempty"`
		CertThumbprint                    *interface{}            `json:"certThumbprint,omitempty"`
		EmbeddedCertData                  *interface{}            `json:"embeddedCertData,omitempty"`
		EnableServerCertificateValidation *bool                   `json:"enableServerCertificateValidation,omitempty"`
		EncryptedCredential               *string                 `json:"encryptedCredential,omitempty"`
		Url                               interface{}             `json:"url"`
		UserName                          *interface{}            `json:"userName,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.AuthHeaders = decoded.AuthHeaders
	s.AuthenticationType = decoded.AuthenticationType
	s.CertThumbprint = decoded.CertThumbprint
	s.EmbeddedCertData = decoded.EmbeddedCertData
	s.EnableServerCertificateValidation = decoded.EnableServerCertificateValidation
	s.EncryptedCredential = decoded.EncryptedCredential
	s.Url = decoded.Url
	s.UserName = decoded.UserName

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling HTTPLinkedServiceTypeProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["password"]; ok {
		impl, err := UnmarshalSecretBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Password' for 'HTTPLinkedServiceTypeProperties': %+v", err)
		}
		s.Password = impl
	}

	return nil
}
