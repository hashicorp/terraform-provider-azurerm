package registries

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EncryptionProperty struct {
	KeyVaultProperties *KeyVaultProperties `json:"keyVaultProperties,omitempty"`
	Status             *EncryptionStatus   `json:"status,omitempty"`
}
