package logic

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type Registration struct{}

// Name is the name of this Service
func (r Registration) Name() string {
	return "Logic"
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"azurerm_logic_app_workflow": dataSourceArmLogicAppWorkflow(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"azurerm_logic_app_action_custom":        resourceArmLogicAppActionCustom(),
		"azurerm_logic_app_action_http":          resourceArmLogicAppActionHTTP(),
		"azurerm_logic_app_trigger_custom":       resourceArmLogicAppTriggerCustom(),
		"azurerm_logic_app_trigger_http_request": resourceArmLogicAppTriggerHttpRequest(),
		"azurerm_logic_app_trigger_recurrence":   resourceArmLogicAppTriggerRecurrence(),
		"azurerm_logic_app_workflow":             resourceArmLogicAppWorkflow()}
}
