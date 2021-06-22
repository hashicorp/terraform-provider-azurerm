package resource

import (
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/sdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
)

var _ sdk.TypedServiceRegistration = Registration{}
var _ sdk.UntypedServiceRegistration = Registration{}

type Registration struct{}

// Name is the name of this Service
func (r Registration) Name() string {
	return "Resources"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"Base",
		"Management",
		"Template",
	}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_resources":             dataSourceResources(),
		"azurerm_resource_group":        dataSourceResourceGroup(),
		"azurerm_template_spec_version": dataSourceTemplateSpecVersion(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_management_lock":                      resourceManagementLock(),
		"azurerm_management_group_template_deployment": managementGroupTemplateDeploymentResource(),
		"azurerm_resource_group":                       resourceResourceGroup(),
		"azurerm_resource_group_template_deployment":   resourceGroupTemplateDeploymentResource(),
		"azurerm_subscription_template_deployment":     subscriptionTemplateDeploymentResource(),
		"azurerm_template_deployment":                  resourceTemplateDeployment(),
		"azurerm_tenant_template_deployment":           tenantTemplateDeploymentResource(),
	}
}

// PackagePath is the relative path to this package
func (r Registration) PackagePath() string {
	return "TODO: do we need this?"
}

// DataSources returns a list of Data Sources supported by this Service
func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{}
}

// Resources returns a list of Resources supported by this Service
func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		ResourceProviderRegistrationResource{},
	}
}
