package applicationgateways

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApplicationGatewayClientAuthConfiguration struct {
	VerifyClientCertIssuerDN *bool                                      `json:"verifyClientCertIssuerDN,omitempty"`
	VerifyClientRevocation   *ApplicationGatewayClientRevocationOptions `json:"verifyClientRevocation,omitempty"`
}
