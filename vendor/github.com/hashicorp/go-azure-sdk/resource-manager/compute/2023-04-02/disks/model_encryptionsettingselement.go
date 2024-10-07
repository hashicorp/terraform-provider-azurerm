package disks

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EncryptionSettingsElement struct {
	DiskEncryptionKey *KeyVaultAndSecretReference `json:"diskEncryptionKey,omitempty"`
	KeyEncryptionKey  *KeyVaultAndKeyReference    `json:"keyEncryptionKey,omitempty"`
}
