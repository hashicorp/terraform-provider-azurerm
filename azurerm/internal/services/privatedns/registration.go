package privatedns

import (
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

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
		"azurerm_private_dns_zone": dataSourcePrivateDnsZone(),
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
