package virtualnetworkgatewayconnections

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CertificateAuthentication struct {
	InboundAuthCertificateChain       *[]string `json:"inboundAuthCertificateChain,omitempty"`
	InboundAuthCertificateSubjectName *string   `json:"inboundAuthCertificateSubjectName,omitempty"`
	OutboundAuthCertificate           *string   `json:"outboundAuthCertificate,omitempty"`
}
