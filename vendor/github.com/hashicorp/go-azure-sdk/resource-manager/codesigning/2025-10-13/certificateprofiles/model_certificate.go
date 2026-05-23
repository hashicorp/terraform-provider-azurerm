package certificateprofiles

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Certificate struct {
	CreatedDate      *string            `json:"createdDate,omitempty"`
	EnhancedKeyUsage *string            `json:"enhancedKeyUsage,omitempty"`
	ExpiryDate       *string            `json:"expiryDate,omitempty"`
	Revocation       *Revocation        `json:"revocation,omitempty"`
	SerialNumber     *string            `json:"serialNumber,omitempty"`
	Status           *CertificateStatus `json:"status,omitempty"`
	SubjectName      *string            `json:"subjectName,omitempty"`
	Thumbprint       *string            `json:"thumbprint,omitempty"`
}
