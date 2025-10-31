package automationaccount

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EncryptionProperties struct {
	Identity           *EncryptionPropertiesIdentity `json:"identity,omitempty"`
	KeySource          *EncryptionKeySourceType      `json:"keySource,omitempty"`
	KeyVaultProperties *KeyVaultProperties           `json:"keyVaultProperties,omitempty"`
}
