package logic

import (
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

// Name is the name of this Service
func (r Registration) Name() string {
	return "Logic"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"Logic App",
	}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_logic_app_workflow":            dataSourceLogicAppWorkflow(),
		"azurerm_logic_app_integration_account": dataSourceLogicAppIntegrationAccount(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_logic_app_action_custom":         resourceLogicAppActionCustom(),
		"azurerm_logic_app_action_http":           resourceLogicAppActionHTTP(),
		"azurerm_logic_app_integration_account":   resourceLogicAppIntegrationAccount(),
		"azurerm_logic_app_trigger_custom":        resourceLogicAppTriggerCustom(),
		"azurerm_logic_app_trigger_http_request":  resourceLogicAppTriggerHttpRequest(),
		"azurerm_logic_app_trigger_recurrence":    resourceLogicAppTriggerRecurrence(),
		"azurerm_logic_app_workflow":              resourceLogicAppWorkflow(),
		"azurerm_integration_service_environment": resourceIntegrationServiceEnvironment(),
	}
}
