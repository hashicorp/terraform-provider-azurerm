package caches

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CacheEncryptionSettings struct {
	KeyEncryptionKey                  *KeyVaultKeyReference `json:"keyEncryptionKey,omitempty"`
	RotationToLatestKeyVersionEnabled *bool                 `json:"rotationToLatestKeyVersionEnabled,omitempty"`
}
