// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package point_to_site_vpn_gateway

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	virtualwanClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/virtualwan/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourcePointToSiteVPNGatewayRead(d *pluginsdk.ResourceData, meta any) error {
	client := virtualwanClient.From(meta).VirtualWANs
	ctx, cancel := timeouts.ForRead(account.StopContext(meta), d)
	defer cancel()

	id, err := commonids.ParseVirtualWANP2SVPNGatewayID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.P2sVpnGatewaysGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state!", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.GatewayName)
	d.Set("resource_group_name", id.ResourceGroupName)
	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))

		if props := model.Properties; props != nil {
			d.Set("dns_servers", pluginsdk.FlattenSlice(props.CustomDnsServers))
			if err := d.Set("connection_configuration", flattenPointToSiteVPNGatewayConnectionConfiguration(props.P2SConnectionConfigurations)); err != nil {
				return fmt.Errorf("setting `connection_configuration`: %+v", err)
			}

			scaleUnit := 0
			if props.VpnGatewayScaleUnit != nil {
				scaleUnit = int(*props.VpnGatewayScaleUnit)
			}
			d.Set("scale_unit", scaleUnit)

			virtualHubId := ""
			if props.VirtualHub != nil && props.VirtualHub.Id != nil {
				virtualHubId = *props.VirtualHub.Id
			}
			d.Set("virtual_hub_id", virtualHubId)

			vpnServerConfigurationId := ""
			if props.VpnServerConfiguration != nil && props.VpnServerConfiguration.Id != nil {
				vpnServerConfigurationId = *props.VpnServerConfiguration.Id
			}
			d.Set("vpn_server_configuration_id", vpnServerConfigurationId)

			d.Set("routing_preference_internet_enabled", pointer.From(props.IsRoutingPreferenceInternet))
		}
		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}
	return nil
}
