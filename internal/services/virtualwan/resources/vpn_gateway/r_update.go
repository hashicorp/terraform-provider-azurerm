// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package vpn_gateway

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/virtualwans"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceVPNGatewayUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).VirtualWan.VirtualWANs
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := virtualwans.ParseVpnGatewayID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.VpnGatewayName, VPNGatewayResourceName)
	defer locks.UnlockByName(id.VpnGatewayName, VPNGatewayResourceName)

	existing, err := client.VpnGatewaysGet(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	if existing.Model == nil || existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `model.Properties` was nil", *id)
	}

	model := *existing.Model
	if d.HasChange("scale_unit") {
		model.Properties.VpnGatewayScaleUnit = pointer.To(int64(d.Get("scale_unit").(int)))
	}
	if d.HasChange("tags") {
		model.Tags = tags.Expand(d.Get("tags").(map[string]any))
	}
	if d.HasChange("bgp_route_translation_for_nat_enabled") {
		model.Properties.EnableBgpRouteTranslationForNat = pointer.To(d.Get("bgp_route_translation_for_nat_enabled").(bool))
	}

	bgpSettingsRaw := d.Get("bgp_settings").([]any)
	if len(bgpSettingsRaw) > 0 {
		val := bgpSettingsRaw[0].(map[string]any)

		if d.HasChange("bgp_settings.0.instance_0_bgp_peering_address") {
			if input := val["instance_0_bgp_peering_address"].([]any); len(input) > 0 {
				val := input[0].(map[string]any)
				(*model.Properties.BgpSettings.BgpPeeringAddresses)[0].CustomBgpIPAddresses = pluginsdk.ExpandStringSlice(val["custom_ips"].(*pluginsdk.Set).List())
			}
		}
		if d.HasChange("bgp_settings.0.instance_1_bgp_peering_address") {
			if input := val["instance_1_bgp_peering_address"].([]any); len(input) > 0 {
				val := input[0].(map[string]any)
				(*model.Properties.BgpSettings.BgpPeeringAddresses)[1].CustomBgpIPAddresses = pluginsdk.ExpandStringSlice(val["custom_ips"].(*pluginsdk.Set).List())
			}
		}
	}

	if err := client.VpnGatewaysCreateOrUpdateThenPoll(ctx, *id, model); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceVPNGatewayRead(d, meta)
}
