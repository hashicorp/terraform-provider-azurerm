package virtualnetworkgateways

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VpnClientConfiguration struct {
	AadAudience                       *string                             `json:"aadAudience,omitempty"`
	AadIssuer                         *string                             `json:"aadIssuer,omitempty"`
	AadTenant                         *string                             `json:"aadTenant,omitempty"`
	RadiusServerAddress               *string                             `json:"radiusServerAddress,omitempty"`
	RadiusServerSecret                *string                             `json:"radiusServerSecret,omitempty"`
	RadiusServers                     *[]RadiusServer                     `json:"radiusServers,omitempty"`
	VngClientConnectionConfigurations *[]VngClientConnectionConfiguration `json:"vngClientConnectionConfigurations,omitempty"`
	VpnAuthenticationTypes            *[]VpnAuthenticationType            `json:"vpnAuthenticationTypes,omitempty"`
	VpnClientAddressPool              *AddressSpace                       `json:"vpnClientAddressPool,omitempty"`
	VpnClientIPsecPolicies            *[]IPsecPolicy                      `json:"vpnClientIpsecPolicies,omitempty"`
	VpnClientProtocols                *[]VpnClientProtocol                `json:"vpnClientProtocols,omitempty"`
	VpnClientRevokedCertificates      *[]VpnClientRevokedCertificate      `json:"vpnClientRevokedCertificates,omitempty"`
	VpnClientRootCertificates         *[]VpnClientRootCertificate         `json:"vpnClientRootCertificates,omitempty"`
}
