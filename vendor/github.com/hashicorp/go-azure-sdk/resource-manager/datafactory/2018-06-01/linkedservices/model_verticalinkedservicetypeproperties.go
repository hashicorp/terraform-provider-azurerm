package linkedservices

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VerticaLinkedServiceTypeProperties struct {
	ConnectionString    *interface{}                  `json:"connectionString,omitempty"`
	Database            *interface{}                  `json:"database,omitempty"`
	EncryptedCredential *string                       `json:"encryptedCredential,omitempty"`
	Port                *int64                        `json:"port,omitempty"`
	Pwd                 *AzureKeyVaultSecretReference `json:"pwd,omitempty"`
	Server              *interface{}                  `json:"server,omitempty"`
	Uid                 *interface{}                  `json:"uid,omitempty"`
}
