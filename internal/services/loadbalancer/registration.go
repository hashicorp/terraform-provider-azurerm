// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package loadbalancer

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var (
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
		"azurerm_lb":                      dataSourceArmLoadBalancer(),
		"azurerm_lb_backend_address_pool": dataSourceArmLoadBalancerBackendAddressPool(),
		"azurerm_lb_rule":                 dataSourceArmLoadBalancerRule(),
		"azurerm_lb_outbound_rule":        dataSourceArmLoadBalancerOutboundRule(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_lb_backend_address_pool": resourceArmLoadBalancerBackendAddressPool(),
		"azurerm_lb_nat_pool":             resourceArmLoadBalancerNatPool(),
		"azurerm_lb_nat_rule":             resourceArmLoadBalancerNatRule(),
		"azurerm_lb_probe":                resourceArmLoadBalancerProbe(),
		"azurerm_lb_outbound_rule":        resourceArmLoadBalancerOutboundRule(),
		"azurerm_lb_rule":                 resourceArmLoadBalancerRule(),
		"azurerm_lb":                      resourceArmLoadBalancer(),
	}
}

// Resources returns a list of Resources supported by this Service
func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		BackendAddressPoolAddressResource{},
	}
}
