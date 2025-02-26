package linkedservices

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OracleLinkedServiceTypeProperties struct {
	AuthenticationType        *OracleAuthenticationType     `json:"authenticationType,omitempty"`
	ConnectionString          string                        `json:"connectionString"`
	CryptoChecksumClient      *string                       `json:"cryptoChecksumClient,omitempty"`
	CryptoChecksumTypesClient *string                       `json:"cryptoChecksumTypesClient,omitempty"`
	EnableBulkLoad            *bool                         `json:"enableBulkLoad,omitempty"`
	EncryptedCredential       *string                       `json:"encryptedCredential,omitempty"`
	EncryptionClient          *string                       `json:"encryptionClient,omitempty"`
	EncryptionTypesClient     *string                       `json:"encryptionTypesClient,omitempty"`
	FetchSize                 *int64                        `json:"fetchSize,omitempty"`
	FetchTswtzAsTimestamp     *bool                         `json:"fetchTswtzAsTimestamp,omitempty"`
	InitialLobFetchSize       *int64                        `json:"initialLobFetchSize,omitempty"`
	InitializationString      *string                       `json:"initializationString,omitempty"`
	Password                  *AzureKeyVaultSecretReference `json:"password,omitempty"`
	Server                    *string                       `json:"server,omitempty"`
	StatementCacheSize        *int64                        `json:"statementCacheSize,omitempty"`
	SupportV1DataTypes        *bool                         `json:"supportV1DataTypes,omitempty"`
	Username                  *string                       `json:"username,omitempty"`
}
