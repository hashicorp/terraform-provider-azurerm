package webapplicationfirewallpolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApplicationGatewaySslProfilePropertiesFormat struct {
	ClientAuthConfiguration   *ApplicationGatewayClientAuthConfiguration `json:"clientAuthConfiguration,omitempty"`
	ProvisioningState         *ProvisioningState                         `json:"provisioningState,omitempty"`
	SslPolicy                 *ApplicationGatewaySslPolicy               `json:"sslPolicy,omitempty"`
	TrustedClientCertificates *[]SubResource                             `json:"trustedClientCertificates,omitempty"`
}
