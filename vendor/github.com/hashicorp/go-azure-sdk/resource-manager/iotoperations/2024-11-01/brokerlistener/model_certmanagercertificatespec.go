package brokerlistener

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CertManagerCertificateSpec struct {
	Duration    *string                `json:"duration,omitempty"`
	IssuerRef   CertManagerIssuerRef   `json:"issuerRef"`
	PrivateKey  *CertManagerPrivateKey `json:"privateKey,omitempty"`
	RenewBefore *string                `json:"renewBefore,omitempty"`
	San         *SanForCert            `json:"san,omitempty"`
	SecretName  *string                `json:"secretName,omitempty"`
}
