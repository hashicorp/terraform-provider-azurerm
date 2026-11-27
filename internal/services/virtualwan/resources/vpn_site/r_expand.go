// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package vpn_site

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/virtualwans"
)

func expandVpnSiteDeviceProperties(vendor, model string) *virtualwans.DeviceProperties {
	if vendor == "" && model == "" {
		return nil
	}
	output := &virtualwans.DeviceProperties{}
	if vendor != "" {
		output.DeviceVendor = &vendor
	}
	if model != "" {
		output.DeviceModel = &model
	}

	return output
}

func expandVpnSiteAddressSpace(input []any) *virtualwans.AddressSpace {
	if len(input) == 0 {
		return nil
	}

	addressPrefixes := make([]string, 0)
	for _, addr := range input {
		addressPrefixes = append(addressPrefixes, addr.(string))
	}

	return &virtualwans.AddressSpace{
		AddressPrefixes: &addressPrefixes,
	}
}

func expandVpnSiteLinks(input []any) *[]virtualwans.VpnSiteLink {
	if len(input) == 0 {
		return nil
	}

	result := make([]virtualwans.VpnSiteLink, 0)
	for _, e := range input {
		if e == nil {
			continue
		}
		e := e.(map[string]any)
		link := virtualwans.VpnSiteLink{
			Name: pointer.To(e["name"].(string)),
			Properties: &virtualwans.VpnSiteLinkProperties{
				LinkProperties: &virtualwans.VpnLinkProviderProperties{
					LinkSpeedInMbps: pointer.To(int64(e["speed_in_mbps"].(int))),
				},
			},
		}

		if v, ok := e["provider_name"]; ok {
			link.Properties.LinkProperties.LinkProviderName = pointer.To(v.(string))
		}
		if v, ok := e["ip_address"]; ok {
			link.Properties.IPAddress = pointer.To(v.(string))
		}
		if v, ok := e["fqdn"]; ok {
			link.Properties.Fqdn = pointer.To(v.(string))
		}
		if v, ok := e["bgp"]; ok {
			link.Properties.BgpProperties = expandVpnSiteVpnLinkBgpSettings(v.([]any))
		}

		result = append(result, link)
	}

	return &result
}

func expandVpnSiteVpnLinkBgpSettings(input []any) *virtualwans.VpnLinkBgpSettings {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	v := input[0].(map[string]any)

	return &virtualwans.VpnLinkBgpSettings{
		Asn:               pointer.To(int64(v["asn"].(int))),
		BgpPeeringAddress: pointer.To(v["peering_address"].(string)),
	}
}

func expandVpnSiteO365Policy(input []any) *virtualwans.O365PolicyProperties {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	o365Policy := input[0].(map[string]any)

	return &virtualwans.O365PolicyProperties{
		BreakOutCategories: expandVpnSiteO365TrafficCategoryPolicy(o365Policy["traffic_category"].([]any)),
	}
}

func expandVpnSiteO365TrafficCategoryPolicy(input []any) *virtualwans.O365BreakOutCategoryPolicies {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	trafficCategory := input[0].(map[string]any)

	return &virtualwans.O365BreakOutCategoryPolicies{
		Allow:    pointer.To(trafficCategory["allow_endpoint_enabled"].(bool)),
		Default:  pointer.To(trafficCategory["default_endpoint_enabled"].(bool)),
		Optimize: pointer.To(trafficCategory["optimize_endpoint_enabled"].(bool)),
	}
}
