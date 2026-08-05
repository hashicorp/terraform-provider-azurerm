package certificates

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CertificateIssuerSetParameters struct {
	Attributes  *IssuerAttributes    `json:"attributes,omitempty"`
	Credentials *IssuerCredentials   `json:"credentials,omitempty"`
	OrgDetails  *OrganizationDetails `json:"org_details,omitempty"`
	Provider    string               `json:"provider"`
}
