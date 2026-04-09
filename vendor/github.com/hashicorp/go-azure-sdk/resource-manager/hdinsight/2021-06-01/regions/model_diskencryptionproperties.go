package regions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DiskEncryptionProperties struct {
	EncryptionAlgorithm *JsonWebKeyEncryptionAlgorithm `json:"encryptionAlgorithm,omitempty"`
	EncryptionAtHost    *bool                          `json:"encryptionAtHost,omitempty"`
	KeyName             *string                        `json:"keyName,omitempty"`
	KeyVersion          *string                        `json:"keyVersion,omitempty"`
	MsiResourceId       *string                        `json:"msiResourceId,omitempty"`
	VaultUri            *string                        `json:"vaultUri,omitempty"`
}
