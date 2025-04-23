package linkedservices

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureStorageLinkedServiceTypeProperties struct {
	AccountKey          *AzureKeyVaultSecretReference `json:"accountKey,omitempty"`
	ConnectionString    *interface{}                  `json:"connectionString,omitempty"`
	EncryptedCredential *string                       `json:"encryptedCredential,omitempty"`
	SasToken            *AzureKeyVaultSecretReference `json:"sasToken,omitempty"`
	SasUri              *interface{}                  `json:"sasUri,omitempty"`
}
