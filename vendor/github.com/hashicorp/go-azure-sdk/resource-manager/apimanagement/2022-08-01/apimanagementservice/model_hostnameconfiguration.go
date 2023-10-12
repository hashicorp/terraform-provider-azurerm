package apimanagementservice

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type HostnameConfiguration struct {
	Certificate                *CertificateInformation `json:"certificate,omitempty"`
	CertificatePassword        *string                 `json:"certificatePassword,omitempty"`
	CertificateSource          *CertificateSource      `json:"certificateSource,omitempty"`
	CertificateStatus          *CertificateStatus      `json:"certificateStatus,omitempty"`
	DefaultSslBinding          *bool                   `json:"defaultSslBinding,omitempty"`
	EncodedCertificate         *string                 `json:"encodedCertificate,omitempty"`
	HostName                   string                  `json:"hostName"`
	IdentityClientId           *string                 `json:"identityClientId,omitempty"`
	KeyVaultId                 *string                 `json:"keyVaultId,omitempty"`
	NegotiateClientCertificate *bool                   `json:"negotiateClientCertificate,omitempty"`
	Type                       HostnameType            `json:"type"`
}
