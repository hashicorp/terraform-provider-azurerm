package accounts

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AccountEncryption struct {
	Identity           *ResourceIdentity        `json:"identity,omitempty"`
	KeyVaultProperties *KeyVaultProperties      `json:"keyVaultProperties,omitempty"`
	Status             *string                  `json:"status,omitempty"`
	Type               AccountEncryptionKeyType `json:"type"`
}
