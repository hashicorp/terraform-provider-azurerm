package netappaccounts

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type KeyVaultProperties struct {
	KeyName            string          `json:"keyName"`
	KeyVaultId         *string         `json:"keyVaultId,omitempty"`
	KeyVaultResourceId *string         `json:"keyVaultResourceId,omitempty"`
	KeyVaultUri        string          `json:"keyVaultUri"`
	Status             *KeyVaultStatus `json:"status,omitempty"`
}
