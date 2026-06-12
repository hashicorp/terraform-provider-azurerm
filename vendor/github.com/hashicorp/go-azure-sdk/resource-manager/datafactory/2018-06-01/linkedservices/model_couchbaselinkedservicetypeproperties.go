package linkedservices

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CouchbaseLinkedServiceTypeProperties struct {
	ConnectionString    *interface{}                  `json:"connectionString,omitempty"`
	CredString          *AzureKeyVaultSecretReference `json:"credString,omitempty"`
	EncryptedCredential *string                       `json:"encryptedCredential,omitempty"`
}
