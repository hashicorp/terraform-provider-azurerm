package applicationgateways

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApplicationGatewayPropertiesFormat struct {
	AuthenticationCertificates          *[]ApplicationGatewayAuthenticationCertificate         `json:"authenticationCertificates,omitempty"`
	AutoscaleConfiguration              *ApplicationGatewayAutoscaleConfiguration              `json:"autoscaleConfiguration,omitempty"`
	BackendAddressPools                 *[]ApplicationGatewayBackendAddressPool                `json:"backendAddressPools,omitempty"`
	BackendHTTPSettingsCollection       *[]ApplicationGatewayBackendHTTPSettings               `json:"backendHttpSettingsCollection,omitempty"`
	BackendSettingsCollection           *[]ApplicationGatewayBackendSettings                   `json:"backendSettingsCollection,omitempty"`
	CustomErrorConfigurations           *[]ApplicationGatewayCustomError                       `json:"customErrorConfigurations,omitempty"`
	DefaultPredefinedSslPolicy          *ApplicationGatewaySslPolicyName                       `json:"defaultPredefinedSslPolicy,omitempty"`
	EnableFips                          *bool                                                  `json:"enableFips,omitempty"`
	EnableHTTP2                         *bool                                                  `json:"enableHttp2,omitempty"`
	FirewallPolicy                      *SubResource                                           `json:"firewallPolicy,omitempty"`
	ForceFirewallPolicyAssociation      *bool                                                  `json:"forceFirewallPolicyAssociation,omitempty"`
	FrontendIPConfigurations            *[]ApplicationGatewayFrontendIPConfiguration           `json:"frontendIPConfigurations,omitempty"`
	FrontendPorts                       *[]ApplicationGatewayFrontendPort                      `json:"frontendPorts,omitempty"`
	GatewayIPConfigurations             *[]ApplicationGatewayIPConfiguration                   `json:"gatewayIPConfigurations,omitempty"`
	GlobalConfiguration                 *ApplicationGatewayGlobalConfiguration                 `json:"globalConfiguration,omitempty"`
	HTTPListeners                       *[]ApplicationGatewayHTTPListener                      `json:"httpListeners,omitempty"`
	Listeners                           *[]ApplicationGatewayListener                          `json:"listeners,omitempty"`
	LoadDistributionPolicies            *[]ApplicationGatewayLoadDistributionPolicy            `json:"loadDistributionPolicies,omitempty"`
	OperationalState                    *ApplicationGatewayOperationalState                    `json:"operationalState,omitempty"`
	PrivateEndpointConnections          *[]ApplicationGatewayPrivateEndpointConnection         `json:"privateEndpointConnections,omitempty"`
	PrivateLinkConfigurations           *[]ApplicationGatewayPrivateLinkConfiguration          `json:"privateLinkConfigurations,omitempty"`
	Probes                              *[]ApplicationGatewayProbe                             `json:"probes,omitempty"`
	ProvisioningState                   *ProvisioningState                                     `json:"provisioningState,omitempty"`
	RedirectConfigurations              *[]ApplicationGatewayRedirectConfiguration             `json:"redirectConfigurations,omitempty"`
	RequestRoutingRules                 *[]ApplicationGatewayRequestRoutingRule                `json:"requestRoutingRules,omitempty"`
	ResourceGuid                        *string                                                `json:"resourceGuid,omitempty"`
	RewriteRuleSets                     *[]ApplicationGatewayRewriteRuleSet                    `json:"rewriteRuleSets,omitempty"`
	RoutingRules                        *[]ApplicationGatewayRoutingRule                       `json:"routingRules,omitempty"`
	Sku                                 *ApplicationGatewaySku                                 `json:"sku,omitempty"`
	SslCertificates                     *[]ApplicationGatewaySslCertificate                    `json:"sslCertificates,omitempty"`
	SslPolicy                           *ApplicationGatewaySslPolicy                           `json:"sslPolicy,omitempty"`
	SslProfiles                         *[]ApplicationGatewaySslProfile                        `json:"sslProfiles,omitempty"`
	TrustedClientCertificates           *[]ApplicationGatewayTrustedClientCertificate          `json:"trustedClientCertificates,omitempty"`
	TrustedRootCertificates             *[]ApplicationGatewayTrustedRootCertificate            `json:"trustedRootCertificates,omitempty"`
	UrlPathMaps                         *[]ApplicationGatewayUrlPathMap                        `json:"urlPathMaps,omitempty"`
	WebApplicationFirewallConfiguration *ApplicationGatewayWebApplicationFirewallConfiguration `json:"webApplicationFirewallConfiguration,omitempty"`
}
