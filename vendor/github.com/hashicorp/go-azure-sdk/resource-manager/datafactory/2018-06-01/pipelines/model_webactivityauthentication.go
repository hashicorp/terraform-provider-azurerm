package pipelines

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WebActivityAuthentication struct {
	Credential *CredentialReference `json:"credential,omitempty"`
	Password   SecretBase           `json:"password"`
	Pfx        SecretBase           `json:"pfx"`
	Resource   *interface{}         `json:"resource,omitempty"`
	Type       *string              `json:"type,omitempty"`
	UserTenant *interface{}         `json:"userTenant,omitempty"`
	Username   *interface{}         `json:"username,omitempty"`
}

var _ json.Unmarshaler = &WebActivityAuthentication{}

func (s *WebActivityAuthentication) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		Credential *CredentialReference `json:"credential,omitempty"`
		Resource   *interface{}         `json:"resource,omitempty"`
		Type       *string              `json:"type,omitempty"`
		UserTenant *interface{}         `json:"userTenant,omitempty"`
		Username   *interface{}         `json:"username,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.Credential = decoded.Credential
	s.Resource = decoded.Resource
	s.Type = decoded.Type
	s.UserTenant = decoded.UserTenant
	s.Username = decoded.Username

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling WebActivityAuthentication into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["password"]; ok {
		impl, err := UnmarshalSecretBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Password' for 'WebActivityAuthentication': %+v", err)
		}
		s.Password = impl
	}

	if v, ok := temp["pfx"]; ok {
		impl, err := UnmarshalSecretBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Pfx' for 'WebActivityAuthentication': %+v", err)
		}
		s.Pfx = impl
	}

	return nil
}
