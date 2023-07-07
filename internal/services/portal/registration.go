// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package portal

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

var _ sdk.UntypedServiceRegistrationWithAGitHubLabel = Registration{}

func (r Registration) AssociatedGitHubLabel() string {
	return "service/portal"
}

// Name is the name of this Service
func (r Registration) Name() string {
	return "Portal"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"Portal",
	}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_portal_dashboard": dataSourcePortalDashboard(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	resources := map[string]*pluginsdk.Resource{
		"azurerm_portal_dashboard":            resourcePortalDashboard(),
		"azurerm_portal_tenant_configuration": resourcePortalTenantConfiguration(),
	}

	if !features.FourPointOhBeta() {
		resources["azurerm_dashboard"] = resourceLegacyDashboard()
	}

	return resources
}
