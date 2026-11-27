// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package vpn_server_configuration

import (
	"errors"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/virtualwans"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	virtualwanClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/virtualwan/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceVPNServerConfigurationCreate(d *pluginsdk.ResourceData, meta any) error {
	client := virtualwanClient.From(meta).VirtualWANs
	subscriptionId := account.From(meta).SubscriptionId
	ctx, cancel := timeouts.ForCreate(account.StopContext(meta), d)
	defer cancel()

	id := virtualwans.NewVpnServerConfigurationID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if !features.From(meta).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.VpnServerConfigurationsGet(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_vpn_server_configuration", id.ID())
		}
	}

	aadAuthenticationRaw := d.Get("azure_active_directory_authentication").([]any)
	aadAuthentication := expandVpnServerConfigurationAADAuthentication(aadAuthenticationRaw)

	clientRevokedCertsRaw := d.Get("client_revoked_certificate").(*pluginsdk.Set).List()
	clientRevokedCerts := expandVpnServerConfigurationClientRevokedCertificates(clientRevokedCertsRaw)

	clientRootCertsRaw := d.Get("client_root_certificate").(*pluginsdk.Set).List()
	clientRootCerts := expandVpnServerConfigurationClientRootCertificates(clientRootCertsRaw)

	ipSecPoliciesRaw := d.Get("ipsec_policy").([]any)
	ipSecPolicies := expandVpnServerConfigurationIPSecPolicies(ipSecPoliciesRaw)

	radius := expandVpnServerConfigurationRadius(d.Get("radius").([]any))

	vpnProtocolsRaw := d.Get("vpn_protocols").(*pluginsdk.Set).List()
	vpnProtocols := expandVpnServerConfigurationVPNProtocols(vpnProtocolsRaw)

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

	props := virtualwans.VpnServerConfigurationProperties{
		AadAuthenticationParameters:  aadAuthentication,
		VpnAuthenticationTypes:       &vpnAuthenticationTypes,
		VpnClientRootCertificates:    clientRootCerts,
		VpnClientRevokedCertificates: clientRevokedCerts,
		VpnClientIPsecPolicies:       ipSecPolicies,
		VpnProtocols:                 vpnProtocols,
	}

	if supportsAAD && aadAuthentication == nil {
		return errors.New("`azure_active_directory_authentication` must be specified when `vpn_authentication_type` is set to `AAD`")
	}

	// parameter:VpnServerConfigVpnClientRootCertificates is not specified when VpnAuthenticationType as Certificate is selected.
	if supportsCertificates && len(clientRootCertsRaw) == 0 {
		return errors.New("`client_root_certificate` must be specified when `vpn_authentication_type` is set to `Certificate`")
	}

	if supportsRadius {
		if radius == nil {
			return errors.New("`radius` must be specified when `vpn_authentication_type` is set to `Radius`")
		}

		if radius.servers != nil && len(*radius.servers) != 0 {
			props.RadiusServers = radius.servers
		}

		props.RadiusServerAddress = pointer.To(radius.address)
		props.RadiusServerSecret = pointer.To(radius.secret)

		props.RadiusClientRootCertificates = radius.clientRootCertificates
		props.RadiusServerRootCertificates = radius.serverRootCertificates
	}

	location := location.Normalize(d.Get("location").(string))
	t := d.Get("tags").(map[string]any)
	parameters := virtualwans.VpnServerConfiguration{
		Location:   pointer.To(location),
		Properties: &props,
		Tags:       tags.Expand(t),
	}

	if err := client.VpnServerConfigurationsCreateOrUpdateCallbackThenPoll(ctx, id, parameters, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceVPNServerConfigurationRead(d, meta)
}
