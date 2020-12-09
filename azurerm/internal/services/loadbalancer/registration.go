package loadbalancer

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type Registration struct{}

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

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"azurerm_lb":                      dataSourceArmLoadBalancer(),
		"azurerm_lb_backend_address_pool": dataSourceArmLoadBalancerBackendAddressPool(),
		"azurerm_lb_rule":                 dataSourceArmLoadBalancerRule(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"azurerm_lb_backend_address_pool": resourceArmLoadBalancerBackendAddressPool(),
		"azurerm_lb_nat_pool":             resourceArmLoadBalancerNatPool(),
		"azurerm_lb_nat_rule":             resourceArmLoadBalancerNatRule(),
		"azurerm_lb_probe":                resourceArmLoadBalancerProbe(),
		"azurerm_lb_outbound_rule":        resourceArmLoadBalancerOutboundRule(),
		"azurerm_lb_rule":                 resourceArmLoadBalancerRule(),
		"azurerm_lb":                      resourceArmLoadBalancer(),
	}
}
