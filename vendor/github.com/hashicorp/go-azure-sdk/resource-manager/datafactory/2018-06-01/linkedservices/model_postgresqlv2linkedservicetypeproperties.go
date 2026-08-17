package linkedservices

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PostgreSqlV2LinkedServiceTypeProperties struct {
	AuthenticationType     interface{}                   `json:"authenticationType"`
	CommandTimeout         *interface{}                  `json:"commandTimeout,omitempty"`
	ConnectionTimeout      *interface{}                  `json:"connectionTimeout,omitempty"`
	Database               interface{}                   `json:"database"`
	Encoding               *interface{}                  `json:"encoding,omitempty"`
	EncryptedCredential    *string                       `json:"encryptedCredential,omitempty"`
	LogParameters          *interface{}                  `json:"logParameters,omitempty"`
	Password               *AzureKeyVaultSecretReference `json:"password,omitempty"`
	Pooling                *interface{}                  `json:"pooling,omitempty"`
	Port                   *interface{}                  `json:"port,omitempty"`
	ReadBufferSize         *interface{}                  `json:"readBufferSize,omitempty"`
	Schema                 *interface{}                  `json:"schema,omitempty"`
	Server                 interface{}                   `json:"server"`
	SslCertificate         *interface{}                  `json:"sslCertificate,omitempty"`
	SslKey                 *interface{}                  `json:"sslKey,omitempty"`
	SslMode                interface{}                   `json:"sslMode"`
	SslPassword            *interface{}                  `json:"sslPassword,omitempty"`
	Timezone               *interface{}                  `json:"timezone,omitempty"`
	TrustServerCertificate *interface{}                  `json:"trustServerCertificate,omitempty"`
	Username               interface{}                   `json:"username"`
}
