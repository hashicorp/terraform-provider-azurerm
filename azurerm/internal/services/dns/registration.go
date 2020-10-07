package dns

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type Registration struct{}

// Name is the name of this Service
func (r Registration) Name() string {
	return "DNS"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"DNS",
	}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"azurerm_dns_zone": dataSourceArmDnsZone(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"azurerm_dns_a_record":     resourceArmDnsARecord(),
		"azurerm_dns_aaaa_record":  resourceArmDnsAAAARecord(),
		"azurerm_dns_caa_record":   resourceArmDnsCaaRecord(),
		"azurerm_dns_cname_record": resourceArmDnsCNameRecord(),
		"azurerm_dns_mx_record":    resourceArmDnsMxRecord(),
		"azurerm_dns_ns_record":    resourceArmDnsNsRecord(),
		"azurerm_dns_ptr_record":   resourceArmDnsPtrRecord(),
		"azurerm_dns_srv_record":   resourceArmDnsSrvRecord(),
		"azurerm_dns_txt_record":   resourceArmDnsTxtRecord(),
		"azurerm_dns_zone":         resourceArmDnsZone()}
}
