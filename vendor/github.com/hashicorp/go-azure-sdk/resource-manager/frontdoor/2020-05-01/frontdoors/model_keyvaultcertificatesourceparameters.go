package frontdoors

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type KeyVaultCertificateSourceParameters struct {
	SecretName    *string                                   `json:"secretName,omitempty"`
	SecretVersion *string                                   `json:"secretVersion,omitempty"`
	Vault         *KeyVaultCertificateSourceParametersVault `json:"vault,omitempty"`
}
