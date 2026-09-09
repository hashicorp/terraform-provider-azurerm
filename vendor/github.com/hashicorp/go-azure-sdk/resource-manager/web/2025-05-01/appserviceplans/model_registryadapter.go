package appserviceplans

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RegistryAdapter struct {
	KeyVaultSecretReference *KeyVaultReferenceWithStatus `json:"keyVaultSecretReference,omitempty"`
	RegistryKey             *string                      `json:"registryKey,omitempty"`
	Type                    *RegistryAdapterType         `json:"type,omitempty"`
}
