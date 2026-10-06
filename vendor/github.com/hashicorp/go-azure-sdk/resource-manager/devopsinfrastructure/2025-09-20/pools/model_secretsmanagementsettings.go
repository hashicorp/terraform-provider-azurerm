package pools

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SecretsManagementSettings struct {
	CertificateStoreLocation *string                     `json:"certificateStoreLocation,omitempty"`
	CertificateStoreName     *CertificateStoreNameOption `json:"certificateStoreName,omitempty"`
	KeyExportable            bool                        `json:"keyExportable"`
	ObservedCertificates     []string                    `json:"observedCertificates"`
}
