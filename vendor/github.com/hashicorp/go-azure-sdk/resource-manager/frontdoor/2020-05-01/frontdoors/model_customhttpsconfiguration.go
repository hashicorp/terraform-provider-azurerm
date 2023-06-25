package frontdoors

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CustomHTTPSConfiguration struct {
	CertificateSource                    FrontDoorCertificateSource            `json:"certificateSource"`
	FrontDoorCertificateSourceParameters *FrontDoorCertificateSourceParameters `json:"frontDoorCertificateSourceParameters,omitempty"`
	KeyVaultCertificateSourceParameters  *KeyVaultCertificateSourceParameters  `json:"keyVaultCertificateSourceParameters,omitempty"`
	MinimumTlsVersion                    MinimumTLSVersion                     `json:"minimumTlsVersion"`
	ProtocolType                         FrontDoorTlsProtocolType              `json:"protocolType"`
}
