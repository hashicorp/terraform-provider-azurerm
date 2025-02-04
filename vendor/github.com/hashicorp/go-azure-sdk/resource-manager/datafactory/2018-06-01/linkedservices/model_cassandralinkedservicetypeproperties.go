package linkedservices

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CassandraLinkedServiceTypeProperties struct {
	AuthenticationType  *string    `json:"authenticationType,omitempty"`
	EncryptedCredential *string    `json:"encryptedCredential,omitempty"`
	Host                string     `json:"host"`
	Password            SecretBase `json:"password"`
	Port                *int64     `json:"port,omitempty"`
	Username            *string    `json:"username,omitempty"`
}

var _ json.Unmarshaler = &CassandraLinkedServiceTypeProperties{}

func (s *CassandraLinkedServiceTypeProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		AuthenticationType  *string `json:"authenticationType,omitempty"`
		EncryptedCredential *string `json:"encryptedCredential,omitempty"`
		Host                string  `json:"host"`
		Port                *int64  `json:"port,omitempty"`
		Username            *string `json:"username,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.AuthenticationType = decoded.AuthenticationType
	s.EncryptedCredential = decoded.EncryptedCredential
	s.Host = decoded.Host
	s.Port = decoded.Port
	s.Username = decoded.Username

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling CassandraLinkedServiceTypeProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["password"]; ok {
		impl, err := UnmarshalSecretBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Password' for 'CassandraLinkedServiceTypeProperties': %+v", err)
		}
		s.Password = impl
	}

	return nil
}
