package certificates

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type X509CertificateProperties struct {
	Ekus           *[]string                `json:"ekus,omitempty"`
	KeyUsage       *[]KeyUsageType          `json:"key_usage,omitempty"`
	Sans           *SubjectAlternativeNames `json:"sans,omitempty"`
	Subject        *string                  `json:"subject,omitempty"`
	ValidityMonths *int64                   `json:"validity_months,omitempty"`
}
