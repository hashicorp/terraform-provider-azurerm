// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package loadbalancer

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loadbalancer/resources/lb"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loadbalancer/resources/lb_backend_address_pool"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loadbalancer/resources/lb_backend_address_pool_address"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loadbalancer/resources/lb_nat_pool"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loadbalancer/resources/lb_nat_rule"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loadbalancer/resources/lb_outbound_rule"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loadbalancer/resources/lb_probe"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loadbalancer/resources/lb_rule"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var (
	_ sdk.FrameworkServiceRegistration             = Registration{}
	_ sdk.TypedServiceRegistrationWithAGitHubLabel = Registration{}
	_ sdk.UntypedServiceRegistration               = Registration{}
)

type Registration struct{}

func (r Registration) AssociatedGitHubLabel() string {
	return "service/load-balancers"
}

// Name is the name of this Service
func (r Registration) Name() string {
	return "Load Balancer"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"Load Balancer",
	}
}

// DataSources returns a list of Data Sources supported by this Service
func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_lb":                      lb.RegisterDataSource(),
		"azurerm_lb_backend_address_pool": lb_backend_address_pool.RegisterDataSource(),
		"azurerm_lb_outbound_rule":        lb_outbound_rule.RegisterDataSource(),
		"azurerm_lb_rule":                 lb_rule.RegisterDataSource(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_lb":                      lb.RegisterResource(),
		"azurerm_lb_backend_address_pool": lb_backend_address_pool.RegisterResource(),
		"azurerm_lb_nat_pool":             lb_nat_pool.RegisterResource(),
		"azurerm_lb_nat_rule":             lb_nat_rule.RegisterResource(),
		"azurerm_lb_outbound_rule":        lb_outbound_rule.RegisterResource(),
		"azurerm_lb_probe":                lb_probe.RegisterResource(),
		"azurerm_lb_rule":                 lb_rule.RegisterResource(),
	}
}

// Resources returns a list of Resources supported by this Service
func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		lb_backend_address_pool_address.Resource{},
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
		lb_backend_address_pool.ListResource{},
		lb_probe.ListResource{},
	}
}
