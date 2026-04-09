package dpscertificate

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CertificateProperties struct {
	Certificate *string `json:"certificate,omitempty"`
	Created     *string `json:"created,omitempty"`
	Expiry      *string `json:"expiry,omitempty"`
	IsVerified  *bool   `json:"isVerified,omitempty"`
	Subject     *string `json:"subject,omitempty"`
	Thumbprint  *string `json:"thumbprint,omitempty"`
	Updated     *string `json:"updated,omitempty"`
}
