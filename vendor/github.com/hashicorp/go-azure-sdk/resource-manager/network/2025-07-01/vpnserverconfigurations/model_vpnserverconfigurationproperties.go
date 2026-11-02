package vpnserverconfigurations

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VpnServerConfigurationProperties struct {
	AadAuthenticationParameters  *AadAuthenticationParameters                  `json:"aadAuthenticationParameters,omitempty"`
	ConfigurationPolicyGroups    *[]VpnServerConfigurationPolicyGroup          `json:"configurationPolicyGroups,omitempty"`
	Etag                         *string                                       `json:"etag,omitempty"`
	Name                         *string                                       `json:"name,omitempty"`
	P2sVpnGateways               *[]P2SVpnGateway                              `json:"p2SVpnGateways,omitempty"`
	ProvisioningState            *string                                       `json:"provisioningState,omitempty"`
	RadiusClientRootCertificates *[]VpnServerConfigRadiusClientRootCertificate `json:"radiusClientRootCertificates,omitempty"`
	RadiusServerAddress          *string                                       `json:"radiusServerAddress,omitempty"`
	RadiusServerRootCertificates *[]VpnServerConfigRadiusServerRootCertificate `json:"radiusServerRootCertificates,omitempty"`
	RadiusServerSecret           *string                                       `json:"radiusServerSecret,omitempty"`
	RadiusServers                *[]RadiusServer                               `json:"radiusServers,omitempty"`
	VpnAuthenticationTypes       *[]VpnAuthenticationType                      `json:"vpnAuthenticationTypes,omitempty"`
	VpnClientIPsecPolicies       *[]IPsecPolicy                                `json:"vpnClientIpsecPolicies,omitempty"`
	VpnClientRevokedCertificates *[]VpnServerConfigVpnClientRevokedCertificate `json:"vpnClientRevokedCertificates,omitempty"`
	VpnClientRootCertificates    *[]VpnServerConfigVpnClientRootCertificate    `json:"vpnClientRootCertificates,omitempty"`
	VpnProtocols                 *[]VpnGatewayTunnelingProtocol                `json:"vpnProtocols,omitempty"`
}
