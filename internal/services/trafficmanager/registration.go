// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package trafficmanager

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

var _ sdk.UntypedServiceRegistrationWithAGitHubLabel = Registration{}

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
		"azurerm_traffic_manager_geographical_location": dataSourceArmTrafficManagerGeographicalLocation(),
		"azurerm_traffic_manager_profile":               dataSourceArmTrafficManagerProfile(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_traffic_manager_azure_endpoint":    resourceAzureEndpoint(),
		"azurerm_traffic_manager_external_endpoint": resourceExternalEndpoint(),
		"azurerm_traffic_manager_nested_endpoint":   resourceNestedEndpoint(),
		"azurerm_traffic_manager_profile":           resourceArmTrafficManagerProfile(),
	}
}
