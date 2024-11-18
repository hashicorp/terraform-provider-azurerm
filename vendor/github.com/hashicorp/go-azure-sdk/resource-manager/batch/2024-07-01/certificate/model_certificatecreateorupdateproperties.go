package certificate

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CertificateCreateOrUpdateProperties struct {
	Data                string             `json:"data"`
	Format              *CertificateFormat `json:"format,omitempty"`
	Password            *string            `json:"password,omitempty"`
	Thumbprint          *string            `json:"thumbprint,omitempty"`
	ThumbprintAlgorithm *string            `json:"thumbprintAlgorithm,omitempty"`
}
