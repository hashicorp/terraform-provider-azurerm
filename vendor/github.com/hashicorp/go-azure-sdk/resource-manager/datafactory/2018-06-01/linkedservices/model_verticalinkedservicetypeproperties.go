package linkedservices

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VerticaLinkedServiceTypeProperties struct {
	ConnectionString    *string                       `json:"connectionString,omitempty"`
	Database            *string                       `json:"database,omitempty"`
	EncryptedCredential *string                       `json:"encryptedCredential,omitempty"`
	Port                *int64                        `json:"port,omitempty"`
	Pwd                 *AzureKeyVaultSecretReference `json:"pwd,omitempty"`
	Server              *string                       `json:"server,omitempty"`
	Uid                 *string                       `json:"uid,omitempty"`
}
