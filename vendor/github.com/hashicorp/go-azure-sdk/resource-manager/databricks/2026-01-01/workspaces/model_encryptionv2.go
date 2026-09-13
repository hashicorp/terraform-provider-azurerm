package workspaces

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EncryptionV2 struct {
	KeySource          EncryptionKeySource             `json:"keySource"`
	KeyVaultProperties *EncryptionV2KeyVaultProperties `json:"keyVaultProperties,omitempty"`
}
