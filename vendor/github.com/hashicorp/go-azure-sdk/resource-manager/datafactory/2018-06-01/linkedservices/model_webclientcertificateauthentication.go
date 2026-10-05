package linkedservices

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ WebLinkedServiceTypeProperties = WebClientCertificateAuthentication{}

type WebClientCertificateAuthentication struct {
	Password SecretBase `json:"password"`
	Pfx      SecretBase `json:"pfx"`

	// Fields inherited from WebLinkedServiceTypeProperties

	AuthenticationType WebAuthenticationType `json:"authenticationType"`
	Url                interface{}           `json:"url"`
}

func (s WebClientCertificateAuthentication) WebLinkedServiceTypeProperties() BaseWebLinkedServiceTypePropertiesImpl {
	return BaseWebLinkedServiceTypePropertiesImpl{
		AuthenticationType: s.AuthenticationType,
		Url:                s.Url,
	}
}

var _ json.Marshaler = WebClientCertificateAuthentication{}

func (s WebClientCertificateAuthentication) MarshalJSON() ([]byte, error) {
	type wrapper WebClientCertificateAuthentication
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling WebClientCertificateAuthentication: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling WebClientCertificateAuthentication: %+v", err)
	}

	decoded["authenticationType"] = "ClientCertificate"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling WebClientCertificateAuthentication: %+v", err)
	}

	return encoded, nil
}

var _ json.Unmarshaler = &WebClientCertificateAuthentication{}

func (s *WebClientCertificateAuthentication) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		AuthenticationType WebAuthenticationType `json:"authenticationType"`
		Url                interface{}           `json:"url"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.AuthenticationType = decoded.AuthenticationType
	s.Url = decoded.Url

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling WebClientCertificateAuthentication into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["password"]; ok {
		impl, err := UnmarshalSecretBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Password' for 'WebClientCertificateAuthentication': %+v", err)
		}
		s.Password = impl
	}

	if v, ok := temp["pfx"]; ok {
		impl, err := UnmarshalSecretBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Pfx' for 'WebClientCertificateAuthentication': %+v", err)
		}
		s.Pfx = impl
	}

	return nil
}
