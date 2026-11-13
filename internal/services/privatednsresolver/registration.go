// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package privatednsresolver

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/privatednsresolver/resources/private_dns_resolver"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/privatednsresolver/resources/private_dns_resolver_dns_forwarding_ruleset"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/privatednsresolver/resources/private_dns_resolver_forwarding_rule"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/privatednsresolver/resources/private_dns_resolver_inbound_endpoint"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/privatednsresolver/resources/private_dns_resolver_outbound_endpoint"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/privatednsresolver/resources/private_dns_resolver_virtual_network_link"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

var (
	_ sdk.FrameworkServiceRegistration               = Registration{}
	_ sdk.TypedServiceRegistrationWithAGitHubLabel   = Registration{}
	_ sdk.UntypedServiceRegistrationWithAGitHubLabel = Registration{}
)

func (r Registration) AssociatedGitHubLabel() string {
	return "service/private-dns-resolver"
}

// Name is the name of this Service
func (r Registration) Name() string {
	return "Private DNS Resolver"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"Private DNS Resolver",
	}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{}
}

// DataSources returns a list of Data Sources supported by this Service
func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{
		private_dns_resolver_dns_forwarding_ruleset.DataSource{},
		private_dns_resolver.DataSource{},
		private_dns_resolver_forwarding_rule.DataSource{},
		private_dns_resolver_inbound_endpoint.DataSource{},
		private_dns_resolver_outbound_endpoint.DataSource{},
		private_dns_resolver_virtual_network_link.DataSource{},
	}
}

// Resources returns a list of Resources supported by this Service
func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		private_dns_resolver_dns_forwarding_ruleset.Resource{},
		private_dns_resolver.Resource{},
		private_dns_resolver_forwarding_rule.Resource{},
		private_dns_resolver_inbound_endpoint.Resource{},
		private_dns_resolver_outbound_endpoint.Resource{},
		private_dns_resolver_virtual_network_link.Resource{},
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
		private_dns_resolver_forwarding_rule.ListResource{},
	}
}
