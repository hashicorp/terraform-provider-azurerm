package linkedservices

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PostgreSqlV2LinkedServiceTypeProperties struct {
	AuthenticationType     interface{}                   `json:"authenticationType"`
	CommandTimeout         *int64                        `json:"commandTimeout,omitempty"`
	ConnectionTimeout      *int64                        `json:"connectionTimeout,omitempty"`
	Database               interface{}                   `json:"database"`
	Encoding               *interface{}                  `json:"encoding,omitempty"`
	EncryptedCredential    *string                       `json:"encryptedCredential,omitempty"`
	LogParameters          *bool                         `json:"logParameters,omitempty"`
	Password               *AzureKeyVaultSecretReference `json:"password,omitempty"`
	Pooling                *bool                         `json:"pooling,omitempty"`
	Port                   *int64                        `json:"port,omitempty"`
	ReadBufferSize         *int64                        `json:"readBufferSize,omitempty"`
	Schema                 *interface{}                  `json:"schema,omitempty"`
	Server                 interface{}                   `json:"server"`
	SslCertificate         *interface{}                  `json:"sslCertificate,omitempty"`
	SslKey                 *interface{}                  `json:"sslKey,omitempty"`
	SslMode                int64                         `json:"sslMode"`
	SslPassword            *interface{}                  `json:"sslPassword,omitempty"`
	Timezone               *interface{}                  `json:"timezone,omitempty"`
	TrustServerCertificate *bool                         `json:"trustServerCertificate,omitempty"`
	Username               interface{}                   `json:"username"`
}
