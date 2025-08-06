package replicationprotecteditems

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DiskEncryptionKeyInfo struct {
	KeyVaultResourceArmId *string `json:"keyVaultResourceArmId,omitempty"`
	SecretIdentifier      *string `json:"secretIdentifier,omitempty"`
}
