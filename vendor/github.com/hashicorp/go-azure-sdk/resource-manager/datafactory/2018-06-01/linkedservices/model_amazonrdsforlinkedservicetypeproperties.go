package linkedservices

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AmazonRdsForLinkedServiceTypeProperties struct {
	AuthenticationType        *AmazonRdsForOracleAuthenticationType `json:"authenticationType,omitempty"`
	ConnectionString          *interface{}                          `json:"connectionString,omitempty"`
	CryptoChecksumClient      *interface{}                          `json:"cryptoChecksumClient,omitempty"`
	CryptoChecksumTypesClient *interface{}                          `json:"cryptoChecksumTypesClient,omitempty"`
	EnableBulkLoad            *bool                                 `json:"enableBulkLoad,omitempty"`
	EncryptedCredential       *string                               `json:"encryptedCredential,omitempty"`
	EncryptionClient          *interface{}                          `json:"encryptionClient,omitempty"`
	EncryptionTypesClient     *interface{}                          `json:"encryptionTypesClient,omitempty"`
	FetchSize                 *int64                                `json:"fetchSize,omitempty"`
	FetchTswtzAsTimestamp     *bool                                 `json:"fetchTswtzAsTimestamp,omitempty"`
	InitialLobFetchSize       *int64                                `json:"initialLobFetchSize,omitempty"`
	InitializationString      *interface{}                          `json:"initializationString,omitempty"`
	Password                  SecretBase                            `json:"password"`
	Server                    *interface{}                          `json:"server,omitempty"`
	StatementCacheSize        *int64                                `json:"statementCacheSize,omitempty"`
	SupportV1DataTypes        *bool                                 `json:"supportV1DataTypes,omitempty"`
	Username                  *interface{}                          `json:"username,omitempty"`
}

var _ json.Unmarshaler = &AmazonRdsForLinkedServiceTypeProperties{}

func (s *AmazonRdsForLinkedServiceTypeProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		AuthenticationType        *AmazonRdsForOracleAuthenticationType `json:"authenticationType,omitempty"`
		ConnectionString          *interface{}                          `json:"connectionString,omitempty"`
		CryptoChecksumClient      *interface{}                          `json:"cryptoChecksumClient,omitempty"`
		CryptoChecksumTypesClient *interface{}                          `json:"cryptoChecksumTypesClient,omitempty"`
		EnableBulkLoad            *bool                                 `json:"enableBulkLoad,omitempty"`
		EncryptedCredential       *string                               `json:"encryptedCredential,omitempty"`
		EncryptionClient          *interface{}                          `json:"encryptionClient,omitempty"`
		EncryptionTypesClient     *interface{}                          `json:"encryptionTypesClient,omitempty"`
		FetchSize                 *int64                                `json:"fetchSize,omitempty"`
		FetchTswtzAsTimestamp     *bool                                 `json:"fetchTswtzAsTimestamp,omitempty"`
		InitialLobFetchSize       *int64                                `json:"initialLobFetchSize,omitempty"`
		InitializationString      *interface{}                          `json:"initializationString,omitempty"`
		Server                    *interface{}                          `json:"server,omitempty"`
		StatementCacheSize        *int64                                `json:"statementCacheSize,omitempty"`
		SupportV1DataTypes        *bool                                 `json:"supportV1DataTypes,omitempty"`
		Username                  *interface{}                          `json:"username,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.AuthenticationType = decoded.AuthenticationType
	s.ConnectionString = decoded.ConnectionString
	s.CryptoChecksumClient = decoded.CryptoChecksumClient
	s.CryptoChecksumTypesClient = decoded.CryptoChecksumTypesClient
	s.EnableBulkLoad = decoded.EnableBulkLoad
	s.EncryptedCredential = decoded.EncryptedCredential
	s.EncryptionClient = decoded.EncryptionClient
	s.EncryptionTypesClient = decoded.EncryptionTypesClient
	s.FetchSize = decoded.FetchSize
	s.FetchTswtzAsTimestamp = decoded.FetchTswtzAsTimestamp
	s.InitialLobFetchSize = decoded.InitialLobFetchSize
	s.InitializationString = decoded.InitializationString
	s.Server = decoded.Server
	s.StatementCacheSize = decoded.StatementCacheSize
	s.SupportV1DataTypes = decoded.SupportV1DataTypes
	s.Username = decoded.Username

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling AmazonRdsForLinkedServiceTypeProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["password"]; ok {
		impl, err := UnmarshalSecretBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Password' for 'AmazonRdsForLinkedServiceTypeProperties': %+v", err)
		}
		s.Password = impl
	}

	return nil
}
