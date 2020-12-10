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
		"azurerm_dns_zone": dataSourceDnsZone(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"azurerm_dns_a_record":     resourceDnsARecord(),
		"azurerm_dns_aaaa_record":  resourceDnsAAAARecord(),
		"azurerm_dns_caa_record":   resourceDnsCaaRecord(),
		"azurerm_dns_cname_record": resourceDnsCNameRecord(),
		"azurerm_dns_mx_record":    resourceDnsMxRecord(),
		"azurerm_dns_ns_record":    resourceDnsNsRecord(),
		"azurerm_dns_ptr_record":   resourceDnsPtrRecord(),
		"azurerm_dns_srv_record":   resourceDnsSrvRecord(),
		"azurerm_dns_txt_record":   resourceDnsTxtRecord(),
		"azurerm_dns_zone":         resourceDnsZone(),
	}
}
