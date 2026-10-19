// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package expressroute

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
	return "ExpressRoute"
}

func (r Registration) AssociatedGitHubLabel() string {
	return "service/express-route"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"Network",
	}
}

func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{}
}

func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_express_route_circuit":         dataSourceExpressRouteCircuit(),
		"azurerm_express_route_circuit_peering": dataSourceExpressRouteCircuitPeering(),
		"azurerm_route_filter":                  dataSourceRouteFilter(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_express_route_circuit":               resourceExpressRouteCircuit(),
		"azurerm_express_route_circuit_authorization": resourceExpressRouteCircuitAuthorization(),
		"azurerm_express_route_circuit_connection":    resourceExpressRouteCircuitConnection(),
		"azurerm_express_route_circuit_peering":       resourceExpressRouteCircuitPeering(),
		"azurerm_express_route_port":                  resourceArmExpressRoutePort(),
		"azurerm_express_route_port_authorization":    resourceExpressRoutePortAuthorization(),
		"azurerm_route_filter":                        resourceRouteFilter(),
	}
}

func (r Registration) Actions() []func() action.Action {
	return []func() action.Action{}
}

func (r Registration) ListResources() []sdk.FrameworkListWrappedResource {
	return []sdk.FrameworkListWrappedResource{}
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
