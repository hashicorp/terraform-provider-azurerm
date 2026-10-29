package hcpopenshiftclusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type KmsEncryptionProfile struct {
	ActiveKey    KmsKey             `json:"activeKey"`
	KeyVaultType *KmsKeyVaultType   `json:"keyVaultType,omitempty"`
	VaultName    string             `json:"vaultName"`
	Visibility   KeyVaultVisibility `json:"visibility"`
}
