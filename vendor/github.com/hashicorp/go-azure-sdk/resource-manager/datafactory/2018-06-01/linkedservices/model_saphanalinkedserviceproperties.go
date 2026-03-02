package linkedservices

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SapHanaLinkedServiceProperties struct {
	AuthenticationType  *SapHanaAuthenticationType `json:"authenticationType,omitempty"`
	ConnectionString    *interface{}               `json:"connectionString,omitempty"`
	EncryptedCredential *string                    `json:"encryptedCredential,omitempty"`
	Password            SecretBase                 `json:"password"`
	Server              *interface{}               `json:"server,omitempty"`
	UserName            *interface{}               `json:"userName,omitempty"`
}

var _ json.Unmarshaler = &SapHanaLinkedServiceProperties{}

func (s *SapHanaLinkedServiceProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		AuthenticationType  *SapHanaAuthenticationType `json:"authenticationType,omitempty"`
		ConnectionString    *interface{}               `json:"connectionString,omitempty"`
		EncryptedCredential *string                    `json:"encryptedCredential,omitempty"`
		Server              *interface{}               `json:"server,omitempty"`
		UserName            *interface{}               `json:"userName,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.AuthenticationType = decoded.AuthenticationType
	s.ConnectionString = decoded.ConnectionString
	s.EncryptedCredential = decoded.EncryptedCredential
	s.Server = decoded.Server
	s.UserName = decoded.UserName

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling SapHanaLinkedServiceProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["password"]; ok {
		impl, err := UnmarshalSecretBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Password' for 'SapHanaLinkedServiceProperties': %+v", err)
		}
		s.Password = impl
	}

	return nil
}
