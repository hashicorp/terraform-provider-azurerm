package gatewayhostnameconfiguration

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GatewayHostnameConfigurationContractProperties struct {
	CertificateId              *string `json:"certificateId,omitempty"`
	HTTP2Enabled               *bool   `json:"http2Enabled,omitempty"`
	Hostname                   *string `json:"hostname,omitempty"`
	NegotiateClientCertificate *bool   `json:"negotiateClientCertificate,omitempty"`
	Tls10Enabled               *bool   `json:"tls10Enabled,omitempty"`
	Tls11Enabled               *bool   `json:"tls11Enabled,omitempty"`
}
