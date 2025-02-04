package linkedservices

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PostgreSqlV2LinkedServiceTypeProperties struct {
	AuthenticationType     string                        `json:"authenticationType"`
	CommandTimeout         *int64                        `json:"commandTimeout,omitempty"`
	ConnectionTimeout      *int64                        `json:"connectionTimeout,omitempty"`
	Database               string                        `json:"database"`
	Encoding               *string                       `json:"encoding,omitempty"`
	EncryptedCredential    *string                       `json:"encryptedCredential,omitempty"`
	LogParameters          *bool                         `json:"logParameters,omitempty"`
	Password               *AzureKeyVaultSecretReference `json:"password,omitempty"`
	Pooling                *bool                         `json:"pooling,omitempty"`
	Port                   *int64                        `json:"port,omitempty"`
	ReadBufferSize         *int64                        `json:"readBufferSize,omitempty"`
	Schema                 *string                       `json:"schema,omitempty"`
	Server                 string                        `json:"server"`
	SslCertificate         *string                       `json:"sslCertificate,omitempty"`
	SslKey                 *string                       `json:"sslKey,omitempty"`
	SslMode                int64                         `json:"sslMode"`
	SslPassword            *string                       `json:"sslPassword,omitempty"`
	Timezone               *string                       `json:"timezone,omitempty"`
	TrustServerCertificate *bool                         `json:"trustServerCertificate,omitempty"`
	Username               string                        `json:"username"`
}
