// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package privatedns

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/privatedns/resources/private_dns_a_record"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/privatedns/resources/private_dns_aaaa_record"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/privatedns/resources/private_dns_cname_record"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/privatedns/resources/private_dns_mx_record"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/privatedns/resources/private_dns_ptr_record"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/privatedns/resources/private_dns_soa_record"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/privatedns/resources/private_dns_srv_record"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/privatedns/resources/private_dns_txt_record"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/privatedns/resources/private_dns_zone"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/privatedns/resources/private_dns_zone_virtual_network_link"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

func (r Registration) Actions() []func() action.Action {
	return []func() action.Action{}
}

func (r Registration) ListResources() []sdk.FrameworkListWrappedResource {
	return []sdk.FrameworkListWrappedResource{
		private_dns_a_record.ListResource{},
		private_dns_cname_record.ListResource{},
		private_dns_zone.ListResource{},
		private_dns_zone_virtual_network_link.ListResource{},
	}
}

func (r Registration) EphemeralResources() []func() ephemeral.EphemeralResource {
	return []func() ephemeral.EphemeralResource{}
}

func (r Registration) FrameworkDataSources() []sdk.FrameworkWrappedDataSource {
	return []sdk.FrameworkWrappedDataSource{}
}

func (r Registration) FrameworkResources() []sdk.FrameworkWrappedResource {
	return []sdk.FrameworkWrappedResource{}
}

var _ sdk.FrameworkServiceRegistration = Registration{}

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
		"azurerm_private_dns_a_record":                  private_dns_a_record.RegisterDataSource(),
		"azurerm_private_dns_aaaa_record":               private_dns_aaaa_record.RegisterDataSource(),
		"azurerm_private_dns_cname_record":              private_dns_cname_record.RegisterDataSource(),
		"azurerm_private_dns_mx_record":                 private_dns_mx_record.RegisterDataSource(),
		"azurerm_private_dns_ptr_record":                private_dns_ptr_record.RegisterDataSource(),
		"azurerm_private_dns_soa_record":                private_dns_soa_record.RegisterDataSource(),
		"azurerm_private_dns_srv_record":                private_dns_srv_record.RegisterDataSource(),
		"azurerm_private_dns_txt_record":                private_dns_txt_record.RegisterDataSource(),
		"azurerm_private_dns_zone":                      private_dns_zone.RegisterDataSource(),
		"azurerm_private_dns_zone_virtual_network_link": private_dns_zone_virtual_network_link.RegisterDataSource(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_private_dns_a_record":                  private_dns_a_record.RegisterResource(),
		"azurerm_private_dns_aaaa_record":               private_dns_aaaa_record.RegisterResource(),
		"azurerm_private_dns_cname_record":              private_dns_cname_record.RegisterResource(),
		"azurerm_private_dns_mx_record":                 private_dns_mx_record.RegisterResource(),
		"azurerm_private_dns_ptr_record":                private_dns_ptr_record.RegisterResource(),
		"azurerm_private_dns_srv_record":                private_dns_srv_record.RegisterResource(),
		"azurerm_private_dns_txt_record":                private_dns_txt_record.RegisterResource(),
		"azurerm_private_dns_zone":                      private_dns_zone.RegisterResource(),
		"azurerm_private_dns_zone_virtual_network_link": private_dns_zone_virtual_network_link.RegisterResource(),
	}
}
