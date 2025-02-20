package linkedservices

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MySqlLinkedServiceTypeProperties struct {
	AllowZeroDateTime   *bool                         `json:"allowZeroDateTime,omitempty"`
	ConnectionString    *string                       `json:"connectionString,omitempty"`
	ConnectionTimeout   *int64                        `json:"connectionTimeout,omitempty"`
	ConvertZeroDateTime *bool                         `json:"convertZeroDateTime,omitempty"`
	Database            *string                       `json:"database,omitempty"`
	DriverVersion       *string                       `json:"driverVersion,omitempty"`
	EncryptedCredential *string                       `json:"encryptedCredential,omitempty"`
	GuidFormat          *string                       `json:"guidFormat,omitempty"`
	Password            *AzureKeyVaultSecretReference `json:"password,omitempty"`
	Port                *int64                        `json:"port,omitempty"`
	Server              *string                       `json:"server,omitempty"`
	SslCert             *string                       `json:"sslCert,omitempty"`
	SslKey              *string                       `json:"sslKey,omitempty"`
	SslMode             *int64                        `json:"sslMode,omitempty"`
	TreatTinyAsBoolean  *bool                         `json:"treatTinyAsBoolean,omitempty"`
	UseSystemTrustStore *int64                        `json:"useSystemTrustStore,omitempty"`
	Username            *string                       `json:"username,omitempty"`
}
