package linkedservices

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MySqlLinkedServiceTypeProperties struct {
	AllowZeroDateTime   *interface{}                  `json:"allowZeroDateTime,omitempty"`
	ConnectionString    *interface{}                  `json:"connectionString,omitempty"`
	ConnectionTimeout   *interface{}                  `json:"connectionTimeout,omitempty"`
	ConvertZeroDateTime *interface{}                  `json:"convertZeroDateTime,omitempty"`
	Database            *interface{}                  `json:"database,omitempty"`
	DriverVersion       *interface{}                  `json:"driverVersion,omitempty"`
	EncryptedCredential *string                       `json:"encryptedCredential,omitempty"`
	GuidFormat          *interface{}                  `json:"guidFormat,omitempty"`
	Password            *AzureKeyVaultSecretReference `json:"password,omitempty"`
	Port                *interface{}                  `json:"port,omitempty"`
	Server              *interface{}                  `json:"server,omitempty"`
	SslCert             *interface{}                  `json:"sslCert,omitempty"`
	SslKey              *interface{}                  `json:"sslKey,omitempty"`
	SslMode             *interface{}                  `json:"sslMode,omitempty"`
	TreatTinyAsBoolean  *interface{}                  `json:"treatTinyAsBoolean,omitempty"`
	UseSystemTrustStore *interface{}                  `json:"useSystemTrustStore,omitempty"`
	Username            *interface{}                  `json:"username,omitempty"`
}
