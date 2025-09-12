package linkedservices

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MariaDBLinkedServiceTypeProperties struct {
	ConnectionString    *interface{}                  `json:"connectionString,omitempty"`
	Database            *interface{}                  `json:"database,omitempty"`
	DriverVersion       *interface{}                  `json:"driverVersion,omitempty"`
	EncryptedCredential *string                       `json:"encryptedCredential,omitempty"`
	Password            *AzureKeyVaultSecretReference `json:"password,omitempty"`
	Port                *int64                        `json:"port,omitempty"`
	Server              *interface{}                  `json:"server,omitempty"`
	SslMode             *int64                        `json:"sslMode,omitempty"`
	UseSystemTrustStore *int64                        `json:"useSystemTrustStore,omitempty"`
	Username            *interface{}                  `json:"username,omitempty"`
}
