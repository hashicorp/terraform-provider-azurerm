package linkedservices

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OracleLinkedServiceTypeProperties struct {
	AuthenticationType        *OracleAuthenticationType     `json:"authenticationType,omitempty"`
	ConnectionString          *interface{}                  `json:"connectionString,omitempty"`
	CryptoChecksumClient      *interface{}                  `json:"cryptoChecksumClient,omitempty"`
	CryptoChecksumTypesClient *interface{}                  `json:"cryptoChecksumTypesClient,omitempty"`
	EnableBulkLoad            *bool                         `json:"enableBulkLoad,omitempty"`
	EncryptedCredential       *string                       `json:"encryptedCredential,omitempty"`
	EncryptionClient          *interface{}                  `json:"encryptionClient,omitempty"`
	EncryptionTypesClient     *interface{}                  `json:"encryptionTypesClient,omitempty"`
	FetchSize                 *int64                        `json:"fetchSize,omitempty"`
	FetchTswtzAsTimestamp     *bool                         `json:"fetchTswtzAsTimestamp,omitempty"`
	InitialLobFetchSize       *int64                        `json:"initialLobFetchSize,omitempty"`
	InitializationString      *interface{}                  `json:"initializationString,omitempty"`
	Password                  *AzureKeyVaultSecretReference `json:"password,omitempty"`
	Server                    *interface{}                  `json:"server,omitempty"`
	StatementCacheSize        *int64                        `json:"statementCacheSize,omitempty"`
	SupportV1DataTypes        *bool                         `json:"supportV1DataTypes,omitempty"`
	Username                  *interface{}                  `json:"username,omitempty"`
}
