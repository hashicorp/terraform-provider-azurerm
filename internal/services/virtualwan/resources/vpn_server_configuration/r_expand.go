// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package vpn_server_configuration

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/virtualwans"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func expandVpnServerConfigurationAADAuthentication(input []any) *virtualwans.AadAuthenticationParameters {
	if len(input) == 0 {
		return nil
	}

	v := input[0].(map[string]any)
	return &virtualwans.AadAuthenticationParameters{
		AadAudience: pointer.To(v["audience"].(string)),
		AadIssuer:   pointer.To(v["issuer"].(string)),
		AadTenant:   pointer.To(v["tenant"].(string)),
	}
}

func expandVpnServerConfigurationClientRootCertificates(input []any) *[]virtualwans.VpnServerConfigVpnClientRootCertificate {
	clientRootCertificates := make([]virtualwans.VpnServerConfigVpnClientRootCertificate, 0)

	for _, v := range input {
		raw := v.(map[string]any)
		clientRootCertificates = append(clientRootCertificates, virtualwans.VpnServerConfigVpnClientRootCertificate{
			Name:           pointer.To(raw["name"].(string)),
			PublicCertData: pointer.To(raw["public_cert_data"].(string)),
		})
	}

	return &clientRootCertificates
}

func expandVpnServerConfigurationClientRevokedCertificates(input []any) *[]virtualwans.VpnServerConfigVpnClientRevokedCertificate {
	clientRevokedCertificates := make([]virtualwans.VpnServerConfigVpnClientRevokedCertificate, 0)

	for _, v := range input {
		raw := v.(map[string]any)
		clientRevokedCertificates = append(clientRevokedCertificates, virtualwans.VpnServerConfigVpnClientRevokedCertificate{
			Name:       pointer.To(raw["name"].(string)),
			Thumbprint: pointer.To(raw["thumbprint"].(string)),
		})
	}

	return &clientRevokedCertificates
}

func expandVpnServerConfigurationIPSecPolicies(input []any) *[]virtualwans.IPsecPolicy {
	ipSecPolicies := make([]virtualwans.IPsecPolicy, 0)

	for _, raw := range input {
		v := raw.(map[string]any)
		ipSecPolicies = append(ipSecPolicies, virtualwans.IPsecPolicy{
			DhGroup:             virtualwans.DhGroup(v["dh_group"].(string)),
			IkeEncryption:       virtualwans.IkeEncryption(v["ike_encryption"].(string)),
			IkeIntegrity:        virtualwans.IkeIntegrity(v["ike_integrity"].(string)),
			IPsecEncryption:     virtualwans.IPsecEncryption(v["ipsec_encryption"].(string)),
			IPsecIntegrity:      virtualwans.IPsecIntegrity(v["ipsec_integrity"].(string)),
			PfsGroup:            virtualwans.PfsGroup(v["pfs_group"].(string)),
			SaLifeTimeSeconds:   int64(v["sa_lifetime_seconds"].(int)),
			SaDataSizeKilobytes: int64(v["sa_data_size_kilobytes"].(int)),
		})
	}

	return &ipSecPolicies
}

func expandVpnServerConfigurationRadius(input []any) *vpnServerConfigurationRadius {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	val := input[0].(map[string]any)

	clientRootCertificates := make([]virtualwans.VpnServerConfigRadiusClientRootCertificate, 0)
	clientRootCertsRaw := val["client_root_certificate"].(*pluginsdk.Set).List()
	for _, raw := range clientRootCertsRaw {
		v := raw.(map[string]any)
		clientRootCertificates = append(clientRootCertificates, virtualwans.VpnServerConfigRadiusClientRootCertificate{
			Name:       pointer.To(v["name"].(string)),
			Thumbprint: pointer.To(v["thumbprint"].(string)),
		})
	}

	serverRootCertificates := make([]virtualwans.VpnServerConfigRadiusServerRootCertificate, 0)
	serverRootCertsRaw := val["server_root_certificate"].(*pluginsdk.Set).List()
	for _, raw := range serverRootCertsRaw {
		v := raw.(map[string]any)
		serverRootCertificates = append(serverRootCertificates, virtualwans.VpnServerConfigRadiusServerRootCertificate{
			Name:           pointer.To(v["name"].(string)),
			PublicCertData: pointer.To(v["public_cert_data"].(string)),
		})
	}

	radiusServers := make([]virtualwans.RadiusServer, 0)
	address := ""
	secret := ""

	if val["server"] != nil {
		radiusServersRaw := val["server"].([]any)
		for _, raw := range radiusServersRaw {
			v := raw.(map[string]any)
			radiusServers = append(radiusServers, virtualwans.RadiusServer{
				RadiusServerAddress: v["address"].(string),
				RadiusServerSecret:  pointer.To(v["secret"].(string)),
				RadiusServerScore:   pointer.To(int64(v["score"].(int))),
			})
		}
	}

	return &vpnServerConfigurationRadius{
		address:                address,
		secret:                 secret,
		servers:                &radiusServers,
		clientRootCertificates: &clientRootCertificates,
		serverRootCertificates: &serverRootCertificates,
	}
}

func expandVpnServerConfigurationVPNProtocols(input []any) *[]virtualwans.VpnGatewayTunnelingProtocol {
	vpnProtocols := make([]virtualwans.VpnGatewayTunnelingProtocol, 0)

	for _, v := range input {
		vpnProtocols = append(vpnProtocols, virtualwans.VpnGatewayTunnelingProtocol(v.(string)))
	}

	return &vpnProtocols
}
