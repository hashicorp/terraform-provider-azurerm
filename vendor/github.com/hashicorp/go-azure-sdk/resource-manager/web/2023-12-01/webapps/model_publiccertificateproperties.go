package webapps

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PublicCertificateProperties struct {
	Blob                      *string                    `json:"blob,omitempty"`
	PublicCertificateLocation *PublicCertificateLocation `json:"publicCertificateLocation,omitempty"`
	Thumbprint                *string                    `json:"thumbprint,omitempty"`
}
