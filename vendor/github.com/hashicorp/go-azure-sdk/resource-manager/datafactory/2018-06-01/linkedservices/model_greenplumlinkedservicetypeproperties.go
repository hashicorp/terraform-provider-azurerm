package linkedservices

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
	Port                *int64                        `json:"port,omitempty"`
	Pwd                 *AzureKeyVaultSecretReference `json:"pwd,omitempty"`
	SslMode             *int64                        `json:"sslMode,omitempty"`
	Username            *interface{}                  `json:"username,omitempty"`
}
