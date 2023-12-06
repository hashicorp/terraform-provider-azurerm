package nginxcertificate

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NginxCertificateProperties struct {
	CertificateVirtualPath *string            `json:"certificateVirtualPath,omitempty"`
	KeyVaultSecretId       *string            `json:"keyVaultSecretId,omitempty"`
	KeyVirtualPath         *string            `json:"keyVirtualPath,omitempty"`
	ProvisioningState      *ProvisioningState `json:"provisioningState,omitempty"`
}
