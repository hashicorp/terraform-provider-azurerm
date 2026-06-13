package synonymmaps

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SearchResourceEncryptionKey struct {
	AccessCredentials  *AzureActiveDirectoryApplicationCredentials `json:"accessCredentials,omitempty"`
	KeyVaultKeyName    string                                      `json:"keyVaultKeyName"`
	KeyVaultKeyVersion string                                      `json:"keyVaultKeyVersion"`
	KeyVaultUri        string                                      `json:"keyVaultUri"`
}
