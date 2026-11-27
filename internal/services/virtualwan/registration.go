// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package virtualwan

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/virtualwan/resources/express_route_connection"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/virtualwan/resources/express_route_gateway"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/virtualwan/resources/point_to_site_vpn_gateway"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/virtualwan/resources/route_map"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/virtualwan/resources/route_server"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/virtualwan/resources/route_server_bgp_connection"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/virtualwan/resources/virtual_hub"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/virtualwan/resources/virtual_hub_bgp_connection"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/virtualwan/resources/virtual_hub_connection"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/virtualwan/resources/virtual_hub_ip"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/virtualwan/resources/virtual_hub_route_table"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/virtualwan/resources/virtual_hub_route_table_route"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/virtualwan/resources/virtual_hub_routing_intent"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/virtualwan/resources/virtual_hub_security_partner_provider"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/virtualwan/resources/virtual_wan"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/virtualwan/resources/vpn_gateway"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/virtualwan/resources/vpn_gateway_connection"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/virtualwan/resources/vpn_gateway_nat_rule"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/virtualwan/resources/vpn_server_configuration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/virtualwan/resources/vpn_server_configuration_policy_group"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/virtualwan/resources/vpn_site"
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
		vpn_server_configuration.DataSource{},
	}
}

func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		route_map.Resource{},
		virtual_hub_routing_intent.Resource{},
	}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_virtual_hub":             virtual_hub.RegisterDataSource(),
		"azurerm_virtual_hub_connection":  virtual_hub_connection.RegisterDataSource(),
		"azurerm_virtual_hub_route_table": virtual_hub_route_table.RegisterDataSource(),
		"azurerm_virtual_wan":             virtual_wan.RegisterDataSource(),
		"azurerm_vpn_gateway":             vpn_gateway.RegisterDataSource(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_express_route_connection":              express_route_connection.RegisterResource(),
		"azurerm_express_route_gateway":                 express_route_gateway.RegisterResource(),
		"azurerm_point_to_site_vpn_gateway":             point_to_site_vpn_gateway.RegisterResource(),
		"azurerm_route_server":                          route_server.RegisterResource(),
		"azurerm_route_server_bgp_connection":           route_server_bgp_connection.RegisterResource(),
		"azurerm_virtual_hub":                           virtual_hub.RegisterResource(),
		"azurerm_virtual_hub_bgp_connection":            virtual_hub_bgp_connection.RegisterResource(),
		"azurerm_virtual_hub_connection":                virtual_hub_connection.RegisterResource(),
		"azurerm_virtual_hub_ip":                        virtual_hub_ip.RegisterResource(),
		"azurerm_virtual_hub_route_table":               virtual_hub_route_table.RegisterResource(),
		"azurerm_virtual_hub_route_table_route":         virtual_hub_route_table_route.RegisterResource(),
		"azurerm_virtual_hub_security_partner_provider": virtual_hub_security_partner_provider.RegisterResource(),
		"azurerm_virtual_wan":                           virtual_wan.RegisterResource(),
		"azurerm_vpn_gateway":                           vpn_gateway.RegisterResource(),
		"azurerm_vpn_gateway_connection":                vpn_gateway_connection.RegisterResource(),
		"azurerm_vpn_gateway_nat_rule":                  vpn_gateway_nat_rule.RegisterResource(),
		"azurerm_vpn_server_configuration":              vpn_server_configuration.RegisterResource(),
		"azurerm_vpn_server_configuration_policy_group": vpn_server_configuration_policy_group.RegisterResource(),
		"azurerm_vpn_site":                              vpn_site.RegisterResource(),
	}
}

func (r Registration) Actions() []func() action.Action {
	return []func() action.Action{}
}

func (r Registration) ListResources() []sdk.FrameworkListWrappedResource {
	return []sdk.FrameworkListWrappedResource{
		virtual_hub_connection.ListResource{},
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
