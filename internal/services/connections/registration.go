// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package connections

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

var (
	_ sdk.TypedServiceRegistration                   = Registration{}
	_ sdk.UntypedServiceRegistrationWithAGitHubLabel = Registration{}
)

func (r Registration) AssociatedGitHubLabel() string {
	return "service/connections"
}

// Name is the name of this Service
func (r Registration) Name() string {
	return "Connections"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"Connections",
	}
}

// DataSources returns the typed Data Sources supported by this Service
func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{
		ApiConnectionDataSource{},
	}
}

// Resources returns the typed Resources supported by this Service
func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		// Add typed resources here when implemented
	}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_managed_api": dataSourceManagedApi(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_api_connection": resourceApiConnection(),
	}
}
