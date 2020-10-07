package resource

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

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
func (r Registration) SupportedDataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"azurerm_resources":      dataSourceArmResources(),
		"azurerm_resource_group": dataSourceArmResourceGroup(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"azurerm_management_lock":     resourceArmManagementLock(),
		"azurerm_resource_group":      resourceArmResourceGroup(),
		"azurerm_template_deployment": resourceArmTemplateDeployment(),

		"azurerm_resource_group_template_deployment": resourceGroupTemplateDeploymentResource(),
		"azurerm_subscription_template_deployment":   subscriptionTemplateDeploymentResource(),
	}
}
