// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package vpn_server_configuration

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/virtualwans"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceVPNServerConfigurationUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).VirtualWan.VirtualWANs
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := virtualwans.ParseVpnServerConfigurationID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.VpnServerConfigurationsGet(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", id)
	}

	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", id)
	}

	payload := existing.Model

	if d.HasChange("azure_active_directory_authentication") {
		payload.Properties.AadAuthenticationParameters = expandVpnServerConfigurationAADAuthentication(d.Get("azure_active_directory_authentication").([]any))
	}

	if d.HasChange("client_revoked_certificate") {
		payload.Properties.VpnClientRevokedCertificates = expandVpnServerConfigurationClientRevokedCertificates(d.Get("client_revoked_certificate").(*pluginsdk.Set).List())
	}

	if d.HasChange("client_root_certificate") {
		payload.Properties.VpnClientRootCertificates = expandVpnServerConfigurationClientRootCertificates(d.Get("client_root_certificate").(*pluginsdk.Set).List())
	}

	if d.HasChange("ipsec_policy") {
		payload.Properties.VpnClientIPsecPolicies = expandVpnServerConfigurationIPSecPolicies(d.Get("ipsec_policy").([]any))
	}

	if d.HasChange("vpn_protocols") {
		payload.Properties.VpnProtocols = expandVpnServerConfigurationVPNProtocols(d.Get("vpn_protocols").(*pluginsdk.Set).List())
	}

	supportsAAD := false
	supportsCertificates := false
	supportsRadius := false

	vpnAuthenticationTypesRaw := d.Get("vpn_authentication_types").([]any)
	vpnAuthenticationTypes := make([]virtualwans.VpnAuthenticationType, 0)
	for _, v := range vpnAuthenticationTypesRaw {
		authType := virtualwans.VpnAuthenticationType(v.(string))

		switch authType {
		case virtualwans.VpnAuthenticationTypeAAD:
			supportsAAD = true

		case virtualwans.VpnAuthenticationTypeCertificate:
			supportsCertificates = true

		case virtualwans.VpnAuthenticationTypeRadius:
			supportsRadius = true

		default:
			return fmt.Errorf("unsupported `vpn_authentication_type`: %q", authType)
		}

		vpnAuthenticationTypes = append(vpnAuthenticationTypes, authType)
	}

	if d.HasChange("vpn_authentication_types") {
		payload.Properties.VpnAuthenticationTypes = &vpnAuthenticationTypes
	}

	if d.HasChange("radius") {
		// if radius has changed, we'll nil out the radius attributes and update them to new values if needed
		payload.Properties.RadiusServerAddress = nil
		payload.Properties.RadiusServerSecret = nil
		payload.Properties.RadiusClientRootCertificates = nil
		payload.Properties.RadiusServerRootCertificates = nil
		payload.Properties.RadiusServers = nil

		radius := expandVpnServerConfigurationRadius(d.Get("radius").([]any))
		if supportsRadius {
			if radius == nil {
				return fmt.Errorf("`radius` must be specified when `vpn_authentication_type` is set to `Radius`")
			}

			if radius.servers != nil && len(*radius.servers) != 0 {
				payload.Properties.RadiusServers = radius.servers
			}

			payload.Properties.RadiusServerAddress = pointer.To(radius.address)
			payload.Properties.RadiusServerSecret = pointer.To(radius.secret)

			payload.Properties.RadiusClientRootCertificates = radius.clientRootCertificates
			payload.Properties.RadiusServerRootCertificates = radius.serverRootCertificates
		}
	}

	if supportsAAD && payload.Properties.AadAuthenticationParameters == nil {
		return fmt.Errorf("`azure_active_directory_authentication` must be specified when `vpn_authentication_type` is set to `AAD`")
	}

	// parameter:VpnServerConfigVpnClientRootCertificates is not specified when VpnAuthenticationType as Certificate is selected.
	if supportsCertificates && payload.Properties.VpnClientRootCertificates != nil && len(*payload.Properties.VpnClientRootCertificates) == 0 {
		return fmt.Errorf("`client_root_certificate` must be specified when `vpn_authentication_type` is set to `Certificate`")
	}

	if d.HasChange("tags") {
		payload.Tags = tags.Expand(d.Get("tags").(map[string]any))
	}

	if err := client.VpnServerConfigurationsCreateOrUpdateThenPoll(ctx, *id, *payload); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceVPNServerConfigurationRead(d, meta)
}
