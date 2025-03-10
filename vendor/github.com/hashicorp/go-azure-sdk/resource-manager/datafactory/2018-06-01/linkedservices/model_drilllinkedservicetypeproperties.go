package linkedservices

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DrillLinkedServiceTypeProperties struct {
	ConnectionString    *interface{}                  `json:"connectionString,omitempty"`
	EncryptedCredential *string                       `json:"encryptedCredential,omitempty"`
	Pwd                 *AzureKeyVaultSecretReference `json:"pwd,omitempty"`
}
