package apimanagementservice

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CertificateConfiguration struct {
	Certificate         *CertificateInformation `json:"certificate,omitempty"`
	CertificatePassword *string                 `json:"certificatePassword,omitempty"`
	EncodedCertificate  *string                 `json:"encodedCertificate,omitempty"`
	StoreName           StoreName               `json:"storeName"`
}
