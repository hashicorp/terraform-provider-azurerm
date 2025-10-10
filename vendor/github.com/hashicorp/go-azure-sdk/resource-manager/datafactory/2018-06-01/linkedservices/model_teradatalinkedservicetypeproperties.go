package linkedservices

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TeradataLinkedServiceTypeProperties struct {
	AuthenticationType  *TeradataAuthenticationType `json:"authenticationType,omitempty"`
	CharacterSet        *interface{}                `json:"characterSet,omitempty"`
	ConnectionString    *interface{}                `json:"connectionString,omitempty"`
	EncryptedCredential *string                     `json:"encryptedCredential,omitempty"`
	HTTPSPortNumber     *int64                      `json:"httpsPortNumber,omitempty"`
	MaxRespSize         *int64                      `json:"maxRespSize,omitempty"`
	Password            SecretBase                  `json:"password"`
	PortNumber          *int64                      `json:"portNumber,omitempty"`
	Server              *interface{}                `json:"server,omitempty"`
	SslMode             *interface{}                `json:"sslMode,omitempty"`
	UseDataEncryption   *int64                      `json:"useDataEncryption,omitempty"`
	Username            *interface{}                `json:"username,omitempty"`
}

var _ json.Unmarshaler = &TeradataLinkedServiceTypeProperties{}

func (s *TeradataLinkedServiceTypeProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		AuthenticationType  *TeradataAuthenticationType `json:"authenticationType,omitempty"`
		CharacterSet        *interface{}                `json:"characterSet,omitempty"`
		ConnectionString    *interface{}                `json:"connectionString,omitempty"`
		EncryptedCredential *string                     `json:"encryptedCredential,omitempty"`
		HTTPSPortNumber     *int64                      `json:"httpsPortNumber,omitempty"`
		MaxRespSize         *int64                      `json:"maxRespSize,omitempty"`
		PortNumber          *int64                      `json:"portNumber,omitempty"`
		Server              *interface{}                `json:"server,omitempty"`
		SslMode             *interface{}                `json:"sslMode,omitempty"`
		UseDataEncryption   *int64                      `json:"useDataEncryption,omitempty"`
		Username            *interface{}                `json:"username,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.AuthenticationType = decoded.AuthenticationType
	s.CharacterSet = decoded.CharacterSet
	s.ConnectionString = decoded.ConnectionString
	s.EncryptedCredential = decoded.EncryptedCredential
	s.HTTPSPortNumber = decoded.HTTPSPortNumber
	s.MaxRespSize = decoded.MaxRespSize
	s.PortNumber = decoded.PortNumber
	s.Server = decoded.Server
	s.SslMode = decoded.SslMode
	s.UseDataEncryption = decoded.UseDataEncryption
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
