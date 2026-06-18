package keys

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type KeyProperties struct {
	IsActiveCMK *bool   `json:"isActiveCMK,omitempty"`
	KeyVaultURL *string `json:"keyVaultUrl,omitempty"`
}
