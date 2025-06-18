package linkedservices

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SybaseLinkedServiceTypeProperties struct {
	AuthenticationType  *SybaseAuthenticationType `json:"authenticationType,omitempty"`
	Database            interface{}               `json:"database"`
	EncryptedCredential *string                   `json:"encryptedCredential,omitempty"`
	Password            SecretBase                `json:"password"`
	Schema              *interface{}              `json:"schema,omitempty"`
	Server              interface{}               `json:"server"`
	Username            *interface{}              `json:"username,omitempty"`
}

var _ json.Unmarshaler = &SybaseLinkedServiceTypeProperties{}

func (s *SybaseLinkedServiceTypeProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		AuthenticationType  *SybaseAuthenticationType `json:"authenticationType,omitempty"`
		Database            interface{}               `json:"database"`
		EncryptedCredential *string                   `json:"encryptedCredential,omitempty"`
		Schema              *interface{}              `json:"schema,omitempty"`
		Server              interface{}               `json:"server"`
		Username            *interface{}              `json:"username,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.AuthenticationType = decoded.AuthenticationType
	s.Database = decoded.Database
	s.EncryptedCredential = decoded.EncryptedCredential
	s.Schema = decoded.Schema
	s.Server = decoded.Server
	s.Username = decoded.Username

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling SybaseLinkedServiceTypeProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["password"]; ok {
		impl, err := UnmarshalSecretBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Password' for 'SybaseLinkedServiceTypeProperties': %+v", err)
		}
		s.Password = impl
	}

	return nil
}
