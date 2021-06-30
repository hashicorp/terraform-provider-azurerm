package loadbalancer

import (
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/sdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
)

var _ sdk.TypedServiceRegistration = Registration{}
var _ sdk.UntypedServiceRegistration = Registration{}

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

// PackagePath is the relative path to this package
func (r Registration) PackagePath() string {
	return "TODO: do we need this?"
}

// Resources returns a list of Resources supported by this Service
func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		BackendAddressPoolAddressResource{},
	}
}
