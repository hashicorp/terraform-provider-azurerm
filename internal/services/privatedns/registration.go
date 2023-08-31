// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package privatedns

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

func (r Registration) AssociatedGitHubLabel() string {
	return "service/dns"
}

// Name is the name of this Service
func (r Registration) Name() string {
	return "Private DNS"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"Private DNS",
	}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_private_dns_zone":                      dataSourcePrivateDnsZone(),
		"azurerm_private_dns_a_record":                  dataSourcePrivateDnsARecord(),
		"azurerm_private_dns_aaaa_record":               dataSourcePrivateDnsAaaaRecord(),
		"azurerm_private_dns_cname_record":              dataSourcePrivateDnsCNameRecord(),
		"azurerm_private_dns_mx_record":                 dataSourcePrivateDnsMxRecord(),
		"azurerm_private_dns_ptr_record":                dataSourcePrivateDnsPtrRecord(),
		"azurerm_private_dns_soa_record":                dataSourcePrivateDnsSoaRecord(),
		"azurerm_private_dns_srv_record":                dataSourcePrivateDnsSrvRecord(),
		"azurerm_private_dns_txt_record":                dataSourcePrivateDnsTxtRecord(),
		"azurerm_private_dns_zone_virtual_network_link": dataSourcePrivateDnsZoneVirtualNetworkLink(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_private_dns_zone":                      resourcePrivateDnsZone(),
		"azurerm_private_dns_a_record":                  resourcePrivateDnsARecord(),
		"azurerm_private_dns_aaaa_record":               resourcePrivateDnsAaaaRecord(),
		"azurerm_private_dns_cname_record":              resourcePrivateDnsCNameRecord(),
		"azurerm_private_dns_mx_record":                 resourcePrivateDnsMxRecord(),
		"azurerm_private_dns_ptr_record":                resourcePrivateDnsPtrRecord(),
		"azurerm_private_dns_srv_record":                resourcePrivateDnsSrvRecord(),
		"azurerm_private_dns_txt_record":                resourcePrivateDnsTxtRecord(),
		"azurerm_private_dns_zone_virtual_network_link": resourcePrivateDnsZoneVirtualNetworkLink(),
	}
}
