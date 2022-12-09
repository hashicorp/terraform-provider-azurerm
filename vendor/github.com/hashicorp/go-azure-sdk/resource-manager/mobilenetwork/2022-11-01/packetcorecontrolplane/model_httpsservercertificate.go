package packetcorecontrolplane

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type HTTPSServerCertificate struct {
	CertificateUrl string                   `json:"certificateUrl"`
	Provisioning   *CertificateProvisioning `json:"provisioning,omitempty"`
}
