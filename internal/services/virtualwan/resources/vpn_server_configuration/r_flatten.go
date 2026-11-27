// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package vpn_server_configuration

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/virtualwans"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func flattenVpnServerConfigurationAADAuthentication(input *virtualwans.AadAuthenticationParameters) []any {
	if input == nil {
		return []any{}
	}

	return []any{
		map[string]any{
			"audience": pointer.From(input.AadAudience),
			"issuer":   pointer.From(input.AadIssuer),
			"tenant":   pointer.From(input.AadTenant),
		},
	}
}

func flattenVpnServerConfigurationClientRootCertificates(input *[]virtualwans.VpnServerConfigVpnClientRootCertificate) []any {
	if input == nil {
		return []any{}
	}

	output := make([]any, 0)

	for _, v := range *input {
		output = append(output, map[string]any{
			"name":             pointer.From(v.Name),
			"public_cert_data": pointer.From(v.PublicCertData),
		})
	}

	return output
}

func flattenVpnServerConfigurationClientRevokedCertificates(input *[]virtualwans.VpnServerConfigVpnClientRevokedCertificate) []any {
	if input == nil {
		return []any{}
	}

	output := make([]any, 0)
	for _, v := range *input {
		output = append(output, map[string]any{
			"name":       pointer.From(v.Name),
			"thumbprint": pointer.From(v.Thumbprint),
		})
	}
	return output
}

func flattenVpnServerConfigurationIPSecPolicies(input *[]virtualwans.IPsecPolicy) []any {
	if input == nil {
		return []any{}
	}

	output := make([]any, 0)
	for _, v := range *input {
		output = append(output, map[string]any{
			"dh_group":               string(v.DhGroup),
			"ipsec_encryption":       string(v.IPsecEncryption),
			"ipsec_integrity":        string(v.IPsecIntegrity),
			"ike_encryption":         string(v.IkeEncryption),
			"ike_integrity":          string(v.IkeIntegrity),
			"pfs_group":              string(v.PfsGroup),
			"sa_data_size_kilobytes": int(v.SaDataSizeKilobytes),
			"sa_lifetime_seconds":    int(v.SaLifeTimeSeconds),
		})
	}
	return output
}

func flattenVpnServerConfigurationRadius(input *virtualwans.VpnServerConfigurationProperties, d *pluginsdk.ResourceData) []any {
	if input == nil || (pointer.From(input.RadiusServerAddress) == "" && (input.RadiusServers == nil || len(*input.RadiusServers) == 0)) {
		return []any{}
	}

	clientRootCertificates := make([]any, 0)
	if input.RadiusClientRootCertificates != nil {
		for _, v := range *input.RadiusClientRootCertificates {
			clientRootCertificates = append(clientRootCertificates, map[string]any{
				"name":       pointer.From(v.Name),
				"thumbprint": pointer.From(v.Thumbprint),
			})
		}
	}

	serverRootCertificates := make([]any, 0)
	if input.RadiusServerRootCertificates != nil {
		for _, v := range *input.RadiusServerRootCertificates {
			serverRootCertificates = append(serverRootCertificates, map[string]any{
				"name":             pointer.From(v.Name),
				"public_cert_data": pointer.From(v.PublicCertData),
			})
		}
	}

	servers := make([]any, 0)
	if input.RadiusServers != nil && len(*input.RadiusServers) > 0 {
		for i, v := range *input.RadiusServers {
			servers = append(servers, map[string]any{
				"address": v.RadiusServerAddress,
				// setting this because the azure api does not return the secret, so need to read it in from the config
				"secret": d.Get(fmt.Sprintf("radius.0.server.%d.secret", i)),
				"score":  pointer.From(v.RadiusServerScore),
			})
		}
	}

	return []any{
		map[string]any{
			"client_root_certificate": clientRootCertificates,
			"server_root_certificate": serverRootCertificates,
			"server":                  servers,
		},
	}
}

func flattenVpnServerConfigurationVPNProtocols(input *[]virtualwans.VpnGatewayTunnelingProtocol) []any {
	if input == nil {
		return []any{}
	}

	output := make([]any, 0)

	for _, v := range *input {
		output = append(output, string(v))
	}

	return output
}
