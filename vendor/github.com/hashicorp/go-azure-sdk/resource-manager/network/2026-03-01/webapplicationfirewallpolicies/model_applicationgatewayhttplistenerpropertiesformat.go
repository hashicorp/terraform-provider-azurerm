package webapplicationfirewallpolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApplicationGatewayHTTPListenerPropertiesFormat struct {
	CustomErrorConfigurations   *[]ApplicationGatewayCustomError `json:"customErrorConfigurations,omitempty"`
	FirewallPolicy              *SubResource                     `json:"firewallPolicy,omitempty"`
	FrontendIPConfiguration     *SubResource                     `json:"frontendIPConfiguration,omitempty"`
	FrontendPort                *SubResource                     `json:"frontendPort,omitempty"`
	HostName                    *string                          `json:"hostName,omitempty"`
	HostNames                   *[]string                        `json:"hostNames,omitempty"`
	Protocol                    *ApplicationGatewayProtocol      `json:"protocol,omitempty"`
	ProvisioningState           *ProvisioningState               `json:"provisioningState,omitempty"`
	RequireServerNameIndication *bool                            `json:"requireServerNameIndication,omitempty"`
	SslCertificate              *SubResource                     `json:"sslCertificate,omitempty"`
	SslProfile                  *SubResource                     `json:"sslProfile,omitempty"`
}
