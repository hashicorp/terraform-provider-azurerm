// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package vpn_site

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/virtualwans"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func flattenVpnSiteAddressSpace(input *virtualwans.AddressSpace) []any {
	if input == nil {
		return make([]any, 0)
	}
	return pluginsdk.FlattenSlice(input.AddressPrefixes)
}

func flattenVpnSiteLinks(input *[]virtualwans.VpnSiteLink) []any {
	if input == nil {
		return make([]any, 0)
	}

	output := make([]any, 0)

	for _, e := range *input {
		var (
			ipAddress        string
			fqdn             string
			linkProviderName string
			linkSpeed        int
			bgpProperty      []any
		)

		if prop := e.Properties; prop != nil {
			if prop.IPAddress != nil {
				ipAddress = *prop.IPAddress
			}

			if prop.Fqdn != nil {
				fqdn = *prop.Fqdn
			}

			if linkProp := prop.LinkProperties; linkProp != nil {
				if linkProp.LinkProviderName != nil {
					linkProviderName = *linkProp.LinkProviderName
				}
				if linkProp.LinkSpeedInMbps != nil {
					linkSpeed = int(*linkProp.LinkSpeedInMbps)
				}
			}

			bgpProperty = flattenVpnSiteVpnSiteBgpSettings(prop.BgpProperties)
		}

		link := map[string]any{
			"name":          pointer.From(e.Name),
			"id":            pointer.From(e.Id),
			"provider_name": linkProviderName,
			"speed_in_mbps": linkSpeed,
			"ip_address":    ipAddress,
			"fqdn":          fqdn,
			"bgp":           bgpProperty,
		}

		output = append(output, link)
	}

	return output
}

func flattenVpnSiteVpnSiteBgpSettings(input *virtualwans.VpnLinkBgpSettings) []any {
	if input == nil {
		return make([]any, 0)
	}

	var asn int
	if input.Asn != nil {
		asn = int(*input.Asn)
	}

	return []any{
		map[string]any{
			"asn":             asn,
			"peering_address": pointer.From(input.BgpPeeringAddress),
		},
	}
}

func flattenVpnSiteO365Policy(input *virtualwans.O365PolicyProperties) []any {
	if input == nil {
		return make([]any, 0)
	}

	trafficCategory := make([]any, 0)
	if input.BreakOutCategories != nil {
		trafficCategory = flattenVpnSiteO365TrafficCategoryPolicy(input.BreakOutCategories)
	}

	return []any{
		map[string]any{
			"traffic_category": trafficCategory,
		},
	}
}

func flattenVpnSiteO365TrafficCategoryPolicy(input *virtualwans.O365BreakOutCategoryPolicies) []any {
	if input == nil {
		return make([]any, 0)
	}

	return []any{
		map[string]any{
			"allow_endpoint_enabled":    pointer.From(input.Allow),
			"default_endpoint_enabled":  pointer.From(input.Default),
			"optimize_endpoint_enabled": pointer.From(input.Optimize),
		},
	}
}
