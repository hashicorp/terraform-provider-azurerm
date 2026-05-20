package serviceresource

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GetTdeCertificatesSqlTaskOutput struct {
	Base64EncodedCertificates *map[string][]string   `json:"base64EncodedCertificates,omitempty"`
	ValidationErrors          *[]ReportableException `json:"validationErrors,omitempty"`
}
