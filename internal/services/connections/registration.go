// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package connections

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ sdk.UntypedServiceRegistration = Registration{}

type Registration struct{}

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

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_api_connection": dataSourceApiConnection(),
		"azurerm_managed_api":    dataSourceManagedApi(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_api_connection": resourceApiConnection(),
	}
}
