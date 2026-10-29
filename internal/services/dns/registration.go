// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package dns

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/dns/resources/dns_a_record"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/dns/resources/dns_aaaa_record"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/dns/resources/dns_caa_record"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/dns/resources/dns_cname_record"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/dns/resources/dns_mx_record"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/dns/resources/dns_ns_record"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/dns/resources/dns_ptr_record"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/dns/resources/dns_soa_record"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/dns/resources/dns_srv_record"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/dns/resources/dns_txt_record"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/dns/resources/dns_zone"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

var (
	_ sdk.FrameworkServiceRegistration               = Registration{}
	_ sdk.TypedServiceRegistration                   = Registration{}
	_ sdk.UntypedServiceRegistrationWithAGitHubLabel = Registration{}
)

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
		"azurerm_dns_a_record":     dns_a_record.RegisterDataSource(),
		"azurerm_dns_aaaa_record":  dns_aaaa_record.RegisterDataSource(),
		"azurerm_dns_caa_record":   dns_caa_record.RegisterDataSource(),
		"azurerm_dns_cname_record": dns_cname_record.RegisterDataSource(),
		"azurerm_dns_mx_record":    dns_mx_record.RegisterDataSource(),
		"azurerm_dns_ns_record":    dns_ns_record.RegisterDataSource(),
		"azurerm_dns_ptr_record":   dns_ptr_record.RegisterDataSource(),
		"azurerm_dns_soa_record":   dns_soa_record.RegisterDataSource(),
		"azurerm_dns_srv_record":   dns_srv_record.RegisterDataSource(),
		"azurerm_dns_txt_record":   dns_txt_record.RegisterDataSource(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_dns_a_record":     dns_a_record.RegisterResource(),
		"azurerm_dns_aaaa_record":  dns_aaaa_record.RegisterResource(),
		"azurerm_dns_caa_record":   dns_caa_record.RegisterResource(),
		"azurerm_dns_cname_record": dns_cname_record.RegisterResource(),
		"azurerm_dns_mx_record":    dns_mx_record.RegisterResource(),
		"azurerm_dns_ns_record":    dns_ns_record.RegisterResource(),
		"azurerm_dns_ptr_record":   dns_ptr_record.RegisterResource(),
		"azurerm_dns_srv_record":   dns_srv_record.RegisterResource(),
		"azurerm_dns_txt_record":   dns_txt_record.RegisterResource(),
	}
}

// DataSources returns a list of Data Sources supported by this Service
func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{
		dns_zone.DataSource{},
	}
}

// Resources returns a list of Resources supported by this Service
func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		dns_zone.Resource{},
	}
}

func (r Registration) Actions() []func() action.Action {
	return []func() action.Action{}
}

func (r Registration) FrameworkResources() []sdk.FrameworkWrappedResource {
	return []sdk.FrameworkWrappedResource{}
}

func (r Registration) FrameworkDataSources() []sdk.FrameworkWrappedDataSource {
	return []sdk.FrameworkWrappedDataSource{}
}

func (r Registration) EphemeralResources() []func() ephemeral.EphemeralResource {
	return []func() ephemeral.EphemeralResource{}
}

func (r Registration) ListResources() []sdk.FrameworkListWrappedResource {
	return []sdk.FrameworkListWrappedResource{
		dns_a_record.ListResource{},
		dns_cname_record.ListResource{},
		dns_txt_record.ListResource{},
	}
}
