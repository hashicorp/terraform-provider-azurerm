package linkedservices

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MySqlLinkedServiceTypeProperties struct {
	AllowZeroDateTime   *bool                         `json:"allowZeroDateTime,omitempty"`
	ConnectionString    *interface{}                  `json:"connectionString,omitempty"`
	ConnectionTimeout   *int64                        `json:"connectionTimeout,omitempty"`
	ConvertZeroDateTime *bool                         `json:"convertZeroDateTime,omitempty"`
	Database            *interface{}                  `json:"database,omitempty"`
	DriverVersion       *interface{}                  `json:"driverVersion,omitempty"`
	EncryptedCredential *string                       `json:"encryptedCredential,omitempty"`
	GuidFormat          *interface{}                  `json:"guidFormat,omitempty"`
	Password            *AzureKeyVaultSecretReference `json:"password,omitempty"`
	Port                *int64                        `json:"port,omitempty"`
	Server              *interface{}                  `json:"server,omitempty"`
	SslCert             *interface{}                  `json:"sslCert,omitempty"`
	SslKey              *interface{}                  `json:"sslKey,omitempty"`
	SslMode             *int64                        `json:"sslMode,omitempty"`
	TreatTinyAsBoolean  *bool                         `json:"treatTinyAsBoolean,omitempty"`
	UseSystemTrustStore *int64                        `json:"useSystemTrustStore,omitempty"`
	Username            *interface{}                  `json:"username,omitempty"`
}
