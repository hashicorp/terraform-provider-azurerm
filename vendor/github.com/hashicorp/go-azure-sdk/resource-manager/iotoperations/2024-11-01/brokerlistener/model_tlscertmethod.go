package brokerlistener

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TlsCertMethod struct {
	CertManagerCertificateSpec *CertManagerCertificateSpec `json:"certManagerCertificateSpec,omitempty"`
	Manual                     *X509ManualCertificate      `json:"manual,omitempty"`
	Mode                       TlsCertMethodMode           `json:"mode"`
}
