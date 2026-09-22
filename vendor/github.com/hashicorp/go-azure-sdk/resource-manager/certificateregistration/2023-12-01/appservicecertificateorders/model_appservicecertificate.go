package appservicecertificateorders

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AppServiceCertificate struct {
	KeyVaultId         *string               `json:"keyVaultId,omitempty"`
	KeyVaultSecretName *string               `json:"keyVaultSecretName,omitempty"`
	ProvisioningState  *KeyVaultSecretStatus `json:"provisioningState,omitempty"`
}
