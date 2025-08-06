package resourceproviders

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CustomDnsSuffixConfigurationProperties struct {
	CertificateURL            *string                           `json:"certificateUrl,omitempty"`
	DnsSuffix                 *string                           `json:"dnsSuffix,omitempty"`
	KeyVaultReferenceIdentity *string                           `json:"keyVaultReferenceIdentity,omitempty"`
	ProvisioningDetails       *string                           `json:"provisioningDetails,omitempty"`
	ProvisioningState         *CustomDnsSuffixProvisioningState `json:"provisioningState,omitempty"`
}
