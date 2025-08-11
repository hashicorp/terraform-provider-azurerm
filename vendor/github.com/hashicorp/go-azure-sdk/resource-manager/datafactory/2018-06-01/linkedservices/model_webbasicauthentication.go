package linkedservices

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ WebLinkedServiceTypeProperties = WebBasicAuthentication{}

type WebBasicAuthentication struct {
	Password SecretBase  `json:"password"`
	Username interface{} `json:"username"`

	// Fields inherited from WebLinkedServiceTypeProperties

	AuthenticationType WebAuthenticationType `json:"authenticationType"`
	Url                interface{}           `json:"url"`
}

func (s WebBasicAuthentication) WebLinkedServiceTypeProperties() BaseWebLinkedServiceTypePropertiesImpl {
	return BaseWebLinkedServiceTypePropertiesImpl{
		AuthenticationType: s.AuthenticationType,
		Url:                s.Url,
	}
}

var _ json.Marshaler = WebBasicAuthentication{}

func (s WebBasicAuthentication) MarshalJSON() ([]byte, error) {
	type wrapper WebBasicAuthentication
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling WebBasicAuthentication: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling WebBasicAuthentication: %+v", err)
	}

	decoded["authenticationType"] = "Basic"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling WebBasicAuthentication: %+v", err)
	}

	return encoded, nil
}

var _ json.Unmarshaler = &WebBasicAuthentication{}

func (s *WebBasicAuthentication) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		Username           interface{}           `json:"username"`
		AuthenticationType WebAuthenticationType `json:"authenticationType"`
		Url                interface{}           `json:"url"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.Username = decoded.Username
	s.AuthenticationType = decoded.AuthenticationType
	s.Url = decoded.Url

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling WebBasicAuthentication into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["password"]; ok {
		impl, err := UnmarshalSecretBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Password' for 'WebBasicAuthentication': %+v", err)
		}
		s.Password = impl
	}

	return nil
}
