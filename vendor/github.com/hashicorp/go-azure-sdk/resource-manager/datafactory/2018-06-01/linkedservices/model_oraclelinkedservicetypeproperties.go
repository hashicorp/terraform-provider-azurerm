package linkedservices

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OracleLinkedServiceTypeProperties struct {
	AuthenticationType        *OracleAuthenticationType     `json:"authenticationType,omitempty"`
	ConnectionString          *interface{}                  `json:"connectionString,omitempty"`
	CryptoChecksumClient      *interface{}                  `json:"cryptoChecksumClient,omitempty"`
	CryptoChecksumTypesClient *interface{}                  `json:"cryptoChecksumTypesClient,omitempty"`
	EnableBulkLoad            *interface{}                  `json:"enableBulkLoad,omitempty"`
	EncryptedCredential       *string                       `json:"encryptedCredential,omitempty"`
	EncryptionClient          *interface{}                  `json:"encryptionClient,omitempty"`
	EncryptionTypesClient     *interface{}                  `json:"encryptionTypesClient,omitempty"`
	FetchSize                 *interface{}                  `json:"fetchSize,omitempty"`
	FetchTswtzAsTimestamp     *interface{}                  `json:"fetchTswtzAsTimestamp,omitempty"`
	InitialLobFetchSize       *interface{}                  `json:"initialLobFetchSize,omitempty"`
	InitializationString      *interface{}                  `json:"initializationString,omitempty"`
	Password                  *AzureKeyVaultSecretReference `json:"password,omitempty"`
	Server                    *interface{}                  `json:"server,omitempty"`
	StatementCacheSize        *interface{}                  `json:"statementCacheSize,omitempty"`
	SupportV1DataTypes        *interface{}                  `json:"supportV1DataTypes,omitempty"`
	Username                  *interface{}                  `json:"username,omitempty"`
}
