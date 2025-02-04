package linkedservices

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TeradataLinkedServiceTypeProperties struct {
	AuthenticationType  *TeradataAuthenticationType `json:"authenticationType,omitempty"`
	ConnectionString    *string                     `json:"connectionString,omitempty"`
	EncryptedCredential *string                     `json:"encryptedCredential,omitempty"`
	Password            SecretBase                  `json:"password"`
	Server              *string                     `json:"server,omitempty"`
	Username            *string                     `json:"username,omitempty"`
}

var _ json.Unmarshaler = &TeradataLinkedServiceTypeProperties{}

func (s *TeradataLinkedServiceTypeProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		AuthenticationType  *TeradataAuthenticationType `json:"authenticationType,omitempty"`
		ConnectionString    *string                     `json:"connectionString,omitempty"`
		EncryptedCredential *string                     `json:"encryptedCredential,omitempty"`
		Server              *string                     `json:"server,omitempty"`
		Username            *string                     `json:"username,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.AuthenticationType = decoded.AuthenticationType
	s.ConnectionString = decoded.ConnectionString
	s.EncryptedCredential = decoded.EncryptedCredential
	s.Server = decoded.Server
	s.Username = decoded.Username

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling TeradataLinkedServiceTypeProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["password"]; ok {
		impl, err := UnmarshalSecretBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Password' for 'TeradataLinkedServiceTypeProperties': %+v", err)
		}
		s.Password = impl
	}

	return nil
}
