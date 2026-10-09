// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package virtualwan

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

var (
	_ sdk.TypedServiceRegistrationWithAGitHubLabel   = Registration{}
	_ sdk.UntypedServiceRegistrationWithAGitHubLabel = Registration{}
	_ sdk.FrameworkServiceRegistration               = Registration{}
)

// Name is the name of this Service
func (r Registration) Name() string {
	return "Virtual WAN"
}

func (r Registration) AssociatedGitHubLabel() string {
	return "service/virtual-wan"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"Network",
	}
}

func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{
		VPNServerConfigurationDataSource{},
	}
}

func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		RouteMapResource{},
		VirtualHubRoutingIntentResource{},
	}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_virtual_hub":             dataSourceVirtualHub(),
		"azurerm_virtual_hub_connection":  dataSourceVirtualHubConnection(),
		"azurerm_virtual_hub_route_table": dataSourceVirtualHubRouteTable(),
		"azurerm_virtual_wan":             dataSourceVirtualWan(),
		"azurerm_vpn_gateway":             dataSourceVPNGateway(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_express_route_connection":               resourceExpressRouteConnection(),
		"azurerm_express_route_gateway":                  resourceExpressRouteGateway(),
		"azurerm_point_to_site_vpn_gateway":               resourcePointToSiteVPNGateway(),
		"azurerm_route_server":                            resourceRouteServer(),
		"azurerm_route_server_bgp_connection":             resourceRouteServerBgpConnection(),
		"azurerm_virtual_hub":                             resourceVirtualHub(),
		"azurerm_virtual_hub_bgp_connection":              resourceVirtualHubBgpConnection(),
		"azurerm_virtual_hub_connection":                  resourceVirtualHubConnection(),
		"azurerm_virtual_hub_ip":                          resourceVirtualHubIP(),
		"azurerm_virtual_hub_route_table":                 resourceVirtualHubRouteTable(),
		"azurerm_virtual_hub_route_table_route":           resourceVirtualHubRouteTableRoute(),
		"azurerm_virtual_hub_security_partner_provider":   resourceVirtualHubSecurityPartnerProvider(),
		"azurerm_virtual_wan":                             resourceVirtualWan(),
		"azurerm_vpn_gateway":                             resourceVPNGateway(),
		"azurerm_vpn_gateway_connection":                  resourceVPNGatewayConnection(),
		"azurerm_vpn_gateway_nat_rule":                    resourceVPNGatewayNatRule(),
		"azurerm_vpn_server_configuration":                resourceVPNServerConfiguration(),
		"azurerm_vpn_server_configuration_policy_group":   resourceVPNServerConfigurationPolicyGroup(),
		"azurerm_vpn_site":                                resourceVpnSite(),
	}
}

func (r Registration) Actions() []func() action.Action {
	return []func() action.Action{}
}

func (r Registration) ListResources() []sdk.FrameworkListWrappedResource {
	return []sdk.FrameworkListWrappedResource{
		VirtualHubConnectionListResource{},
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
