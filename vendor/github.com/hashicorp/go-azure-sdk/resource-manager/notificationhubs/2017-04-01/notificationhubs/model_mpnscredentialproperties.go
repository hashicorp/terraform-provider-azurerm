package notificationhubs

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MpnsCredentialProperties struct {
	CertificateKey  *string `json:"certificateKey,omitempty"`
	MpnsCertificate *string `json:"mpnsCertificate,omitempty"`
	Thumbprint      *string `json:"thumbprint,omitempty"`
}
