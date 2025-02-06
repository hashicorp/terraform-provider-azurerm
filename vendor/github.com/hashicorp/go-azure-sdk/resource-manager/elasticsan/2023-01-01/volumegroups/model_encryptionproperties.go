package volumegroups

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EncryptionProperties struct {
	Identity           *EncryptionIdentity `json:"identity,omitempty"`
	KeyVaultProperties *KeyVaultProperties `json:"keyVaultProperties,omitempty"`
}
