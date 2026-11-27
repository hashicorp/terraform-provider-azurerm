// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package vpn_gateway

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/virtualwans"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceVPNGatewayCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).VirtualWan.VirtualWANs
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := virtualwans.NewVpnGatewayID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.VpnGatewaysGet(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_vpn_gateway", id.ID())
		}
	}

	bgpSettingsRaw := d.Get("bgp_settings").([]any)
	bgpSettings := expandVPNGatewayBGPSettings(bgpSettingsRaw)
	payload := virtualwans.VpnGateway{
		Location: pointer.To(location.Normalize(d.Get("location").(string))),
		Properties: &virtualwans.VpnGatewayProperties{
			EnableBgpRouteTranslationForNat: pointer.To(d.Get("bgp_route_translation_for_nat_enabled").(bool)),
			BgpSettings:                     bgpSettings,
			VirtualHub: &virtualwans.SubResource{
				Id: pointer.To(d.Get("virtual_hub_id").(string)),
			},
			VpnGatewayScaleUnit:         pointer.To(int64(d.Get("scale_unit").(int))),
			IsRoutingPreferenceInternet: pointer.To(d.Get("routing_preference").(string) == "Internet"),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]any)),
	}

	if err := client.VpnGatewaysCreateOrUpdateCallbackThenPoll(ctx, id, payload, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}
	d.SetId(id.ID())

	// `vpnGatewayParameters.Properties.bgpSettings.bgpPeeringAddress` customer cannot provide this field during create. This will be set with default value once gateway is created.
	// it could only be updated
	if len(bgpSettingsRaw) > 0 {
		resp, err := client.VpnGatewaysGet(ctx, id)
		if err != nil {
			return fmt.Errorf("retrieving %s: %+v", id, err)
		}
		if resp.Model == nil || resp.Model.Properties == nil {
			return fmt.Errorf("retrieving %s: `model.Properties` was nil", id)
		}
		props := resp.Model.Properties

		if props.BgpSettings != nil && props.BgpSettings.BgpPeeringAddresses != nil {
			val := bgpSettingsRaw[0].(map[string]any)
			input0 := val["instance_0_bgp_peering_address"].([]any)
			input1 := val["instance_1_bgp_peering_address"].([]any)

			if len(input0) > 0 || len(input1) > 0 {
				if len(input0) > 0 && input0[0] != nil {
					val := input0[0].(map[string]any)
					(*props.BgpSettings.BgpPeeringAddresses)[0].CustomBgpIPAddresses = pluginsdk.ExpandStringSlice(val["custom_ips"].(*pluginsdk.Set).List())
				}
				if len(input1) > 0 && input1[0] != nil {
					val := input1[0].(map[string]any)
					(*props.BgpSettings.BgpPeeringAddresses)[1].CustomBgpIPAddresses = pluginsdk.ExpandStringSlice(val["custom_ips"].(*pluginsdk.Set).List())
				}

				resp.Model.Properties = props

				if err := client.VpnGatewaysCreateOrUpdateThenPoll(ctx, id, *resp.Model); err != nil {
					return fmt.Errorf("creating %s: %+v", id, err)
				}
			}
		}
	}

	return resourceVPNGatewayRead(d, meta)
}

func expandVPNGatewayBGPSettings(input []any) *virtualwans.BgpSettings {
	if len(input) == 0 {
		return nil
	}

	val := input[0].(map[string]any)
	return &virtualwans.BgpSettings{
		Asn:        pointer.To(int64(val["asn"].(int))),
		PeerWeight: pointer.To(int64(val["peer_weight"].(int))),
	}
}
