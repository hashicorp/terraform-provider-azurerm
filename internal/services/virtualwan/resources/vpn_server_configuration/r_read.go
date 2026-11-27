// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package vpn_server_configuration

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/virtualwans"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceVPNServerConfigurationRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).VirtualWan.VirtualWANs
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := virtualwans.ParseVpnServerConfigurationID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.VpnServerConfigurationsGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.VpnServerConfigurationName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))

		if props := model.Properties; props != nil {
			if err := d.Set("azure_active_directory_authentication", flattenVpnServerConfigurationAADAuthentication(props.AadAuthenticationParameters)); err != nil {
				return fmt.Errorf("setting `azure_active_directory_authentication`: %+v", err)
			}

			if err := d.Set("client_root_certificate", flattenVpnServerConfigurationClientRootCertificates(props.VpnClientRootCertificates)); err != nil {
				return fmt.Errorf("setting `client_root_certificate`: %+v", err)
			}

			if err := d.Set("client_revoked_certificate", flattenVpnServerConfigurationClientRevokedCertificates(props.VpnClientRevokedCertificates)); err != nil {
				return fmt.Errorf("setting `client_revoked_certificate`: %+v", err)
			}

			if err := d.Set("ipsec_policy", flattenVpnServerConfigurationIPSecPolicies(props.VpnClientIPsecPolicies)); err != nil {
				return fmt.Errorf("setting `ipsec_policy`: %+v", err)
			}

			if err := d.Set("radius", flattenVpnServerConfigurationRadius(props, d)); err != nil {
				return fmt.Errorf("setting `radius`: %+v", err)
			}

			vpnAuthenticationTypes := make([]any, 0)
			if props.VpnAuthenticationTypes != nil {
				for _, v := range *props.VpnAuthenticationTypes {
					vpnAuthenticationTypes = append(vpnAuthenticationTypes, string(v))
				}
			}
			if err := d.Set("vpn_authentication_types", vpnAuthenticationTypes); err != nil {
				return fmt.Errorf("setting `vpn_authentication_types`: %+v", err)
			}

			flattenedVpnProtocols := flattenVpnServerConfigurationVPNProtocols(props.VpnProtocols)
			if err := d.Set("vpn_protocols", pluginsdk.NewSet(pluginsdk.HashString, flattenedVpnProtocols)); err != nil {
				return fmt.Errorf("setting `vpn_protocols`: %+v", err)
			}

			if err := tags.FlattenAndSet(d, model.Tags); err != nil {
				return err
			}
		}
	}

	return nil
}

type vpnServerConfigurationRadius struct {
	address                string
	secret                 string
	servers                *[]virtualwans.RadiusServer
	clientRootCertificates *[]virtualwans.VpnServerConfigRadiusClientRootCertificate
	serverRootCertificates *[]virtualwans.VpnServerConfigRadiusServerRootCertificate
}
