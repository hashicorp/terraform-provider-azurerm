// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package dns

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

var _ sdk.UntypedServiceRegistrationWithAGitHubLabel = Registration{}

func (r Registration) AssociatedGitHubLabel() string {
	return "service/dns"
}

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
func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_dns_a_record":     dataSourceDnsARecord(),
		"azurerm_dns_aaaa_record":  dataSourceDnsAAAARecord(),
		"azurerm_dns_caa_record":   dataSourceDnsCaaRecord(),
		"azurerm_dns_cname_record": dataSourceDnsCNameRecord(),
		"azurerm_dns_mx_record":    dataSourceDnsMxRecord(),
		"azurerm_dns_ns_record":    dataSourceDnsNsRecord(),
		"azurerm_dns_ptr_record":   dataSourceDnsPtrRecord(),
		"azurerm_dns_soa_record":   dataSourceDnsSoaRecord(),
		"azurerm_dns_srv_record":   dataSourceDnsSrvRecord(),
		"azurerm_dns_txt_record":   dataSourceDnsTxtRecord(),
		"azurerm_dns_zone":         dataSourceDnsZone(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
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
