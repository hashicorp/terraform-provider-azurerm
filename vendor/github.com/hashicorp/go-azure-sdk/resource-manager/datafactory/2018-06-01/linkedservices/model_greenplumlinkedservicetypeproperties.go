package linkedservices

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GreenplumLinkedServiceTypeProperties struct {
	AuthenticationType  *GreenplumAuthenticationType  `json:"authenticationType,omitempty"`
	CommandTimeout      *int64                        `json:"commandTimeout,omitempty"`
	ConnectionString    *interface{}                  `json:"connectionString,omitempty"`
	ConnectionTimeout   *int64                        `json:"connectionTimeout,omitempty"`
	Database            *interface{}                  `json:"database,omitempty"`
	EncryptedCredential *string                       `json:"encryptedCredential,omitempty"`
	Host                *interface{}                  `json:"host,omitempty"`
	Password            SecretBase                    `json:"password"`
	Port                *int64                        `json:"port,omitempty"`
	Pwd                 *AzureKeyVaultSecretReference `json:"pwd,omitempty"`
	SslMode             *int64                        `json:"sslMode,omitempty"`
	Username            *interface{}                  `json:"username,omitempty"`
}

var _ json.Unmarshaler = &GreenplumLinkedServiceTypeProperties{}

func (s *GreenplumLinkedServiceTypeProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		AuthenticationType  *GreenplumAuthenticationType  `json:"authenticationType,omitempty"`
		CommandTimeout      *int64                        `json:"commandTimeout,omitempty"`
		ConnectionString    *interface{}                  `json:"connectionString,omitempty"`
		ConnectionTimeout   *int64                        `json:"connectionTimeout,omitempty"`
		Database            *interface{}                  `json:"database,omitempty"`
		EncryptedCredential *string                       `json:"encryptedCredential,omitempty"`
		Host                *interface{}                  `json:"host,omitempty"`
		Port                *int64                        `json:"port,omitempty"`
		Pwd                 *AzureKeyVaultSecretReference `json:"pwd,omitempty"`
		SslMode             *int64                        `json:"sslMode,omitempty"`
		Username            *interface{}                  `json:"username,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.AuthenticationType = decoded.AuthenticationType
	s.CommandTimeout = decoded.CommandTimeout
	s.ConnectionString = decoded.ConnectionString
	s.ConnectionTimeout = decoded.ConnectionTimeout
	s.Database = decoded.Database
	s.EncryptedCredential = decoded.EncryptedCredential
	s.Host = decoded.Host
	s.Port = decoded.Port
	s.Pwd = decoded.Pwd
	s.SslMode = decoded.SslMode
	s.Username = decoded.Username

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling GreenplumLinkedServiceTypeProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["password"]; ok {
		impl, err := UnmarshalSecretBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Password' for 'GreenplumLinkedServiceTypeProperties': %+v", err)
		}
		s.Password = impl
	}

	return nil
}
