package domains

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DomainPropertiesVerificationRecords struct {
	DKIM   *DnsRecord `json:"DKIM,omitempty"`
	DKIM2  *DnsRecord `json:"DKIM2,omitempty"`
	DMARC  *DnsRecord `json:"DMARC,omitempty"`
	Domain *DnsRecord `json:"Domain,omitempty"`
	SPF    *DnsRecord `json:"SPF,omitempty"`
}
