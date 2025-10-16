package virtualnetworkgateways

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VpnClientParameters struct {
	AuthenticationMethod        *AuthenticationMethod  `json:"authenticationMethod,omitempty"`
	ClientRootCertificates      *[]string              `json:"clientRootCertificates,omitempty"`
	ProcessorArchitecture       *ProcessorArchitecture `json:"processorArchitecture,omitempty"`
	RadiusServerAuthCertificate *string                `json:"radiusServerAuthCertificate,omitempty"`
}
