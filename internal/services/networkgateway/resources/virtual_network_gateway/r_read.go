// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_network_gateway

import (
	"bytes"
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/virtualnetworkgateways"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceVirtualNetworkGatewayRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).NetworkGateway.VirtualNetworkGateways
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := virtualnetworkgateways.ParseVirtualNetworkGatewayID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.VirtualNetworkGatewayName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))
		d.Set("edge_zone", flattenEdgeZoneModel(model.ExtendedLocation))

		props := model.Properties

		d.Set("bgp_enabled", props.EnableBgp)

		d.Set("type", pointer.FromEnum(props.GatewayType))
		d.Set("private_ip_address_enabled", props.EnablePrivateIPAddress)
		d.Set("active_active", props.ActiveActive)
		d.Set("bgp_route_translation_for_nat_enabled", props.EnableBgpRouteTranslationForNat)
		d.Set("dns_forwarding_enabled", props.EnableDnsForwarding)
		d.Set("ip_sec_replay_protection_enabled", !*props.DisableIPSecReplayProtection)
		d.Set("remote_vnet_traffic_enabled", props.AllowRemoteVnetTraffic)
		d.Set("virtual_wan_traffic_enabled", props.AllowVirtualWanTraffic)
		d.Set("generation", pointer.FromEnum(props.VpnGatewayGeneration))

		if props.VpnType != nil {
			d.Set("vpn_type", pointer.FromEnum(props.VpnType))
		}

		if props.GatewayDefaultSite != nil {
			d.Set("default_local_network_gateway_id", props.GatewayDefaultSite.Id)
		}

		if props.Sku != nil {
			d.Set("sku", pointer.FromEnum(props.Sku.Name))
		}

		gatewayType := pointer.From(props.GatewayType)
		if err := d.Set("ip_configuration", flattenVirtualNetworkGatewayIPConfigurations(props.IPConfigurations, gatewayType)); err != nil {
			return fmt.Errorf("setting `ip_configuration`: %+v", err)
		}

		minScaleUnit, maxScaleUnit := flattenVirtualNetworkGatewayAutoScaleConfiguration(props.AutoScaleConfiguration)
		if err := d.Set("minimum_scale_unit", minScaleUnit); err != nil {
			return fmt.Errorf("setting `minimum_scale_unit`: %+v", err)
		}
		if err := d.Set("maximum_scale_unit", maxScaleUnit); err != nil {
			return fmt.Errorf("setting: `maximum_scale_unit`: %+v", err)
		}

		if err := d.Set("policy_group", flattenVirtualNetworkGatewayPolicyGroups(props.VirtualNetworkGatewayPolicyGroups)); err != nil {
			return fmt.Errorf("setting `policy_group`: %+v", err)
		}

		vpnClientConfig, err := flattenVirtualNetworkGatewayVpnClientConfig(props.VpnClientConfiguration)
		if err != nil {
			return err
		}
		if err := d.Set("vpn_client_configuration", vpnClientConfig); err != nil {
			return fmt.Errorf("setting `vpn_client_configuration`: %+v", err)
		}

		bgpSettings, err := flattenVirtualNetworkGatewayBgpSettings(props.BgpSettings)
		if err != nil {
			return err
		}
		if err := d.Set("bgp_settings", bgpSettings); err != nil {
			return fmt.Errorf("setting `bgp_settings`: %+v", err)
		}

		if err := d.Set("custom_route", flattenVirtualNetworkGatewayAddressSpace(props.CustomRoutes)); err != nil {
			return fmt.Errorf("setting `custom_route`: %+v", err)
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	return nil
}

func resourceVirtualNetworkGatewayCustomizeDiff(ctx context.Context, d *pluginsdk.ResourceDiff, _ any) error {
	gatewayType := d.Get("type").(string)

	// Validate that public_ip_address_id is not set for ExpressRoute gateways
	if gatewayType == string(virtualnetworkgateways.VirtualNetworkGatewayTypeExpressRoute) {
		ipConfigs := d.Get("ip_configuration").([]any)
		for i, ipConfigRaw := range ipConfigs {
			ipConfig := ipConfigRaw.(map[string]any)
			if publicIPID, ok := ipConfig["public_ip_address_id"].(string); ok && publicIPID != "" {
				return fmt.Errorf("`ip_configuration.%d.public_ip_address_id` cannot be set when `type` is set to `ExpressRoute`", i)
			}
		}
	}

	minScaleUnit := d.Get("minimum_scale_unit").(int)
	maxScaleUnit := d.Get("maximum_scale_unit").(int)
	sku := d.Get("sku").(string)

	if sku == string(virtualnetworkgateways.VirtualNetworkGatewaySkuNameErGwScale) {
		if minScaleUnit == 0 || maxScaleUnit == 0 {
			return fmt.Errorf("`minimum_scale_unit` and `maximum_scale_unit` must be set when `sku` is `%s`", virtualnetworkgateways.VirtualNetworkGatewaySkuNameErGwScale)
		}
	}

	// Use RawConfig to determine if the user explicitly set the scale unit fields,
	// since these are Optional+Computed and d.Get() returns API-stored values from state
	rawConfig := d.GetRawConfig().AsValueMap()
	minIsSet := !rawConfig["minimum_scale_unit"].IsNull()
	maxIsSet := !rawConfig["maximum_scale_unit"].IsNull()

	if minIsSet || maxIsSet {
		if sku != string(virtualnetworkgateways.VirtualNetworkGatewaySkuNameErGwScale) {
			return fmt.Errorf("`minimum_scale_unit` and `maximum_scale_unit` are only supported when `sku` is set to `%s`", virtualnetworkgateways.VirtualNetworkGatewaySkuNameErGwScale)
		}

		if minScaleUnit > maxScaleUnit {
			return fmt.Errorf("`minimum_scale_unit` (%d) cannot be greater than `maximum_scale_unit` (%d)", minScaleUnit, maxScaleUnit)
		}
	}

	// The Azure API can't convert an ExpressRoute gateway between the
	// availability-zone SKUs (ErGw1AZ/ErGw2AZ/ErGw3AZ/ErGwScale) and the
	// non-availability-zone SKUs (Standard/HighPerformance/UltraPerformance) in
	// place; the gateway must be deleted and recreated. An in-place change fails
	// with `ExpressRouteVirtualNetworkGatewayAutoscaleBoundsNotValid`.
	if d.Id() != "" && gatewayType == string(virtualnetworkgateways.VirtualNetworkGatewayTypeExpressRoute) && d.HasChange("sku") {
		oldSku, newSku := d.GetChange("sku")
		if expressRouteGatewaySkuIsAvailabilityZone(oldSku.(string)) != expressRouteGatewaySkuIsAvailabilityZone(newSku.(string)) {
			if err := d.ForceNew("sku"); err != nil {
				return err
			}
		}
	}

	return nil
}

func expressRouteGatewaySkuIsAvailabilityZone(sku string) bool {
	switch virtualnetworkgateways.VirtualNetworkGatewaySkuName(sku) {
	case virtualnetworkgateways.VirtualNetworkGatewaySkuNameErGwOneAZ,
		virtualnetworkgateways.VirtualNetworkGatewaySkuNameErGwTwoAZ,
		virtualnetworkgateways.VirtualNetworkGatewaySkuNameErGwThreeAZ,
		virtualnetworkgateways.VirtualNetworkGatewaySkuNameErGwScale:
		return true
	default:
		return false
	}
}

func hashVirtualNetworkGatewayRootCert(v any) int {
	var buf bytes.Buffer
	m := v.(map[string]any)

	fmt.Fprintf(&buf, "%s-", m["name"].(string))
	fmt.Fprintf(&buf, "%s-", m["public_cert_data"].(string))

	return pluginsdk.HashString(buf.String())
}

func hashVirtualNetworkGatewayRevokedCert(v any) int {
	var buf bytes.Buffer
	m := v.(map[string]any)

	fmt.Fprintf(&buf, "%s-", m["name"].(string))
	fmt.Fprintf(&buf, "%s-", m["thumbprint"].(string))

	return pluginsdk.HashString(buf.String())
}

func validateVirtualNetworkGatewayPolicyBasedVpnSku() pluginsdk.SchemaValidateFunc {
	return validation.StringInSlice([]string{
		string(virtualnetworkgateways.VirtualNetworkGatewaySkuTierBasic),
	}, false)
}

func validateVirtualNetworkGatewayRouteBasedVpnSkuGeneration1() pluginsdk.SchemaValidateFunc {
	return validation.StringInSlice([]string{
		string(virtualnetworkgateways.VirtualNetworkGatewaySkuTierBasic),
		string(virtualnetworkgateways.VirtualNetworkGatewaySkuTierStandard),
		string(virtualnetworkgateways.VirtualNetworkGatewaySkuTierHighPerformance),
		string(virtualnetworkgateways.VirtualNetworkGatewaySkuNameVpnGwOne),
		string(virtualnetworkgateways.VirtualNetworkGatewaySkuNameVpnGwTwo),
		string(virtualnetworkgateways.VirtualNetworkGatewaySkuNameVpnGwThree),
		string(virtualnetworkgateways.VirtualNetworkGatewaySkuNameVpnGwOneAZ),
		string(virtualnetworkgateways.VirtualNetworkGatewaySkuNameVpnGwTwoAZ),
		string(virtualnetworkgateways.VirtualNetworkGatewaySkuNameVpnGwThreeAZ),
	}, false)
}

func validateVirtualNetworkGatewayRouteBasedVpnSkuGeneration2() pluginsdk.SchemaValidateFunc {
	return validation.StringInSlice([]string{
		string(virtualnetworkgateways.VirtualNetworkGatewaySkuNameVpnGwTwo),
		string(virtualnetworkgateways.VirtualNetworkGatewaySkuNameVpnGwThree),
		string(virtualnetworkgateways.VirtualNetworkGatewaySkuNameVpnGwFour),
		string(virtualnetworkgateways.VirtualNetworkGatewaySkuNameVpnGwFive),
		string(virtualnetworkgateways.VirtualNetworkGatewaySkuNameVpnGwTwoAZ),
		string(virtualnetworkgateways.VirtualNetworkGatewaySkuNameVpnGwThreeAZ),
		string(virtualnetworkgateways.VirtualNetworkGatewaySkuNameVpnGwFourAZ),
		string(virtualnetworkgateways.VirtualNetworkGatewaySkuNameVpnGwFiveAZ),
	}, false)
}

func validateVirtualNetworkGatewayExpressRouteSku() pluginsdk.SchemaValidateFunc {
	return validation.StringInSlice([]string{
		string(virtualnetworkgateways.VirtualNetworkGatewaySkuTierStandard),
		string(virtualnetworkgateways.VirtualNetworkGatewaySkuTierHighPerformance),
		string(virtualnetworkgateways.VirtualNetworkGatewaySkuTierUltraPerformance),
		string(virtualnetworkgateways.VirtualNetworkGatewaySkuNameErGwOneAZ),
		string(virtualnetworkgateways.VirtualNetworkGatewaySkuNameErGwTwoAZ),
		string(virtualnetworkgateways.VirtualNetworkGatewaySkuNameErGwThreeAZ),
		string(virtualnetworkgateways.VirtualNetworkGatewaySkuNameErGwScale),
	}, false)
}
