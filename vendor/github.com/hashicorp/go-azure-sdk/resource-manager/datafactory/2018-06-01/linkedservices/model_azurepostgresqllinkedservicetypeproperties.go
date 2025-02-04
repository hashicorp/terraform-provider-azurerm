package linkedservices

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzurePostgreSqlLinkedServiceTypeProperties struct {
	CommandTimeout         *int64                        `json:"commandTimeout,omitempty"`
	ConnectionString       *string                       `json:"connectionString,omitempty"`
	Database               *string                       `json:"database,omitempty"`
	Encoding               *string                       `json:"encoding,omitempty"`
	EncryptedCredential    *string                       `json:"encryptedCredential,omitempty"`
	Password               *AzureKeyVaultSecretReference `json:"password,omitempty"`
	Port                   *int64                        `json:"port,omitempty"`
	ReadBufferSize         *int64                        `json:"readBufferSize,omitempty"`
	Server                 *string                       `json:"server,omitempty"`
	SslMode                *int64                        `json:"sslMode,omitempty"`
	Timeout                *int64                        `json:"timeout,omitempty"`
	Timezone               *string                       `json:"timezone,omitempty"`
	TrustServerCertificate *bool                         `json:"trustServerCertificate,omitempty"`
	Username               *string                       `json:"username,omitempty"`
}
