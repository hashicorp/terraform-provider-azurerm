package clusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type KeyVaultProperties struct {
	KeyName     *string `json:"keyName,omitempty"`
	KeyRsaSize  *int64  `json:"keyRsaSize,omitempty"`
	KeyVaultUri *string `json:"keyVaultUri,omitempty"`
	KeyVersion  *string `json:"keyVersion,omitempty"`
}
