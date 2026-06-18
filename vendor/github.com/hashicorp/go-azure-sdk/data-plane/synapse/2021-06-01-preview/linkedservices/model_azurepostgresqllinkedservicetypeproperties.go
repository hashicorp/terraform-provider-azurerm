package linkedservices

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzurePostgreSqlLinkedServiceTypeProperties struct {
	ConnectionString    *interface{}                  `json:"connectionString,omitempty"`
	EncryptedCredential *interface{}                  `json:"encryptedCredential,omitempty"`
	Password            *AzureKeyVaultSecretReference `json:"password,omitempty"`
}
