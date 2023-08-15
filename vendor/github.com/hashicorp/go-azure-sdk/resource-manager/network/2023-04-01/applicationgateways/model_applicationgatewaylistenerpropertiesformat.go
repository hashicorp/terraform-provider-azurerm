package applicationgateways

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApplicationGatewayListenerPropertiesFormat struct {
	FrontendIPConfiguration *SubResource                `json:"frontendIPConfiguration,omitempty"`
	FrontendPort            *SubResource                `json:"frontendPort,omitempty"`
	Protocol                *ApplicationGatewayProtocol `json:"protocol,omitempty"`
	ProvisioningState       *ProvisioningState          `json:"provisioningState,omitempty"`
	SslCertificate          *SubResource                `json:"sslCertificate,omitempty"`
	SslProfile              *SubResource                `json:"sslProfile,omitempty"`
}
