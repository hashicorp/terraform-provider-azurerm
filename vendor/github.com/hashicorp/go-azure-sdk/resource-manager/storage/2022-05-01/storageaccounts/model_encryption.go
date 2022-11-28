package storageaccounts

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Encryption struct {
	Identity                        *EncryptionIdentity `json:"identity"`
	KeySource                       *KeySource          `json:"keySource,omitempty"`
	Keyvaultproperties              *KeyVaultProperties `json:"keyvaultproperties"`
	RequireInfrastructureEncryption *bool               `json:"requireInfrastructureEncryption,omitempty"`
	Services                        *EncryptionServices `json:"services"`
}
