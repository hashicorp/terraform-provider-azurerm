// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package vpn_gateway

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/virtualwans"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceVPNGatewayRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).VirtualWan.VirtualWANs
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := virtualwans.ParseVpnGatewayID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.VpnGatewaysGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.VpnGatewayName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))

		if props := model.Properties; props != nil {
			if err := d.Set("bgp_settings", flattenVPNGatewayBGPSettings(props.BgpSettings)); err != nil {
				return fmt.Errorf("setting `bgp_settings`: %+v", err)
			}

			d.Set("bgp_route_translation_for_nat_enabled", pointer.From(props.EnableBgpRouteTranslationForNat))
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

			isRoutingPreferenceInternet := "Microsoft Network"
			if props.IsRoutingPreferenceInternet != nil && *props.IsRoutingPreferenceInternet {
				isRoutingPreferenceInternet = "Internet"
			}
			d.Set("routing_preference", isRoutingPreferenceInternet)

			if err := d.Set("ip_configuration", flattenVPNGatewayIpConfiguration(props.IPConfigurations)); err != nil {
				return fmt.Errorf("setting `ip_configuration`: %+v", err)
			}
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	return nil
}
