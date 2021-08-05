package msi

import (
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
)

// TODO: we should probably rename this Identity, or move into Authorization

type Registration struct{}

// Name is the name of this Service
func (r Registration) Name() string {
	return "Managed Service Identities"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"Authorization",
	}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_user_assigned_identity": dataSourceArmUserAssignedIdentity(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_user_assigned_identity": resourceArmUserAssignedIdentity(),
	}
}
