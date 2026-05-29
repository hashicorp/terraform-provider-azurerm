package devices

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Secret struct {
	EncryptedSecret *AsymmetricEncryptedSecret `json:"encryptedSecret,omitempty"`
	KeyVaultId      *string                    `json:"keyVaultId,omitempty"`
}
