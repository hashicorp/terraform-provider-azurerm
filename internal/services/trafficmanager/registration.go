// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package trafficmanager

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/trafficmanager/resources/traffic_manager_azure_endpoint"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/trafficmanager/resources/traffic_manager_external_endpoint"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/trafficmanager/resources/traffic_manager_geographical_location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/trafficmanager/resources/traffic_manager_nested_endpoint"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/trafficmanager/resources/traffic_manager_profile"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

var (
	_ sdk.FrameworkServiceRegistration               = Registration{}
	_ sdk.UntypedServiceRegistrationWithAGitHubLabel = Registration{}
)

func (r Registration) AssociatedGitHubLabel() string {
	return "service/traffic-manager"
}

// Name is the name of this Service
func (r Registration) Name() string {
	return "Traffic Manager"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"Network",
	}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_traffic_manager_geographical_location": traffic_manager_geographical_location.RegisterDataSource(),
		"azurerm_traffic_manager_profile":               traffic_manager_profile.RegisterDataSource(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_traffic_manager_azure_endpoint":    traffic_manager_azure_endpoint.RegisterResource(),
		"azurerm_traffic_manager_external_endpoint": traffic_manager_external_endpoint.RegisterResource(),
		"azurerm_traffic_manager_nested_endpoint":   traffic_manager_nested_endpoint.RegisterResource(),
		"azurerm_traffic_manager_profile":           traffic_manager_profile.RegisterResource(),
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
		traffic_manager_profile.ListResource{},
	}
}
