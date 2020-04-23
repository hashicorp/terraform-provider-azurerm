package privatedns

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
func (r Registration) SupportedDataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"azurerm_private_dns_zone": dataSourceArmPrivateDnsZone(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"azurerm_private_dns_zone":                      resourceArmPrivateDnsZone(),
		"azurerm_private_dns_a_record":                  resourceArmPrivateDnsARecord(),
		"azurerm_private_dns_aaaa_record":               resourceArmPrivateDnsAaaaRecord(),
		"azurerm_private_dns_cname_record":              resourceArmPrivateDnsCNameRecord(),
		"azurerm_private_dns_mx_record":                 resourceArmPrivateDnsMxRecord(),
		"azurerm_private_dns_ptr_record":                resourceArmPrivateDnsPtrRecord(),
		"azurerm_private_dns_srv_record":                resourceArmPrivateDnsSrvRecord(),
		"azurerm_private_dns_txt_record":                resourceArmPrivateDnsTxtRecord(),
		"azurerm_private_dns_zone_virtual_network_link": resourceArmPrivateDnsZoneVirtualNetworkLink(),
	}
}
