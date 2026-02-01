package globalrulestack

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SecurityServices struct {
	AntiSpywareProfile         *string `json:"antiSpywareProfile,omitempty"`
	AntiVirusProfile           *string `json:"antiVirusProfile,omitempty"`
	DnsSubscription            *string `json:"dnsSubscription,omitempty"`
	FileBlockingProfile        *string `json:"fileBlockingProfile,omitempty"`
	OutboundTrustCertificate   *string `json:"outboundTrustCertificate,omitempty"`
	OutboundUnTrustCertificate *string `json:"outboundUnTrustCertificate,omitempty"`
	UrlFilteringProfile        *string `json:"urlFilteringProfile,omitempty"`
	VulnerabilityProfile       *string `json:"vulnerabilityProfile,omitempty"`
}
