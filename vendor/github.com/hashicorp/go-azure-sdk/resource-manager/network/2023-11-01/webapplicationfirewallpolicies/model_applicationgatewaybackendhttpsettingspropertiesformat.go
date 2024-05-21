package webapplicationfirewallpolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApplicationGatewayBackendHTTPSettingsPropertiesFormat struct {
	AffinityCookieName             *string                                `json:"affinityCookieName,omitempty"`
	AuthenticationCertificates     *[]SubResource                         `json:"authenticationCertificates,omitempty"`
	ConnectionDraining             *ApplicationGatewayConnectionDraining  `json:"connectionDraining,omitempty"`
	CookieBasedAffinity            *ApplicationGatewayCookieBasedAffinity `json:"cookieBasedAffinity,omitempty"`
	HostName                       *string                                `json:"hostName,omitempty"`
	Path                           *string                                `json:"path,omitempty"`
	PickHostNameFromBackendAddress *bool                                  `json:"pickHostNameFromBackendAddress,omitempty"`
	Port                           *int64                                 `json:"port,omitempty"`
	Probe                          *SubResource                           `json:"probe,omitempty"`
	ProbeEnabled                   *bool                                  `json:"probeEnabled,omitempty"`
	Protocol                       *ApplicationGatewayProtocol            `json:"protocol,omitempty"`
	ProvisioningState              *ProvisioningState                     `json:"provisioningState,omitempty"`
	RequestTimeout                 *int64                                 `json:"requestTimeout,omitempty"`
	TrustedRootCertificates        *[]SubResource                         `json:"trustedRootCertificates,omitempty"`
}
